package service

import (
	"context"
	"errors"
	"strings"
	"time"

	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/cache"
	"github.com/emrgen/authbase/pkg/model"
	"github.com/emrgen/authbase/pkg/permission"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/emrgen/authbase/x"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// defaultAccessKeyExpireIn is the default expire time for a refresh token
	defaultAccessKeyExpireIn = time.Hour * 24 * 60 // 60 days
)

// NewAccessKeyService creates new an offline token service
func NewAccessKeyService(perm permission.AuthBasePermission, store store.Provider, cache *cache.Redis, keyProvider x.JWTSignerVerifierProvider, verifier x.TokenVerifier) v1.AccessKeyServiceServer {
	return &AccessKeyService{
		perm:        perm,
		store:       store,
		cache:       cache,
		verifier:    verifier,
		keyProvider: keyProvider,
	}
}

var _ v1.AccessKeyServiceServer = new(AccessKeyService)

// AccessKeyService is a service for token. AccessKey is an offline token.
type AccessKeyService struct {
	perm        permission.AuthBasePermission
	store       store.Provider
	cache       *cache.Redis
	verifier    x.TokenVerifier
	keyProvider x.JWTSignerVerifierProvider
	v1.UnimplementedAccessKeyServiceServer
}

// CreateAccessKey creates an offline access key
// 1. user authentication is already done by the middleware
// 2. get project store
// 3. check if the user includes the scopes
func (t *AccessKeyService) CreateAccessKey(ctx context.Context, request *v1.CreateAccessKeyRequest) (*v1.CreateAccessKeyResponse, error) {
	as, err := store.GetProjectStore(ctx, t.store)
	if err != nil {
		return nil, err
	}

	poolID, err := x.GetAuthbasePoolID(ctx)
	if err != nil {
		return nil, err
	}

	// TODO: check if the user has the permission to create the access key for the pool

	accountID, err := x.GetAuthbaseAccountID(ctx)
	if err != nil {
		return nil, err
	}

	expireAfter := defaultAccessKeyExpireIn
	if request.GetExpiresIn() != 0 {
		// expiry time can not be more than 120 days
		// TODO: this should be configurable in the future
		if request.GetExpiresIn() > 60*24*60*60 {
			return nil, errors.New("expires_in can not be more than 60 days")
		}
		expireAfter = time.Second * time.Duration(request.GetExpiresIn())
	}

	//perm, err := as.GetProjectMemberByID(ctx, clientID, accountID)
	//if err != nil {
	//	return nil, err
	//}

	// TODO: check if the user has the permission to create the access key for the pool

	// custom permissions from the downstream service
	scopes := request.GetScopes()

	expireAt := time.Now().Add(expireAfter)
	token := x.NewAccessKey()

	project, err := as.GetPoolByID(ctx, poolID)
	if err != nil {
		return nil, err
	}

	// create a new token
	accessKey := &model.AccessKey{
		ID:        token.ID.String(),
		AccountID: accountID.String(),
		PoolID:    poolID.String(),
		ProjectID: project.ID,
		Name:      request.GetName(),
		Token:     token.Value,
		Scopes:    strings.Join(scopes, ","),
		ExpireAt:  expireAt,
	}

	groupIDs := make([]string, 0)
	for _, group := range request.GetGroups() {
		groupIDs = append(groupIDs, group.GetId())
	}
	groupMembers := make([]*model.GroupMemberAccessKey, 0)
	for _, groupID := range groupIDs {
		groupMember := &model.GroupMemberAccessKey{
			GroupID:     groupID,
			AccessKeyID: accessKey.ID,
		}
		groupMembers = append(groupMembers, groupMember)
	}

	// save the token into the database
	err = as.Transaction(func(tx store.AuthBaseStore) error {
		err = t.cache.Set(token.ID.String(), accessKey.Token, defaultAccessKeyExpireIn)
		if err != nil {
			return err
		}

		err = tx.CreateAccessKey(ctx, accessKey)
		if err != nil {
			return err
		}

		if len(groupMembers) == 0 {
			// create group membership
			err := tx.CreateGroupMemberAccessKey(ctx, groupMembers)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.CreateAccessKeyResponse{
		Token: &v1.AccessKey{
			Id:        accessKey.ID,
			AccessKey: token.String(),
			ProjectId: accessKey.ProjectID,
			CreatedAt: timestamppb.New(time.Now()),
			ExpiresAt: timestamppb.New(expireAt),
		},
	}, nil
}

// GetAccessKey gets a token by id
func (t *AccessKeyService) GetAccessKey(ctx context.Context, request *v1.GetAccessKeyRequest) (*v1.GetAccessKeyResponse, error) {
	as, err := store.GetProjectStore(ctx, t.store)
	if err != nil {
		return nil, err
	}
	id, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, err
	}

	token, err := as.GetAccessKeyByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = t.perm.CheckProjectPermission(ctx, uuid.MustParse(token.ProjectID), "read")
	if err != nil {
		return nil, err
	}

	return &v1.GetAccessKeyResponse{
		Token: &v1.AccessKey{
			Id:        token.ID,
			Name:      token.Name,
			ProjectId: token.ProjectID,
			AccessKey: token.Token,
			CreatedAt: timestamppb.New(token.CreatedAt),
			ExpiresAt: timestamppb.New(token.ExpireAt),
		},
	}, nil
}

// ListAccessKeys lists offline tokens by project id and user id
func (t *AccessKeyService) ListAccessKeys(ctx context.Context, request *v1.ListAccessKeysRequest) (*v1.ListAccessKeysResponse, error) {
	as, err := store.GetProjectStore(ctx, t.store)
	if err != nil {
		return nil, err
	}

	projectID, err := x.GetAuthbaseProjectID(ctx)
	if err != nil {
		return nil, err
	}
	if request.ProjectId != nil {
		projectID = uuid.MustParse(request.GetProjectId())
	}

	err = t.perm.CheckProjectPermission(ctx, projectID, "read")
	if err != nil {
		return nil, err
	}

	accountID, err := x.GetAuthbaseAccountID(ctx)
	if err != nil {
		return nil, err
	}
	if request.AccountId != nil {
		accountID = uuid.MustParse(request.GetAccountId())
	}

	page := &v1.Page{
		Page: 0,
		Size: 20,
	}
	if request.Page != nil {
		page = request.GetPage()
	}
	size := max(page.Size, 20)

	accessKeys, total, err := as.ListAccountAccessKeys(ctx, projectID, accountID, int(page.Page), int(size))
	if err != nil {
		return nil, err
	}

	var keys []*v1.AccessKey
	for _, token := range accessKeys {
		keys = append(keys, &v1.AccessKey{
			Id:        token.ID,
			AccessKey: token.Token,
			AccountId: accountID.String(),
			ProjectId: projectID.String(),
			PoolId:    token.PoolID,
		})
	}

	return &v1.ListAccessKeysResponse{
		Tokens: keys,
		Meta: &v1.Meta{
			Total: int32(total),
			Page:  page.Page,
			Size:  size,
		},
	}, nil
}

// DeleteAccessKey deletes a offline token by id
func (t *AccessKeyService) DeleteAccessKey(ctx context.Context, request *v1.DeleteAccessKeyRequest) (*v1.DeleteAccessKeyResponse, error) {
	as, err := store.GetProjectStore(ctx, t.store)
	if err != nil {
		return nil, err

	}
	id, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, err
	}

	token, err := as.GetAccessKeyByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = t.perm.CheckProjectPermission(ctx, uuid.MustParse(token.ProjectID), "write")
	if err != nil {
		return nil, err
	}

	err = as.DeleteAccessKey(ctx, id)
	if err != nil {
		return nil, err
	}

	logrus.Errorf("delete token %v", request)
	return &v1.DeleteAccessKeyResponse{}, nil
}

// GetAccessKeyAccount get the account from the token
func (t *AccessKeyService) GetAccessKeyAccount(ctx context.Context, request *v1.GetAccessKeyAccountRequest) (*v1.GetAccessKeyAccountResponse, error) {
	as, err := store.GetProjectStore(ctx, t.store)
	if err != nil {
		return nil, err
	}

	accountID, err := x.GetAuthbaseAccountID(ctx)
	if err != nil {
		return nil, err
	}

	account, err := as.GetAccountByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return &v1.GetAccessKeyAccountResponse{
		Account: &v1.Account{
			Id:          account.ID,
			Email:       account.Email,
			VisibleName: account.VisibleName,
			Member:      account.ProjectMember,
			ProjectId:   account.ProjectID,
			PoolId:      account.PoolID,
		},
	}, nil
}

// GetTokenFromAccessKey gets a token from the access key
func (t *AccessKeyService) GetTokenFromAccessKey(ctx context.Context, request *v1.GetTokenFromAccessKeyRequest) (*v1.GetTokenFromAccessKeyResponse, error) {
	accessKey := request.GetAccessKey()
	token, err := x.ParseAccessKey(accessKey)
	if !errors.Is(err, x.ErrInvalidToken) && err != nil {
		return nil, err
	}

	if token == nil {
		return nil, err
	}

	claims, err := t.verifier.VerifyAccessKey(ctx, token.ID, token.Value)
	if err != nil {
		return nil, err
	}

	claims.Jti = uuid.New().String()
	claims.IssuedAt = time.Now()
	claims.ExpireAt = time.Now().Add(x.AccessTokenDuration)

	signer, err := t.keyProvider.GetSigner(claims.PoolID)
	if err != nil {
		return nil, err
	}

	// generate the token from the claims
	jwtToken, err := x.GenerateJWTToken(claims, signer)
	if err != nil {
		return nil, err
	}

	return &v1.GetTokenFromAccessKeyResponse{
		AccessToken:  jwtToken.AccessToken,
		RefreshToken: jwtToken.RefreshToken,
		ExpiresAt:    timestamppb.New(jwtToken.ExpireAt),
	}, nil
}
