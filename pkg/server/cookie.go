package server

import (
	"context"
	"github.com/emrgen/authbase/pkg/cache"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"net/http"
)

// InjectCookie injects a cookie into the response based on the message type
func InjectCookie(store CookieStore) func(ctx context.Context, w http.ResponseWriter, m proto.Message) error {
	return func(ctx context.Context, w http.ResponseWriter, m proto.Message) error {
		//m = m.ProtoReflect().Interface()
		//switch m.(type) {
		//case *v1.LoginResponse:
		//	loginResponse := m.(*v1.LoginResponse)
		//	// marshal the token
		//	token := loginResponse.Token.String()
		//	logrus.Info("Token: ", token)
		//	cookie := http.Cookie{
		//		Name:     "session",
		//		HttpOnly: true,
		//		Value:    token,
		//		SameSite: http.SameSiteStrictMode, // Strict mode is the most secure option
		//	}
		//	http.SetCookie(w, &cookie)
		//case *v1.OAuthLoginResponse:
		//	state := x.GenerateCode()
		//	// create the redirection URL
		//	cookie := http.Cookie{
		//		Name:     "oauthstate",
		//		HttpOnly: true,
		//		Value:    state,
		//		SameSite: http.SameSiteStrictMode, // Strict mode is the most secure option
		//	}
		//	http.SetCookie(w, &cookie)
		//
		//	res := m.(*v1.OAuthLoginResponse)
		//	provider := res.Provider
		//
		//	var googleOauthConfig = &oauth2.Config{
		//		RedirectURL:  res.CallbackUrl,
		//		ClientID:     provider.ClientId,
		//		ClientSecret: provider.ClientSecret,
		//		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		//		Endpoint:     google.Endpoint,
		//	}
		//
		//	err := store.Set(ctx, state)
		//	if err != nil {
		//		return err
		//	}
		//
		//	redirectURL := googleOauthConfig.AuthCodeURL(state)
		//	w.Header().Set("Location", redirectURL)
		//	w.WriteHeader(http.StatusFound)
		//}

		return nil
	}
}

// ExtractCookie extracts a cookie from the request and check if the cookie is valid
func ExtractCookie(store CookieStore) func(ctx context.Context, r *http.Request) metadata.MD {
	return func(ctx context.Context, r *http.Request) metadata.MD {
		cookie, err := r.Cookie("oauthstate")
		if err != nil {
			return metadata.Pairs()
		}

		exists, err := store.Exists(ctx, cookie.Value)
		if err != nil {
			return metadata.Pairs("error", err.Error())
		}

		if !exists {
			return metadata.Pairs("error", "invalid oauthstate")
		}

		return metadata.Pairs("oauthstate", cookie.Value)
	}
}

// CookieStore stores sessions using secure cookies.
type CookieStore struct {
	redis *cache.Redis
}

func NewCookieStore(redis *cache.Redis) CookieStore {
	return CookieStore{redis: redis}
}

func (s *CookieStore) Exists(ctx context.Context, state string) (bool, error) {
	return s.redis.SExists("oauthstate", state)
}

func (s *CookieStore) Set(ctx context.Context, state string) error {
	return s.redis.SAdd("oauthstate", state)
}
