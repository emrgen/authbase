package service

import (
	"context"
	"encoding/json"
	"errors"
	goset "github.com/deckarep/golang-set/v2"
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

// NewAuthService creates a new AuthService
func NewAuthService(store store.Provider, perm permission.AuthBasePermission, mailer mail.MailerProvider, cache *cache.Redis) *AuthService {
	return &AuthService{store: store, perm: perm, mailer: mailer, cache: cache}
}

var _ v1.AuthServiceServer = new(AuthService)

// AuthService is a service that implements the AuthServiceServer interface
type AuthService struct {
	store  store.Provider
	mailer mail.MailerProvider
	cache  *cache.Redis
	perm   permission.AuthBasePermission
	v1.UnimplementedAuthServiceServer
}

func (a *AuthService) AccountEmailExists(ctx context.Context, request *v1.AccountEmailExistsRequest) (*v1.AccountEmailExistsResponse, error) {
	orgID, err := uuid.Parse(request.GetProjectId())
	if err != nil {
		return nil, err
	}

	as, err := store.GetProjectStore(ctx, a.store)
	if err != nil {
		return nil, err
	}
	users, err := as.AccountExists(ctx, orgID, request.GetUsername(), request.GetEmail())
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

	return &v1.AccountEmailExistsResponse{
		EmailExists:    emailExists,
		UsernameExists: usernameExists,
	}, nil
}

// LoginUsingIdp redirects the user to the identity provider for login
func (a *AuthService) LoginUsingIdp(ctx context.Context, request *v1.LoginUsingPasswordRequest) (*v1.LoginUsingPasswordResponse, error) {
	//TODO implement me
	panic("implement me")
}

// GetIdpToken gets the token from the identity provider
func (a *AuthService) GetIdpToken(ctx context.Context, request *v1.GetIdpTokenRequest) (*v1.GetIdpTokenResponse, error) {
	//TODO implement me
	panic("implement me")
}

// RegisterUsingPassword registers a user using a username, email, and password
func (a *AuthService) RegisterUsingPassword(ctx context.Context, request *v1.RegisterUsingPasswordRequest) (*v1.RegisterUsingPasswordResponse, error) {
	email := request.GetEmail()
	username := request.GetUsername()
	password := request.GetPassword()
	orgID := uuid.MustParse(request.GetProjectId())

	as, err := store.GetProjectStore(ctx, a.store)
	if err != nil {
		return nil, err
	}

	users, err := as.AccountExists(ctx, orgID, username, email)
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

	user := &model.Account{
		ID:        uuid.New().String(),
		ProjectID: orgID.String(),
		Username:  username,
		Email:     email,
		Password:  string(hashedPassword),
		Salt:      salt,
		Verified:  false,
	}

	err = as.Transaction(func(tx store.AuthBaseStore) error {
		err = as.CreateAccount(ctx, user)
		if err != nil {
			return err
		}

		// generate a verification code
		code := x.GenerateVerificationCode()
		expireAt := time.Now().Add(24 * time.Hour)

		err = as.CreateVerificationCode(ctx, &model.VerificationCode{
			ID:        uuid.New().String(),
			Code:      code,
			AccountID: user.ID,
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

	return &v1.RegisterUsingPasswordResponse{
		Message: "user registered",
	}, nil
}

// LoginUsingPassword logs in a user and returns an access token and a refresh token
// TODO: should be rate limited to prevent brute force attacks
func (a *AuthService) LoginUsingPassword(ctx context.Context, request *v1.LoginUsingPasswordRequest) (*v1.LoginUsingPasswordResponse, error) {
	email := request.GetEmail()
	password := request.GetPassword()

	clientID := uuid.MustParse(request.GetClientId())
	as, err := store.GetProjectStore(ctx, a.store)
	if err != nil {
		return nil, err
	}
	client, err := as.GetClientByID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	poolID := uuid.MustParse(client.PoolID)

	account, err := as.GetAccountByEmail(ctx, poolID, email)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, errors.New("account not found")
	}

	ok := x.CompareHashAndPassword(account.Password, password, account.Salt)
	if !ok {
		return nil, errors.New("incorrect password")
	}

	//if account.ProjectMember {
	//	perm, err = as.GetProjectMemberByID(ctx, clientID, uuid.MustParse(account.ID))
	//	if err != nil {
	//		return nil, err
	//	}
	//}

	// get the account scopes from the memberships
	accountID := uuid.MustParse(account.ID)
	memberships, err := as.ListGroupMemberByAccount(ctx, poolID, accountID)
	if err != nil {
		return nil, err
	}
	set := goset.NewSet[string]()
	for _, member := range memberships {
		for _, role := range member.Group.Roles {
			set.Add(role.Name)
		}
	}

	roleNames := set.ToSlice()

	// generate tokens
	jti := uuid.New().String()
	token, err := x.GenerateJWTToken(x.Claims{
		Username:  account.Username,
		Email:     account.Email,
		ProjectID: account.ProjectID,
		AccountID: account.ID,
		Audience:  "", // the target website or app that will use the token
		Jti:       jti,
		ExpireAt:  time.Now().Add(x.AccessTokenDuration),
		IssuedAt:  time.Now(),
		Provider:  "authbase",
		Scopes:    roleNames,
	})
	if err != nil {
		return nil, err
	}

	refreshExpireAt := time.Now().Add(x.RefreshTokenDuration)
	// save refresh token to cache, it will be used to validate the refresh token request
	err = a.cache.Set(jti, account.ID, x.RefreshTokenDuration)
	if err != nil {
		return nil, err
	}

	// save the token to the db
	// TODO: save the token to the db in encrypted form
	err = as.Transaction(func(tx store.AuthBaseStore) error {
		err = as.CreateRefreshToken(ctx, &model.RefreshToken{
			Token:     token.RefreshToken,
			ProjectID: account.ProjectID,
			AccountID: account.ID,
			ExpireAt:  refreshExpireAt,
			IssuedAt:  token.IssuedAt,
		})
		if err != nil {
			return err
		}

		// create a new session
		// account the jti as the session id
		// this will allow multiple sessions for a account at the same time
		err = as.CreateSession(ctx, &model.Session{
			ID:        jti,
			AccountID: account.ID,
			ProjectID: account.ProjectID,
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
	return &v1.LoginUsingPasswordResponse{
		User: &v1.Account{
			Id:        account.ID,
			Username:  account.Username,
			Email:     account.Email,
			CreatedAt: timestamppb.New(account.CreatedAt),
			UpdatedAt: timestamppb.New(account.UpdatedAt),
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

// Logout logs out a user by deleting the session from the db, and the refresh tokens from the cache
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

// ForgotPassword sends a password reset link to the user's email or phone number to reset their password
func (a *AuthService) ForgotPassword(ctx context.Context, request *v1.ForgotPasswordRequest) (*v1.ForgotPasswordResponse, error) {
	projectID, err := x.GetAuthbaseProjectID(ctx)
	if err != nil {
		return nil, err
	}

	as, err := store.GetProjectStore(ctx, a.store)
	if err != nil {
		return nil, err
	}

	mailer := a.mailer.Provide(projectID)
	email := request.GetEmail()

	account, err := as.GetAccountByEmail(ctx, projectID, email)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, errors.New("account not found")
	}

	code := x.GenerateVerificationCode()
	expireAt := time.Now().Add(24 * time.Hour)

	err = a.cache.Set(code, account.ID, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	err = as.CreateVerificationCode(ctx, &model.VerificationCode{
		ID:        uuid.New().String(),
		Code:      code,
		AccountID: account.ID,
		ProjectID: account.ProjectID,
		PoolID:    account.PoolID,
		ExpiresAt: expireAt,
	})
	if err != nil {
		return nil, err
	}

	// send email verification code to the user
	go func() {
		logrus.Infof("sending email to %s", email)
		err = mailer.SendMail(email, email, "Reset your password", "reset-password")
		if err != nil {
			logrus.Errorf("failed to send email: %v", err)
		}
	}()

	return &v1.ForgotPasswordResponse{Message: "password reset link sent"}, nil
}

// ResetPassword resets the password of a user using a verification code
func (a *AuthService) ResetPassword(ctx context.Context, request *v1.ResetPasswordRequest) (*v1.ResetPasswordResponse, error) {
	code := request.GetCode()
	password := request.GetNewPassword()

	as, err := store.GetProjectStore(ctx, a.store)
	if err != nil {
		return nil, err
	}

	err = as.Transaction(func(tx store.AuthBaseStore) error {
		code, err := tx.GetVerificationCode(ctx, code)
		if err != nil {
			return err
		}

		if code.ExpiresAt.Before(time.Now()) {
			return errors.New("verification code has expired")
		}

		account, err := tx.GetAccountByID(ctx, uuid.MustParse(code.AccountID))
		if err != nil {
			return err
		}

		salt := x.Keygen()
		hashedPassword, _ := x.HashPassword(password, salt)
		account.Password = string(hashedPassword)
		account.Salt = salt

		// update the account with the new password
		err = tx.UpdateAccount(ctx, account)
		if err != nil {
			return err
		}

		// delete the verification code
		err = tx.DeleteVerificationCode(ctx, code.ID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// TODO: should redirect to the login page
	return &v1.ResetPasswordResponse{Message: "password reset"}, nil
}

// ChangePassword changes the password of a user
func (a *AuthService) ChangePassword(ctx context.Context, request *v1.ChangePasswordRequest) (*v1.ChangePasswordResponse, error) {
	accountID, err := x.GetAuthbaseAccountID(ctx)
	if err != nil {
		return nil, err
	}

	as, err := store.GetProjectStore(ctx, a.store)
	if err != nil {
		return nil, err
	}

	err = as.Transaction(func(tx store.AuthBaseStore) error {
		account, err := tx.GetAccountByID(ctx, accountID)
		if err != nil {
			return err
		}

		salt := x.Keygen()
		hashedPassword, _ := x.HashPassword(request.GetNewPassword(), salt)
		account.Password = string(hashedPassword)
		account.Salt = salt

		err = tx.UpdateAccount(ctx, account)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.ChangePasswordResponse{Message: "password changed"}, nil
}

// Refresh generates a new access token using the refresh token
func (a *AuthService) Refresh(ctx context.Context, request *v1.RefreshRequest) (*v1.RefreshResponse, error) {
	var foundToken bool
	var accountID string
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

	as, err := store.GetProjectStore(ctx, a.store)
	if err != nil {
		return nil, err
	}

	if tokenStr != "" {
		var token model.RefreshToken
		err := json.Unmarshal([]byte(tokenStr), &token)
		if err != nil {
			return nil, err
		}
		accountID = token.AccountID
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
		accountID = token.AccountID
		projectID = token.ProjectID
	}

	user, err := as.GetAccountByID(ctx, uuid.MustParse(accountID))
	if err != nil {
		return nil, err
	}

	jti := uuid.New().String()
	token, err := x.GenerateJWTToken(x.Claims{
		ProjectID: projectID,
		AccountID: user.ID,
		Username:  claims.Username,
		Email:     claims.Email,
		Audience:  claims.Audience,
		Jti:       jti,
		ExpireAt:  time.Now().Add(15 * time.Minute),
		IssuedAt:  time.Now(),
		Provider:  "authbase",
		Scopes:    claims.Scopes,
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

		user, err := tx.GetAccountByID(ctx, uuid.MustParse(code.AccountID))
		if err != nil {
			return err
		}

		user.Verified = true
		user.VerifiedAt = time.Now().UTC()

		err = tx.UpdateAccount(ctx, user)
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
