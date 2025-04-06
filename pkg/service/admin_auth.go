package service

import (
	"context"
	"errors"
	goset "github.com/deckarep/golang-set/v2"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/cache"
	"github.com/emrgen/authbase/pkg/config"
	"github.com/emrgen/authbase/pkg/model"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/emrgen/authbase/x"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type AdminAuthService struct {
	provider    store.Provider
	cfg         *config.AdminProjectConfig
	cache       *cache.Redis
	keyProvider x.JWTSignerVerifierProvider

	v1.UnimplementedAdminAuthServiceServer
}

func NewAdminAuthService(store store.Provider, cfg *config.AdminProjectConfig, keyProvider x.JWTSignerVerifierProvider, cache *cache.Redis,
) *AdminAuthService {
	return &AdminAuthService{
		provider:    store,
		cfg:         cfg,
		cache:       cache,
		keyProvider: keyProvider,
	}
}

var _ v1.AdminAuthServiceServer = (*AdminAuthService)(nil)

// AdminLoginUsingPassword handles the admin login using password
func (a AdminAuthService) AdminLoginUsingPassword(ctx context.Context, request *v1.AdminLoginUsingPasswordRequest) (*v1.AdminLoginUsingPasswordResponse, error) {
	as, err := store.GetProjectStore(ctx, a.provider)
	if err != nil {
		return nil, err
	}

	// TODO: this is very sensitive api
	// rate limit and ip block should be implemented here
	// more than 5 failed attempts should block the ip for 1 hour

	email := request.GetEmail()
	password := request.GetPassword()
	projectName := request.GetProjectName()

	// get default pool for the project
	project, err := as.GetProjectByName(ctx, projectName)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.New("project not found")
	}

	poolID := uuid.MustParse(project.PoolID)

	account, err := as.GetAccountByEmail(ctx, poolID, email)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, errors.New("admin account not found")
	}

	// validate the password
	ok := x.CompareHashAndPassword(password, account.Salt, account.PasswordHash)
	if !ok {
		return nil, errors.New("incorrect password")
	}

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

	// return the token
	return &v1.AdminLoginUsingPasswordResponse{
		Account: &v1.Account{
			Id:          account.ID,
			Username:    account.Username,
			Email:       account.Email,
			ProjectId:   account.ProjectID,
			PoolId:      account.PoolID,
			Member:      account.ProjectMember,
			VisibleName: account.VisibleName,
			CreatedAt:   timestamppb.New(account.CreatedAt),
			UpdatedAt:   timestamppb.New(account.UpdatedAt),
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
