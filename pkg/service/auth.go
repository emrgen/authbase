package service

import (
	"context"
	"encoding/json"
	"errors"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/cache"
	"github.com/emrgen/authbase/pkg/model"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/emrgen/authbase/x"
	"github.com/emrgen/authbase/x/mail"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

var _ v1.AuthServiceServer = new(AuthService)

// AuthService is a service that implements the AuthServiceServer interface
type AuthService struct {
	store  store.Provider
	mailer mail.MailerProvider
	cache  *cache.Redis
	v1.UnimplementedAuthServiceServer
}

// NewAuthService creates a new AuthService
func NewAuthService(store store.Provider, mailer mail.MailerProvider, cache *cache.Redis) *AuthService {
	return &AuthService{store: store, mailer: mailer, cache: cache}
}

// CheckUserAlreadyExists checks if a user with the given email or username already exists in an organization
// NOTE: should be rate limited
func (a *AuthService) CheckUserAlreadyExists(ctx context.Context, request *v1.CheckUserAlreadyExistsRequest) (*v1.CheckEmailAlreadyExistsResponse, error) {
	orgID, err := uuid.Parse(request.GetOrganizationId())
	if err != nil {
		return nil, err
	}

	authStore, err := a.store.Provide(orgID)
	if err != nil {
		return nil, err
	}
	users, err := authStore.UserExists(ctx, orgID, request.GetUsername(), request.GetEmail())
	var emailExists bool
	var usernameExists bool
	for _, user := range users {
		if user.Email == request.GetEmail() {
			emailExists = true
		}

		if user.Username == request.GetUsername() {
			usernameExists = true
		}
	}

	return &v1.CheckEmailAlreadyExistsResponse{
		EmailExists:    emailExists,
		UsernameExists: usernameExists,
	}, nil

}

// Register creates a new user if the username and email are unique
func (a *AuthService) Register(ctx context.Context, request *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	email := request.GetEmail()
	username := request.GetUsername()
	password := request.GetPassword()
	orgID := uuid.MustParse(request.GetOrganizationId())

	authStore, err := a.store.Provide(orgID)
	if err != nil {
		return nil, err
	}

	users, err := authStore.UserExists(ctx, orgID, username, email)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.Email == email {
			return nil, errors.New("email already exists")
		}

		if user.Username == username {
			return nil, errors.New("username already exists")
		}
	}

	salt := x.Keygen()
	hashedPassword, _ := x.HashPassword(password, salt)

	user := &model.User{
		OrganizationID: orgID.String(),
		Username:       username,
		Email:          email,
		Password:       string(hashedPassword),
		Salt:           salt,
		Verified:       false,
	}

	err = authStore.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// generate a verification code
	code := x.GenerateCode()
	expireAt := time.Now().Add(24 * time.Hour)

	err = authStore.CreateVerificationCode(ctx, &model.VerificationCode{
		Code:      code,
		UserID:    user.ID,
		ExpiresAt: expireAt,
	})
	if err != nil {
		return nil, err
	}

	// send email verification code to the user
	err = a.mailer.Provide(orgID).SendMail(email, email, "Verify your email", "verify-email")
	if err != nil {
		return nil, err
	}

	return &v1.RegisterResponse{
		Message: "user registered",
	}, nil
}

// Login logs in a user and returns an access token and a refresh token
func (a *AuthService) Login(ctx context.Context, request *v1.LoginRequest) (*v1.LoginResponse, error) {
	email := request.GetEmail()
	password := request.GetPassword()

	orgID := uuid.MustParse(request.GetOrganizationId())
	authStore, err := a.store.Provide(orgID)
	if err != nil {
		return nil, err
	}

	user, err := authStore.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	ok := x.CompareHashAndPassword(user.Password, password, user.Salt)
	if !ok {
		return nil, errors.New("incorrect password")
	}

	// generate tokens
	jti := uuid.New().String()
	token, err := x.GenerateJWTToken(user.ID, user.OrganizationID, jti)
	if err != nil {
		return nil, err
	}
	logrus.Info("token: ", token)
	refreshExpireAt := time.Now().Add(5 * 24 * time.Hour)
	// save refresh token to cache, it will be used to validate the refresh token request
	err = a.cache.Set(jti, user.ID, 5*24*time.Hour)
	if err != nil {
		return nil, err
	}

	// save the token to the db
	// TODO: save the token to the db in encrypted form
	err = authStore.Transaction(func(tx store.AuthBaseStore) error {
		err = authStore.CreateRefreshToken(ctx, &model.RefreshToken{
			Token:          token.RefreshToken,
			OrganizationID: user.OrganizationID,
			UserID:         user.ID,
			ExpireAt:       refreshExpireAt,
			IssuedAt:       token.IssuedAt,
		})
		if err != nil {
			return err
		}

		// create a new session
		// user the jti as the session id
		// this will allow multiple sessions for a user at the same time
		err = authStore.CreateSession(ctx, &model.Session{
			ID:             jti,
			UserID:         user.ID,
			OrganizationID: user.OrganizationID,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// return tokens
	return &v1.LoginResponse{
		User: &v1.User{
			Id:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
		Token: &v1.AuthToken{
			AccessToken:      token.AccessToken,
			RefreshToken:     token.RefreshToken,
			ExpiresAt:        timestamppb.New(token.ExpireAt),
			IssuedAt:         timestamppb.New(token.IssuedAt),
			RefreshExpiresAt: timestamppb.New(refreshExpireAt),
		},
	}, nil
}

func (a *AuthService) Logout(ctx context.Context, request *v1.LogoutRequest) (*v1.LogoutResponse, error) {
	// check if the token is still valid
	accessToken := request.GetAccessToken()
	claims, err := x.VerifyJWTToken(accessToken)
	if err != nil {
		return nil, err
	}

	jti := uuid.MustParse(claims.Jti)
	orgID := uuid.MustParse(claims.OrgID)
	authStore, err := a.store.Provide(orgID)
	if err != nil {
		return nil, err
	}

	err = authStore.DeleteSession(ctx, jti)
	if err != nil {
		return nil, err
	}

	return &v1.LogoutResponse{Message: "logged out"}, nil

}

// Refresh generates a new access token using the refresh token
func (a *AuthService) Refresh(ctx context.Context, request *v1.RefreshRequest) (*v1.RefreshResponse, error) {
	var foundToken bool
	var userID string
	var organizationID string

	// check if the refresh token is still valid
	refreshToken := request.GetRefreshToken()
	claims, err := x.VerifyJWTToken(refreshToken)
	if err != nil {
		return nil, err
	}

	tokenStr, err := a.cache.Get(claims.Jti)
	if err != nil {
		// if no value in cache check the db
		return nil, err
	}

	orgID := uuid.MustParse(request.GetOrganizationId())
	authStore, err := a.store.Provide(orgID)
	if err != nil {
		return nil, err
	}

	if tokenStr != "" {
		var token model.Token
		err := json.Unmarshal([]byte(tokenStr), &token)
		if err != nil {
			return nil, err
		}
		userID = token.UserID
		organizationID = token.OrganizationID
		foundToken = true
	}

	if !foundToken {
		// check the db
		token, err := authStore.GetRefreshTokenByID(ctx, refreshToken)
		if err != nil {
			return nil, err
		}

		if token == nil {
			return nil, errors.New("token not found, need to login again")
		}

		foundToken = true
		userID = token.UserID
		organizationID = token.OrganizationID
	}

	jwtToken, err := x.GenerateJWTToken(userID, organizationID, claims.Jti)
	if err != nil {
		return nil, err
	}

	return &v1.RefreshResponse{
		AccessToken:  jwtToken.AccessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    timestamppb.New(jwtToken.ExpireAt),
		IssuedAt:     timestamppb.New(jwtToken.IssuedAt),
	}, nil
}

// VerifyEmail verifies the email of a user and sets the verified field to true
func (a *AuthService) VerifyEmail(ctx context.Context, request *v1.VerifyEmailRequest) (*v1.VerifyEmailResponse, error) {
	var err error
	orgID := uuid.MustParse(request.GetOrganizationId())
	authStore, err := a.store.Provide(orgID)
	if err != nil {
		return nil, err
	}

	// check if the email is already verified
	err = authStore.Transaction(func(tx store.AuthBaseStore) error {
		code, err := tx.GetVerificationCode(ctx, request.GetToken())
		if err != nil {
			return err
		}

		if code.ExpiresAt.Before(time.Now()) {
			return errors.New("verification code has expired")
		}

		user, err := tx.GetUserByID(ctx, uuid.MustParse(code.UserID))
		if err != nil {
			return err
		}

		user.Verified = true
		user.VerifiedAt = time.Now().UTC()

		err = tx.UpdateUser(ctx, user)
		if err != nil {
			return err
		}

		err = tx.DeleteVerificationCode(ctx, request.GetToken())
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// if it is, return an error
	return &v1.VerifyEmailResponse{Message: "email verified"}, nil
}
