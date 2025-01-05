package service

import (
	"context"
	v1 "github.com/emrgen/authbase/apis/v1"
	gopackv1 "github.com/emrgen/gopack/apis/v1"
	"github.com/golang-jwt/jwt/v5"
)

var _ gopackv1.TokenServiceServer = (*TokenService)(nil)

type TokenService struct {
	offlineTokenService v1.OfflineTokenServiceServer
	oauthService        v1.OauthServiceServer
	gopackv1.UnimplementedTokenServiceServer
}

func NewTokenService(offlineTokenService v1.OfflineTokenServiceServer, oauthService v1.OauthServiceServer) *TokenService {
	return &TokenService{
		offlineTokenService: offlineTokenService,
		oauthService:        oauthService,
	}
}

// VerifyToken verifies the token and returns the user id and project id
func (t TokenService) VerifyToken(ctx context.Context, request *gopackv1.VerifyTokenRequest) (*gopackv1.VerifyTokenResponse, error) {
	bearerToken := request.GetToken()
	token, _, err := jwt.NewParser().ParseUnverified(bearerToken, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	claim := token.Claims.(jwt.MapClaims)
	provider, ok := claim["provider"].(string)
	if !ok {
		return nil, err
	}

	switch provider {
	case "authbase":
		res, err := t.offlineTokenService.VerifyOfflineToken(ctx, &v1.OfflineTokenVerifyRequest{Token: bearerToken})
		if err != nil {
			return nil, err
		}

		return &gopackv1.VerifyTokenResponse{Valid: true, UserId: res.GetUserId(), ProjectId: res.GetProjectId()}, nil
	case "google":
		res, err := t.oauthService.OAuthVerifyToken(ctx, &v1.VerifyOAuthTokenRequest{Token: bearerToken})
		if err != nil {
			return nil, err
		}
		return &gopackv1.VerifyTokenResponse{Valid: true, UserId: res.GetUserId(), ProjectId: res.GetProjectId()}, nil
	}

	return &gopackv1.VerifyTokenResponse{Valid: false}, nil
}
