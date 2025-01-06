package service

import (
	"context"
	gox "github.com/emrgen/gopack/x"
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
func NewAccessKeyService(perm permission.AuthBasePermission, store store.Provider, cache *cache.Redis) v1.AccessKeyServiceServer {
	return &AccessKeyService{
		perm:  perm,
		store: store,
		cache: cache,
	}
}

var _ v1.AccessKeyServiceServer = new(AccessKeyService)

// AccessKeyService is a service for token
type AccessKeyService struct {
	perm  permission.AuthBasePermission
	store store.Provider
	cache *cache.Redis
	v1.UnimplementedAccessKeyServiceServer
}

// CreateAccessKey creates an offline access key
// 1. user authentication is already done by the middleware
// 2. get project store
// 3. check if the user has the permission to create a token in the project
func (t *AccessKeyService) CreateAccessKey(ctx context.Context, request *v1.CreateAccessKeyRequest) (*v1.CreateAccessKeyResponse, error) {
	as, err := store.GetProjectStore(ctx, t.store)
	if err != nil {
		return nil, err
	}

	orgID, err := uuid.Parse(request.GetProjectId())
	if err != nil {
		return nil, err
	}

	err = t.perm.CheckProjectPermission(ctx, orgID, "write")
	if err != nil {
		return nil, err
	}

	userID, err := gox.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	expireAfter := defaultAccessKeyExpireIn
	if request.GetExpiresIn() != 0 {
		expireAfter = time.Second * time.Duration(request.GetExpiresIn())
	}

	//perm, err := as.GetProjectMemberByID(ctx, orgID, userID)
	//if err != nil {
	//	return nil, err
	//}

	// TODO: check if the user has the permission to include the scopes
	scopes := make([]string, 0)
	if request.Scopes != nil {
		scopes = request.GetScopes()
	}

	expireAt := time.Now().Add(expireAfter)
	token := x.NewAccessKey()

	// create a new token
	accessKey := &model.AccessKey{
		ID:        token.ID.String(),
		AccountID: userID.String(),
		ProjectID: request.GetProjectId(),
		Name:      request.GetName(),
		Token:     token.Value,
		Scopes:    strings.Join(scopes, ","),
		ExpireAt:  expireAt,
	}

	// save the token into the database
	err = as.Transaction(func(tx store.AuthBaseStore) error {
		err = t.cache.Set(token.ID.String(), accessKey.ProjectID, defaultAccessKeyExpireIn)
		if err != nil {
			return err
		}

		err = tx.CreateAccessKey(ctx, accessKey)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.CreateAccessKeyResponse{
		Id:    accessKey.ID,
		Token: token.String(),
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

	orgID, err := uuid.Parse(request.GetProjectId())
	if err != nil {
		return nil, err
	}

	err = t.perm.CheckProjectPermission(ctx, orgID, "read")
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(request.GetAccountId())
	if err != nil {
		return nil, err
	}

	page := &v1.Page{
		Page: 0,
		Size: 20,
	}
	if request.Page != nil {
		page = request.GetPage()
	}
	size := max(page.Size, 20)

	tokens, total, err := as.ListAccountAccessKeys(ctx, orgID, userID, int(page.Page), int(size))
	if err != nil {
		return nil, err
	}

	var tokenProtos []*v1.AccessKey
	for _, token := range tokens {
		tokenProtos = append(tokenProtos, &v1.AccessKey{
			Id:        token.ID,
			AccessKey: token.Token,
		})
	}

	return &v1.ListAccessKeysResponse{
		Tokens: tokenProtos,
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
		},
	}, nil
}
