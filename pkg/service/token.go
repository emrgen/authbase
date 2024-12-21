package service

import (
	"context"
	"errors"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/cache"
	"github.com/emrgen/authbase/pkg/model"
	"github.com/emrgen/authbase/pkg/permission"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/emrgen/authbase/x"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

const (
	defaultRefreshTokenExpireIn = time.Hour * 24 * 7
)

var _ v1.TokenServiceServer = new(TokenService)

// TokenService is a service for token
type TokenService struct {
	perm  permission.AuthBasePermission
	store store.Provider
	cache *cache.Redis
	v1.UnimplementedTokenServiceServer
}

// NewTokenService creates new an offline token service
func NewTokenService(store store.Provider, cache *cache.Redis) *TokenService {
	return &TokenService{
		store: store,
		cache: cache,
	}
}

// CreateToken creates an offline new token
func (t *TokenService) CreateToken(ctx context.Context, request *v1.CreateTokenRequest) (*v1.CreateTokenResponse, error) {
	as, err := store.GetProjectStore(ctx, t.store)
	if err != nil {
		return nil, err
	}

	orgID, err := uuid.Parse(request.GetOrganizationId())
	if err != nil {
		return nil, err
	}

	err = t.perm.CheckOrganizationPermission(ctx, orgID, "write")
	if err != nil {
		return nil, err
	}

	// create a new token
	token := &model.Token{
		ID:             uuid.New().String(),
		OrganizationID: request.GetOrganizationId(),
		Token:          x.GenerateToken(),
		Name:           request.GetName(),
	}

	logrus.Info("TokenService", token, request.GetOrganizationId())

	if request.GetExpiresIn() != 0 {
		duration := time.Second * time.Duration(request.GetExpiresIn())
		token.ExpireAt = time.Now().Add(duration)
	} else {
		token.ExpireAt = time.Now().Add(defaultRefreshTokenExpireIn)
	}

	// save the token into the database
	err = as.Transaction(func(tx store.AuthBaseStore) error {
		// check if the user exists on the database within the organization
		user, err := tx.GetUserByEmail(ctx, request.GetEmail())
		if err != nil {
			return err
		}

		// check if the user is in the organization
		if request.GetPassword() != user.Password {
			return errors.New("invalid password")
		}

		err = t.cache.Set(token.Token, token.OrganizationID, defaultRefreshTokenExpireIn)
		if err != nil {
			return err
		}

		token.UserID = user.ID
		err = tx.CreateToken(ctx, token)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.CreateTokenResponse{
		Id:    token.ID,
		Token: token.Token,
	}, nil
}

// GetToken gets a token by id
func (t *TokenService) GetToken(ctx context.Context, request *v1.GetTokenRequest) (*v1.GetTokenResponse, error) {
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
func (t *TokenService) ListTokens(ctx context.Context, request *v1.ListTokensRequest) (*v1.ListTokensResponse, error) {
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
func (t *TokenService) DeleteToken(ctx context.Context, request *v1.DeleteTokenRequest) (*v1.DeleteTokenResponse, error) {
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
func (t *TokenService) VerifyToken(ctx context.Context, request *v1.VerifyTokenRequest) (*v1.VerifyTokenResponse, error) {
	token := request.GetToken()
	jwt, err := x.VerifyJWTToken(token)
	if err != nil {
		return nil, err
	}

	return &v1.VerifyTokenResponse{
		OrganizationId: jwt.OrgID,
		UserId:         jwt.UserID,
	}, nil
}
