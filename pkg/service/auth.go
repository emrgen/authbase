package service

import (
	"context"
	"errors"
	goset "github.com/deckarep/golang-set/v2"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/cache"
	"github.com/emrgen/authbase/pkg/model"
	"github.com/emrgen/authbase/pkg/permission"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/emrgen/authbase/x"
	"github.com/emrgen/authbase/x/mail"
	"github.com/emrgen/authbase/x/oauth"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// NewAuthService creates a new AuthService
func NewAuthService(store store.Provider, keyProvider x.JWTSignerVerifierProvider, perm permission.AuthBasePermission, mailer mail.MailerProvider, cache *cache.Redis, verifier *x.StoreBasedUserVerifier) *AuthService {
	return &AuthService{store: store, keyProvider: keyProvider, perm: perm, mailer: mailer, cache: cache, verifier: verifier}
}

var _ v1.AuthServiceServer = new(AuthService)

// AuthService is a service that implements the AuthServiceServer interface
type AuthService struct {
	store       store.Provider
	mailer      mail.MailerProvider
	cache       *cache.Redis
	keyProvider x.JWTSignerVerifierProvider
	perm        permission.AuthBasePermission
	verifier    *x.StoreBasedUserVerifier
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
func (a *AuthService) LoginUsingIdp(ctx context.Context, request *v1.LoginUsingIdpRequest) (*v1.LoginUsingIdpResponse, error) {
	// get the provider details and redirect to the provider
	as, err := store.GetProjectStore(ctx, a.store)
	if err != nil {
		return nil, err
	}

	clientID, err := uuid.Parse(request.GetClientId())

	provider, err := as.GetOauthProviderByName(ctx, clientID, request.Provider)
	if err != nil {
		return nil, err
	}

	// this payload will not reach the client
	// TODO: grpc interceptors will handle the redirect with additional http.Cookie with the oauthstate
	// ref: https://github.com/douglasmakey/oauth2-example
	return &v1.LoginUsingIdpResponse{
		Provider: &v1.OAuthProvider{
			ClientId:     provider.Config.ClientID,
			ClientSecret: provider.Config.ClientSecret,
			RedirectUris: nil,
		},
	}, nil
}

// GetIdpToken gets the token from the identity provider code
func (a *AuthService) GetIdpToken(ctx context.Context, request *v1.GetIdpTokenRequest) (*v1.GetIdpTokenResponse, error) {
	as, err := store.GetProjectStore(ctx, a.store)
	if err != nil {
		return nil, err
	}
	// get the client id and secret from the request
	clientID, err := uuid.Parse(request.GetClientId())
	if err != nil {
		return nil, err
	}
	client, err := as.GetClientByID(ctx, clientID)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, errors.New("client not found")
	}

	// get the provider details
	provider, err := as.GetOauthProviderByName(ctx, clientID, request.Provider)
	if err != nil {
		return nil, err
	}

	oauthProvider, err := oauth.GetProvider(provider.Provider, oauth2.Config{
		ClientID:     provider.Config.ClientID,
		ClientSecret: provider.Config.ClientSecret,
	})
	if err != nil {
		return nil, err
	}

	// get the token from the provider
	token, err := oauthProvider.GetToken(ctx, request.Code)
	if err != nil {
		return nil, err
	}

	return &v1.GetIdpTokenResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    timestamppb.New(time.Now().Add(time.Duration(token.ExpiresIn))),
		IssuedAt:     timestamppb.New(time.Now()),
	}, nil
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
	hashedPassword := x.HashPassword(password, salt)

	user := &model.Account{
		ID:           uuid.New().String(),
		ProjectID:    orgID.String(),
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Salt:         salt,
		Verified:     false,
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
// TODO: /admin/login should we separate the admin login url from the other user logins?
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

	if client == nil {
		return nil, errors.New("client not found")
	}

	poolID := uuid.MustParse(client.PoolID)

	account, err := as.GetAccountByEmail(ctx, poolID, email)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, errors.New("account not found")
	}

	// disabled accounts should not be able to login
	if account.Disabled {
		// TODO: should we return a different error code, may be permission denied, or message?
		return nil, errors.New("account is disabled")
	}

	ok := x.CompareHashAndPassword(password, account.Salt, account.PasswordHash)
	if !ok {
		return nil, errors.New("incorrect password")
	}

	// get the account scopes from the group memberships
	accountID := uuid.MustParse(account.ID)
	memberships, err := as.ListGroupMemberByAccount(ctx, accountID)
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

	signer, err := a.keyProvider.GetSigner(poolID.String())
	if err != nil {
		return nil, err
	}

	// generate tokens for the account
	jti := uuid.New().String() // unique id for the token
	token, err := x.GenerateJWTToken(&x.Claims{
		Username:  account.Username,
		Email:     account.Email,
		ProjectID: account.ProjectID,
		PoolID:    account.PoolID,
		AccountID: account.ID,
		Audience:  "", // TODO: the target website or app that will use the token
		Jti:       jti,
		ExpireAt:  time.Now().Add(x.AccessTokenDuration),
		IssuedAt:  time.Now(),
		Provider:  "authbase", // TODO: what should this be?
		Scopes:    roleNames,  // internal roles
		Roles:     roleNames,
	}, signer)
	if err != nil {
		return nil, err
	}

	refreshExpireAt := time.Now().Add(x.RefreshTokenDuration)
	// save refresh token to cache, it will be used to validate the refresh token request
	// on cache miss, it will check the provider for the refresh token
	err = a.cache.Set(jti, token.RefreshToken, x.RefreshTokenDuration)
	if err != nil {
		return nil, err
	}

	// save the token to the provider
	// TODO: save the token to the provider in encrypted form
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
		// use the jti as the session id
		// this will allow multiple sessions for a account at the same time
		// TODO: should we limit the number of sessions?
		err = as.CreateSession(ctx, &model.Session{
			ID:        jti,
			PoolID:    account.PoolID,
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

// Logout logs out a user by deleting the session from the provider, and the refresh tokens from the cache
func (a *AuthService) Logout(ctx context.Context, request *v1.LogoutRequest) (*v1.LogoutResponse, error) {
	poolID, err := x.GetAuthbasePoolID(ctx)
	if err != nil {
		return nil, err
	}
	// check if the token is still valid
	accessToken, err := x.TokenFromHeader(ctx, "Bearer")
	if err != nil {
		return nil, err
	}
	verifier, err := a.keyProvider.GetVerifier(poolID.String())
	if err != nil {
		return nil, err
	}

	claims, err := x.VerifyJWTToken(accessToken, verifier)
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

	// delete the refresh token from the provider
	err = as.DeleteRefreshToken(ctx, claims.Jti)
	if err != nil {
		return nil, err
	}

	// delete the session from the provider
	err = as.DeleteSession(ctx, jti)
	if err != nil {
		return nil, err
	}

	return &v1.LogoutResponse{Message: "logged out"}, nil
}

// ForgotPassword sends a password reset link to the user's email or phone number to reset their password
func (a *AuthService) ForgotPassword(ctx context.Context, request *v1.ForgotPasswordRequest) (*v1.ForgotPasswordResponse, error) {
	poolID, err := x.GetAuthbaseProjectID(ctx)
	if err != nil {
		return nil, err
	}

	as, err := store.GetProjectStore(ctx, a.store)
	if err != nil {
		return nil, err
	}

	mailer := a.mailer.Provide(poolID)
	email := request.GetEmail()

	account, err := as.GetAccountByEmail(ctx, poolID, email)
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
		hashedPassword := x.HashPassword(password, salt)
		account.PasswordHash = string(hashedPassword)
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
		hashedPassword := x.HashPassword(request.GetNewPassword(), salt)
		account.PasswordHash = string(hashedPassword)
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

	oldClaims, err := x.GetTokenClaims(refreshToken)
	if err != nil {
		return nil, err
	}

	verifier, err := a.keyProvider.GetVerifier(oldClaims.ClientID)

	claims, err := x.VerifyJWTToken(refreshToken, verifier)
	if err != nil {
		return nil, err
	}

	// check if the token is in the cache
	tokenStr, err := a.cache.Get(claims.Jti)
	if err != nil {
		// if no value in cache check the provider
		return nil, err
	}

	as, err := store.GetProjectStore(ctx, a.store)
	if err != nil {
		return nil, err
	}

	if tokenStr != "" {
		accountID = claims.AccountID
		projectID = claims.ProjectID
		foundToken = true
	}

	if !foundToken {
		// check the provider
		token, err := as.GetRefreshTokenByID(ctx, claims.Jti)
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
	claims = &x.Claims{
		ProjectID: projectID,
		PoolID:    claims.PoolID,
		AccountID: user.ID,
		Username:  claims.Username,
		Email:     claims.Email,
		Audience:  claims.Audience,
		Jti:       jti,
		ExpireAt:  time.Now().Add(15 * time.Minute),
		IssuedAt:  time.Now(),
		Provider:  "authbase",
		Scopes:    claims.Scopes,
	}

	signer, err := a.keyProvider.GetSigner(claims.PoolID)
	if err != nil {
		return nil, err
	}

	token, err := x.GenerateJWTToken(claims, signer)
	if err != nil {
		return nil, err
	}
	err = a.cache.Set(claims.Jti, token.RefreshToken, x.ScheduleRefreshTokenExpiry)
	if err != nil {
		return nil, err
	}

	// TODO: should intercept the response and delete old token from cookie
	return &v1.RefreshResponse{
		Tokens: &v1.Tokens{
			AccessToken:      token.AccessToken,
			RefreshToken:     refreshToken,
			ExpiresAt:        timestamppb.New(token.ExpireAt),
			IssuedAt:         timestamppb.New(token.IssuedAt),
			RefreshExpiresAt: timestamppb.New(token.ExpireAt),
		},
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
		// TODO: use redis with TTL to store the verification code, for now we will use the provider (database)
		code, err := tx.GetVerificationCode(ctx, request.GetCode())
		if err != nil {
			return err
		}

		if code.ExpiresAt.Before(time.Now()) {
			return status.Errorf(codes.InvalidArgument, "verification code expired")
		}

		account, err := tx.GetAccountByID(ctx, uuid.MustParse(code.AccountID))
		if err != nil {
			return err
		}

		account.Verified = true
		account.VerifiedAt = time.Now().UTC()

		err = tx.UpdateAccount(ctx, account)
		if err != nil {
			return err
		}

		err = tx.DeleteVerificationCode(ctx, request.GetCode())
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// should redirect to the login page, let the UI handle the redirection
	return &v1.VerifyEmailResponse{Message: "email verified"}, nil
}
