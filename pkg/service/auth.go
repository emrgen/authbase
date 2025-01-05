package service

import (
	"context"
	"encoding/json"
	"errors"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/cache"
	"github.com/emrgen/authbase/pkg/model"
	"github.com/emrgen/authbase/pkg/permission"
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
	perm   permission.AuthBasePermission
	v1.UnimplementedAuthServiceServer
}

// NewAuthService creates a new AuthService
func NewAuthService(store store.Provider, perm permission.AuthBasePermission, mailer mail.MailerProvider, cache *cache.Redis) *AuthService {
	return &AuthService{store: store, perm: perm, mailer: mailer, cache: cache}
}

// UserEmailExists checks if a user with the given email or username already exists in an project
// NOTE: should be rate limited by ip address
func (a *AuthService) UserEmailExists(ctx context.Context, request *v1.UserEmailExistsRequest) (*v1.UserEmailExistsResponse, error) {
	orgID, err := uuid.Parse(request.GetProjectId())
	if err != nil {
		return nil, err
	}

	as, err := store.GetProjectStore(ctx, a.store)
	if err != nil {
		return nil, err
	}
	users, err := as.UserExists(ctx, orgID, request.GetUsername(), request.GetEmail())
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

	return &v1.UserEmailExistsResponse{
		EmailExists:    emailExists,
		UsernameExists: usernameExists,
	}, nil

}

// Register creates a new user if the username and email are unique
func (a *AuthService) Register(ctx context.Context, request *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	email := request.GetEmail()
	username := request.GetUsername()
	password := request.GetPassword()
	orgID := uuid.MustParse(request.GetProjectId())

	as, err := store.GetProjectStore(ctx, a.store)
	if err != nil {
		return nil, err
	}

	users, err := as.UserExists(ctx, orgID, username, email)
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
		ID:        uuid.New().String(),
		ProjectID: orgID.String(),
		Username:  username,
		Email:     email,
		Password:  string(hashedPassword),
		Salt:      salt,
		Verified:  false,
	}

	err = as.Transaction(func(tx store.AuthBaseStore) error {
		err = as.CreateUser(ctx, user)
		if err != nil {
			return err
		}

		// generate a verification code
		code := x.GenerateCode()
		expireAt := time.Now().Add(24 * time.Hour)

		err = as.CreateVerificationCode(ctx, &model.VerificationCode{
			ID:        uuid.New().String(),
			Code:      code,
			UserID:    user.ID,
			ExpiresAt: expireAt,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// send email verification code to the user
	go func() {
		logrus.Infof("sending email to %s", email)
		err = a.mailer.Provide(orgID).SendMail(email, email, "Verify your email", "verify-email")
		if err != nil {
			logrus.Errorf("failed to send email: %v", err)
		}
	}()

	return &v1.RegisterResponse{
		Message: "user registered",
	}, nil
}

// Login logs in a user and returns an access token and a refresh token
func (a *AuthService) Login(ctx context.Context, request *v1.LoginRequest) (*v1.LoginResponse, error) {
	email := request.GetEmail()
	password := request.GetPassword()

	orgID := uuid.MustParse(request.GetProjectId())
	as, err := store.GetProjectStore(ctx, a.store)
	if err != nil {
		return nil, err
	}

	user, err := as.GetUserByEmail(ctx, orgID, email)
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

	perm, err := as.GetProjectMemberByID(ctx, orgID, uuid.MustParse(user.ID))

	// generate tokens
	jti := uuid.New().String()
	token, err := x.GenerateJWTToken(x.Claims{
		Username:   user.Username,
		Email:      user.Email,
		ProjectID:  user.ProjectID,
		UserID:     user.ID,
		Permission: perm.Permission,
		Audience:   "",
		Jti:        jti,
		ExpireAt:   time.Now().Add(5 * time.Hour),
		IssuedAt:   time.Now(),
		Provider:   "authbase",
		Data:       nil,
	})
	if err != nil {
		return nil, err
	}

	refreshExpireAt := time.Now().Add(5 * 24 * time.Hour)
	// save refresh token to cache, it will be used to validate the refresh token request
	err = a.cache.Set(jti, user.ID, 5*24*time.Hour)
	if err != nil {
		return nil, err
	}

	// save the token to the db
	// TODO: save the token to the db in encrypted form
	err = as.Transaction(func(tx store.AuthBaseStore) error {
		err = as.CreateRefreshToken(ctx, &model.RefreshToken{
			Token:     token.RefreshToken,
			ProjectID: user.ProjectID,
			UserID:    user.ID,
			ExpireAt:  refreshExpireAt,
			IssuedAt:  token.IssuedAt,
		})
		if err != nil {
			return err
		}

		// create a new session
		// user the jti as the session id
		// this will allow multiple sessions for a user at the same time
		err = as.CreateSession(ctx, &model.Session{
			ID:        jti,
			UserID:    user.ID,
			ProjectID: user.ProjectID,
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
	as, err := store.GetProjectStore(ctx, a.store)
	if err != nil {
		return nil, err
	}

	// delete the token from the cache
	err = a.cache.Del(claims.Jti)
	if err != nil {
		return nil, err
	}

	// delete the session from the db
	err = as.DeleteSession(ctx, jti)
	if err != nil {
		return nil, err
	}

	return &v1.LogoutResponse{Message: "logged out"}, nil
}

// RevokeAllSessions logs out all sessions of a user
func (a *AuthService) RevokeAllSessions(ctx context.Context, request *v1.RevokeAllSessionsRequest) (*v1.RevokeAllSessionsResponse, error) {
	var err error

	userID := uuid.MustParse(request.GetUserId())
	as, err := store.GetProjectStore(ctx, a.store)
	if err != nil {
		return nil, err
	}

	// user project id
	user, err := as.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	orgID := uuid.MustParse(user.ProjectID)

	err = a.perm.CheckProjectPermission(ctx, orgID, "write")
	if err != nil {
		return nil, err
	}

	sessions, err := as.ListActiveSessions(ctx, userID)
	err = as.DeleteSessionByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &v1.RevokeAllSessionsResponse{
		SessionCount: uint64(len(sessions)),
		Message:      "logged out of all sessions",
	}, nil
}

// Refresh generates a new access token using the refresh token
func (a *AuthService) Refresh(ctx context.Context, request *v1.RefreshRequest) (*v1.RefreshResponse, error) {
	var foundToken bool
	var userID string
	var projectID string

	// check if the refresh token is still valid
	refreshToken := request.GetRefreshToken()
	claims, err := x.VerifyJWTToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// check if the token is in the cache
	tokenStr, err := a.cache.Get(claims.Jti)
	if err != nil {
		// if no value in cache check the db
		return nil, err
	}

	orgID := uuid.MustParse(request.GetProjectId())
	as, err := store.GetProjectStore(ctx, a.store)
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
		projectID = token.ProjectID
		foundToken = true
	}

	if !foundToken {
		// check the db
		token, err := as.GetRefreshTokenByID(ctx, refreshToken)
		if err != nil {
			return nil, err
		}

		if token == nil {
			return nil, errors.New("token not found, need to login again")
		}

		foundToken = true
		userID = token.UserID
		projectID = token.ProjectID
	}

	user, err := as.GetUserByID(ctx, uuid.MustParse(userID))
	if err != nil {
		return nil, err
	}

	perm, err := as.GetProjectMemberByID(ctx, orgID, uuid.MustParse(userID))
	if err != nil {
		return nil, err
	}

	jti := uuid.New().String()
	token, err := x.GenerateJWTToken(x.Claims{
		ProjectID:  projectID,
		UserID:     user.ID,
		Username:   claims.Username,
		Email:      claims.Email,
		Permission: perm.Permission,
		Audience:   claims.Audience,
		Jti:        jti,
		ExpireAt:   time.Now().Add(5 * time.Hour),
		IssuedAt:   time.Now(),
		Provider:   "authbase",
		Data:       nil,
	})
	if err != nil {
		return nil, err
	}

	return &v1.RefreshResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    timestamppb.New(token.ExpireAt),
		IssuedAt:     timestamppb.New(token.IssuedAt),
	}, nil
}

// VerifyEmail verifies the email of a user and sets the verified field to true
func (a *AuthService) VerifyEmail(ctx context.Context, request *v1.VerifyEmailRequest) (*v1.VerifyEmailResponse, error) {
	var err error
	as, err := store.GetProjectStore(ctx, a.store)
	if err != nil {
		return nil, err
	}

	// check if the email is already verified
	err = as.Transaction(func(tx store.AuthBaseStore) error {
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
