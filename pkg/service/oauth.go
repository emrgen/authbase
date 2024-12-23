package service

import (
	"context"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/cache"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/google/uuid"
)

var _ v1.OauthServiceServer = new(OauthService)

// OauthService is a service for oauth
type OauthService struct {
	store store.Provider
	cache *cache.Redis
	v1.UnimplementedOauthServiceServer
}

func NewOauthService(store store.Provider, cache *cache.Redis) *OauthService {
	return &OauthService{store: store, cache: cache}
}

// OAuthLogin authorizes a request and returns a response
func (o *OauthService) OAuthLogin(ctx context.Context, request *v1.OAuthLoginRequest) (*v1.OAuthLoginResponse, error) {
	// get the provider details and redirect to the provider
	as, err := store.GetProjectStore(ctx, o.store)
	if err != nil {
		return nil, err
	}

	orgID, err := uuid.Parse(request.GetOrganizationId())

	provider, err := as.GetOauthProviderByName(ctx, orgID, request.Provider)
	if err != nil {
		return nil, err
	}

	// this payload will not reach the client
	// grpc interceptors will handle the redirect with additional http.Cookie with the oauthstate
	// ref: https://github.com/douglasmakey/oauth2-example
	return &v1.OAuthLoginResponse{
		Provider: &v1.OAuthProvider{
			Name:         provider.Name,
			ClientId:     "",
			ClientSecret: "",
			RedirectUris: nil,
		},
	}, nil
}

// OAuthCallback handles the callback request after authorization
func (o *OauthService) OAuthCallback(ctx context.Context, request *v1.OAuthCallbackRequest) (*v1.OAuthCallbackResponse, error) {
	//cookie, err := x.GetOAuthState(ctx)
	//if err != nil {
	//	return nil, err
	//}
	// get the provider details and exchange the code for a token
	//TODO implement me
	panic("implement me")

	return &v1.OAuthCallbackResponse{}, nil
}

func (o *OauthService) Logout(ctx context.Context, request *v1.OauthLogoutRequest) (*v1.OauthLogoutResponse, error) {
	//TODO implement me
	panic("implement me")
}

// OAuthVerifyToken verifies the token for a user
func (o *OauthService) OAuthVerifyToken(ctx context.Context, request *v1.VerifyOAuthTokenRequest) (*v1.VerifyOAuthTokenResponse, error) {
	//TODO implement me
	panic("implement me")
}
