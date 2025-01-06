package service

//
//var _ v1.OAuth2ServiceServer = new(OAuth2Service)
//
//// OauthService is a service for oauth
//type OAuth2Service struct {
//	store store.Provider
//	cache *cache.Redis
//	v1.UnimplementedOAuth2ServiceServer
//}
//
//func NewOauthService(store store.Provider, cache *cache.Redis) *OAuth2Service {
//	return &OAuth2Service{store: store, cache: cache}
//}
//
//// OAuthLogin authorizes a request and returns a response
//func (o *OAuth2Service) OAuth2Login(ctx context.Context, request *v1.OAuthLoginRequest) (*v1.OAuthLoginResponse, error) {
//	// get the provider details and redirect to the provider
//	as, err := store.GetProjectStore(ctx, o.store)
//	if err != nil {
//		return nil, err
//	}
//
//	orgID, err := uuid.Parse(request.GetProjectId())
//
//	provider, err := as.GetOauthProviderByName(ctx, orgID, request.Provider)
//	if err != nil {
//		return nil, err
//	}
//
//	// this payload will not reach the client
//	// grpc interceptors will handle the redirect with additional http.Cookie with the oauthstate
//	// ref: https://github.com/douglasmakey/oauth2-example
//	return &v1.OAuthLoginResponse{
//		Provider: &v1.OAuthProvider{
//			Name:         provider.Name,
//			ClientId:     provider.Config.ClientID,
//			ClientSecret: provider.Config.ClientSecret,
//			RedirectUris: nil,
//		},
//	}, nil
//}
//
//// OAuthCallback handles the callback request after authorization
//func (o *OAuth2Service) OAuth2Callback(ctx context.Context, request *v1.OAuthCallbackRequest) (*v1.OAuthCallbackResponse, error) {
//	//cookie, err := x.GetOAuthState(ctx)
//	//if err != nil {
//	//	return nil, err
//	//}
//	// get the provider details and exchange the code for a token
//	//TODO implement me
//	panic("implement me")
//
//	return &v1.OAuthCallbackResponse{}, nil
//}
//
//// Logout logs out a user, invalidating the token and removing the session
//func (o *OAuth2Service) Logout(ctx context.Context, request *v1.OAuth2LogoutRequest) (*v1.OauthLogoutResponse, error) {
//	//TODO implement me
//	panic("implement me")
//}
