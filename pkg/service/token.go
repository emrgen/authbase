package service

import (
	"context"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/x"
)

var _ v1.TokenServiceServer = (*TokenService)(nil)

type TokenService struct {
	verifier x.TokenVerifier
	v1.UnimplementedTokenServiceServer
}

func NewTokenService(verifier x.TokenVerifier) *TokenService {
	return &TokenService{
		verifier: verifier,
	}
}

// VerifyToken verifies the token and returns the user id and project id
func (t TokenService) VerifyToken(ctx context.Context, request *v1.VerifyTokenRequest) (*v1.VerifyTokenResponse, error) {
	yes := x.IsAccessKey(request.GetToken())
	if yes {
		key, err := x.ParseAccessKey(request.GetToken())
		if err != nil {
			return nil, err
		}

		accessKey, err := t.verifier.VerifyAccessKey(ctx, key.ID, key.Value)
		if err != nil {
			return nil, err
		}

		return &v1.VerifyTokenResponse{Valid: true, UserId: accessKey.AccountID, ProjectId: accessKey.ProjectID, PoolId: accessKey.PoolID}, nil
	}

	bearerToken := request.GetToken()
	claims, err := x.GetTokenClaims(bearerToken)
	if err != nil {
		return nil, err
	}

	switch claims.Provider {
	case "authbase":
		res, err := t.verifier.VerifyToken(ctx, bearerToken, claims.PoolID)
		if err != nil {
			return nil, err
		}

		return &v1.VerifyTokenResponse{Valid: true, UserId: res.AccountID, ProjectId: res.ProjectID, PoolId: res.ProjectID}, nil
	}

	return &v1.VerifyTokenResponse{Valid: false}, nil
}
