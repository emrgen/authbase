package server

import (
	"context"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"net/http"
)

// InjectCookie injects a cookie into the response
func InjectCookie(ctx context.Context, response http.ResponseWriter, m proto.Message) error {
	m = m.ProtoReflect().Interface()
	switch m.(type) {
	case *v1.LoginResponse:
		loginResponse := m.(*v1.LoginResponse)
		// marshal the token
		token := loginResponse.Token.String()
		logrus.Info("Token: ", token)
		cookie := http.Cookie{
			Name:     "session",
			HttpOnly: true,
			Value:    token,
			SameSite: http.SameSiteStrictMode, // Strict mode is the most secure option
		}
		http.SetCookie(response, &cookie)
	}

	return nil
}

func ExtractCookie(ctx context.Context, request *http.Request) (string, error) {
	cookie, err := request.Cookie("session")
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}
