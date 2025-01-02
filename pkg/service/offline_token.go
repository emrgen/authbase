package service

import (
	"context"
	"errors"
	gox "github.com/emrgen/gopack/x"
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
	// defaultRefreshTokenExpireIn is the default expire time for a refresh token
	defaultRefreshTokenExpireIn = time.Hour * 24 * 60
)

var _ v1.OfflineTokenServiceServer = new(OfflineTokenService)

// OfflineTokenService is a service for token
type OfflineTokenService struct {
	perm  permission.AuthBasePermission
	store store.Provider
	cache *cache.Redis
	v1.UnimplementedOfflineTokenServiceServer
}

// NewOfflineTokenService creates new an offline token service
func NewOfflineTokenService(perm permission.AuthBasePermission, store store.Provider, cache *cache.Redis) *OfflineTokenService {
	return &OfflineTokenService{
		perm:  perm,
		store: store,
		cache: cache,
	}
}

// CreateToken creates an offline new token
// 1. user authentication is already done by the middleware
// 2. get project store
// 3. check if the user has the permission to create a token in the organization
func (t *OfflineTokenService) CreateToken(ctx context.Context, request *v1.CreateTokenRequest) (*v1.CreateTokenResponse, error) {
	as, err := store.GetProjectStore(ctx, t.store)
	if err != nil {
		return nil, err
	}

	orgID, err := uuid.Parse(request.GetOrganizationId())
	if err != nil {
		return nil, err
	}

	err = t.perm.CheckOrganizationPermission(ctx, orgID, "read")
	if err != nil {
		return nil, err
	}

	userID, err := gox.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	expireAfter := defaultRefreshTokenExpireIn
	if request.GetExpiresIn() != 0 {
		duration := time.Second * time.Duration(request.GetExpiresIn())
		expireAfter = duration
	}

	user, err := as.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	perm, err := as.GetPermissionByID(ctx, orgID, userID)
	if err != nil {
		return nil, err
	}

	data := make(map[string]string)
	if request.Data != nil {
		data = request.GetData()
	}

	jti := uuid.New().String()
	token, err := x.GenerateJWTToken(x.Claims{
		Username:       user.Username,
		Email:          user.Email,
		OrganizationID: orgID.String(),
		UserID:         userID.String(),
		Permission:     perm.Permission,
		Audience:       "",
		Jti:            jti,
		ExpireAt:       time.Now().Add(expireAfter),
		IssuedAt:       time.Now(),
		Provider:       "authbase",
		Data:           data,
	})
	if err != nil {
		return nil, err
	}

	// create a new token
	tokenModel := &model.Token{
		ID:             uuid.New().String(),
		OrganizationID: request.GetOrganizationId(),
		Name:           request.GetName(),
		Token:          token.AccessToken,
		ExpireAt:       token.ExpireAt,
	}

	// save the token into the database
	err = as.Transaction(func(tx store.AuthBaseStore) error {
		err = t.cache.Set(jti, tokenModel.OrganizationID, defaultRefreshTokenExpireIn)
		if err != nil {
			return err
		}

		tokenModel.UserID = user.ID
		err = tx.CreateToken(ctx, tokenModel)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.CreateTokenResponse{
		Id:    tokenModel.ID,
		Token: tokenModel.Token,
	}, nil
}

// GetToken gets a token by id
func (t *OfflineTokenService) GetToken(ctx context.Context, request *v1.GetTokenRequest) (*v1.GetTokenResponse, error) {
	as, err := store.GetProjectStore(ctx, t.store)
	if err != nil {
		return nil, err
	}
	id, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, err
	}

	token, err := as.GetTokenByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = t.perm.CheckOrganizationPermission(ctx, uuid.MustParse(token.OrganizationID), "read")
	if err != nil {
		return nil, err
	}

	return &v1.GetTokenResponse{
		Token: &v1.Token{
			Id:             token.ID,
			Name:           token.Name,
			OrganizationId: token.OrganizationID,
			UserId:         token.UserID,
			CreatedAt:      timestamppb.New(token.CreatedAt),
			ExpiresAt:      timestamppb.New(token.ExpireAt),
		},
	}, nil
}

// ListTokens lists offline tokens by organization id and user id
func (t *OfflineTokenService) ListTokens(ctx context.Context, request *v1.ListTokensRequest) (*v1.ListTokensResponse, error) {
	as, err := store.GetProjectStore(ctx, t.store)
	if err != nil {
		return nil, err
	}

	orgID, err := uuid.Parse(request.GetOrganizationId())
	if err != nil {
		return nil, err
	}

	err = t.perm.CheckOrganizationPermission(ctx, orgID, "read")
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(request.GetUserId())
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

	tokens, total, err := as.ListUserTokens(ctx, orgID, userID, int(page.Page), int(size))
	if err != nil {
		return nil, err
	}

	var tokenProtos []*v1.Token
	for _, token := range tokens {
		tokenProtos = append(tokenProtos, &v1.Token{
			Id:    token.ID,
			Token: token.Token,
		})
	}

	return &v1.ListTokensResponse{
		Tokens: tokenProtos,
		Meta: &v1.Meta{
			Total: int32(total),
			Page:  page.Page,
			Size:  size,
		},
	}, nil
}

// DeleteToken deletes a offline token by id
func (t *OfflineTokenService) DeleteToken(ctx context.Context, request *v1.DeleteTokenRequest) (*v1.DeleteTokenResponse, error) {
	as, err := store.GetProjectStore(ctx, t.store)
	if err != nil {
		return nil, err

	}
	id, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, err
	}

	token, err := as.GetTokenByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = t.perm.CheckOrganizationPermission(ctx, uuid.MustParse(token.OrganizationID), "write")
	if err != nil {
		return nil, err
	}

	err = as.DeleteToken(ctx, id)
	if err != nil {
		return nil, err
	}

	logrus.Errorf("delete token %v", request)
	return &v1.DeleteTokenResponse{}, nil
}

// VerifyToken verifies a token and returns the organization id and user id
// no need to check the permission here
func (t *OfflineTokenService) VerifyOfflineToken(ctx context.Context, request *v1.OfflineTokenVerifyRequest) (*v1.OfflineTokenVerifyResponse, error) {
	token := request.GetToken()
	jwt, err := x.VerifyJWTToken(token)
	if err != nil {
		return nil, err
	}

	if jwt.ExpireAt.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return &v1.OfflineTokenVerifyResponse{
		OrganizationId: jwt.OrganizationID,
		UserId:         jwt.UserID,
	}, nil
}
