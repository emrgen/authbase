package service

import (
	"context"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/cache"
	"github.com/emrgen/authbase/pkg/store"
)

var _ v1.OauthServiceServer = new(OauthService)

// OauthService is a service for oauth
type OauthService struct {
	store store.AuthBaseStore
	cache *cache.Redis
	v1.UnimplementedOauthServiceServer
}

func NewOauthService(store store.AuthBaseStore, cache *cache.Redis) *OauthService {
	return &OauthService{store: store, cache: cache}
}

// Authorize authorizes a request and returns a response
func (o *OauthService) Authorize(ctx context.Context, request *v1.AuthorizeRequest) (*v1.AuthorizeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o *OauthService) Logout(ctx context.Context, request *v1.OauthLogoutRequest) (*v1.OauthLogoutResponse, error) {
	//TODO implement me
	panic("implement me")
}

// Callback handles the callback request after authorization
func (o *OauthService) Callback(ctx context.Context, request *v1.CallbackRequest) (*v1.CallbackResponse, error) {
	//TODO implement me
	panic("implement me")
}

// Token returns a token
func (o *OauthService) Token(ctx context.Context, request *v1.TokenRequest) (*v1.TokenResponse, error) {
	//TODO implement me
	panic("implement me")
}
