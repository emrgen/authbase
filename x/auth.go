package x

import (
	"context"
	"errors"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
)

// AuthInterceptor authenticates the request using the provided verifier.
// on success, it sets the accountID and projectID and account permission in the context.
func AuthInterceptor(verifier TokenVerifier) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		switch info.FullMethod {
		case
			v1.AuthService_LoginUsingPassword_FullMethodName,
			v1.TokenService_VerifyToken_FullMethodName:
			break
		case v1.AccessKeyService_CreateAccessKey_FullMethodName:
			request := req.(*v1.CreateAccessKeyRequest)
			orgID, err := uuid.Parse(request.GetProjectId())
			if err != nil {
				return nil, err
			}
			user, err := verifier.VerifyEmailPassword(ctx, orgID, request.Email, request.Password)
			if err != nil {
				return nil, err
			}

			ctx = context.WithValue(ctx, AccountIDKey, uuid.MustParse(user.ID))
			ctx = context.WithValue(ctx, ProjectIDKey, uuid.MustParse(request.GetProjectId()))
		default:
			// TODO: if http cookie is present use that
			// user Bearer token for authentication
			token, err := tokenFromHeader(ctx, "Bearer")

			accessKey, err := ParseAccessKey(token)
			if !errors.Is(err, ErrInvalidToken) && err != nil {
				return nil, err
			}

			if accessKey != nil {
				claims, err := verifier.VerifyAccessKey(ctx, accessKey.ID, accessKey.Value)
				if err != nil {
					return nil, err
				}

				ctx = context.WithValue(ctx, AccountIDKey, uuid.MustParse(claims.AccountID))
				ctx = context.WithValue(ctx, ProjectIDKey, uuid.MustParse(claims.ProjectID))
				ctx = context.WithValue(ctx, PoolIDKey, uuid.MustParse(claims.PoolID))
				ctx = context.WithValue(ctx, ScopesKey, claims.Scopes)
			} else {
				if err != nil {
					return nil, err
				}
				claims, err := VerifyJWTToken(token)
				if err != nil {
					return nil, err
				}

				ctx = context.WithValue(ctx, AccountIDKey, uuid.MustParse(claims.AccountID))
				ctx = context.WithValue(ctx, ProjectIDKey, uuid.MustParse(claims.ProjectID))
				ctx = context.WithValue(ctx, PoolIDKey, uuid.MustParse(claims.PoolID))
				ctx = context.WithValue(ctx, ScopesKey, claims.Scopes)
			}
		}

		return handler(ctx, req)
	}
}

func tokenFromHeader(ctx context.Context, expectedScheme string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("metadata not found")
	}

	val, ok := md["authorization"]
	if !ok {
		return "", errors.New("no authorization header found")
	}

	if len(val) == 0 {
		return "", errors.New("no token found")
	}

	scheme, token, found := strings.Cut(val[0], " ")
	if !found {
		return "", errors.New("bad authorization string")
	}

	if !strings.EqualFold(scheme, expectedScheme) {
		return "", errors.New("request unauthenticated with " + expectedScheme)
	}

	return token, nil
}
