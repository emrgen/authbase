package x

import (
	"context"
	"errors"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

// AuthInterceptor authenticates the request using the provided verifier.
// on success, it sets the accountID and projectID and account permission in the context.
func AuthInterceptor(verifier TokenVerifier, keyProvider JWTSignerVerifierProvider, provider store.Provider) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logrus.Infof("authbase: interceptor method: %s", info.FullMethod)
		switch info.FullMethod {
		case
			v1.AdminAuthService_AdminLoginUsingPassword_FullMethodName,
			v1.AuthService_LoginUsingPassword_FullMethodName,
			v1.AuthService_Refresh_FullMethodName,
			v1.AccessKeyService_GetTokenFromAccessKey_FullMethodName,
			v1.TokenService_VerifyToken_FullMethodName:
			break
		case v1.AccessKeyService_CreateAccessKey_FullMethodName:
			logrus.Infof("authbase: interceptor create access key")
			token, err := tokenFromHeader(ctx, "Bearer")

			logrus.Info(token, err)
			// check if the token present in the header
			if err == nil && token != "" {
				accessKey, err := ParseAccessKey(token)
				if !errors.Is(err, ErrInvalidToken) && err != nil {
					logrus.Errorf("authbase: interceptor error parsing access key: %v", err)
					return nil, err
				}

				if accessKey != nil {
					ctx, _, err = verifyAccessKey(ctx, verifier, accessKey)
					if err != nil {
						return nil, err
					}
				} else {
					ctx, _, err = verifyJwtToken(ctx, keyProvider, token)
					if err == nil {
						return nil, err
					}
				}

				return handler(ctx, req)
			}

			// verify the client secret and password
			as, err := store.GetProjectStore(ctx, provider)
			if err != nil {
				return nil, err
			}
			request := req.(*v1.CreateAccessKeyRequest)
			clientID, err := uuid.Parse(request.GetClientId())
			if err != nil {
				return nil, err
			}
			client, err := as.GetClientByID(ctx, clientID)
			if err != nil {
				return nil, grpc.Errorf(codes.NotFound, "authbase: get client by id failed: %v", err)
			}

			// verify the client secret
			clientSecret := request.GetClientSecret()
			ok := CompareHashAndPassword(clientSecret, client.Salt, client.SecretHash)
			if !ok {
				return nil, errors.New("invalid client secret")
			}

			poolID, err := uuid.Parse(client.PoolID)
			if err != nil {
				return nil, err
			}

			ctx, _, err = verifyPassword(ctx, verifier, poolID, request.GetEmail(), request.GetPassword())
			if err != nil {
				return nil, err
			}
		default:
			if info.FullMethod == v1.AuthService_ChangePassword_FullMethodName {
				request := req.(*v1.ChangePasswordRequest)
				oldPassword := request.GetOldPassword()
				email := request.GetEmail()
				if oldPassword != "" && email != "" {
					poolID, err := uuid.Parse(request.GetPoolId())
					if err != nil {
						return nil, err
					}
					ctx, _, err = verifyPassword(ctx, verifier, poolID, email, oldPassword)
					if err != nil {
						return nil, err
					}
					return handler(ctx, req)
				}
			}

			// TODO: if http cookie is present use that
			// user Bearer token for authentication
			token, err := tokenFromHeader(ctx, "Bearer")

			accessKey, err := ParseAccessKey(token)
			if !errors.Is(err, ErrInvalidToken) && err != nil {
				logrus.Errorf("authbase: interceptor error parsing access key: %v", err)
				return nil, err
			}

			if accessKey != nil {
				ctx, _, err = verifyAccessKey(ctx, verifier, accessKey)
				if err != nil {
					return nil, err
				}
			} else {
				ctx, _, err = verifyJwtToken(ctx, keyProvider, token)
				if err != nil {
					return nil, err
				}
			}
		}

		return handler(ctx, req)
	}
}

func verifyPassword(ctx context.Context, verifier TokenVerifier, poolID uuid.UUID, email, password string) (context.Context, *Claims, error) {
	user, err := verifier.VerifyEmailPassword(ctx, poolID, email, password)
	if err != nil {
		return ctx, nil, status.Error(codes.PermissionDenied, "old password is invalid")
	}

	ctx = context.WithValue(ctx, AccountIDKey, uuid.MustParse(user.ID))
	ctx = context.WithValue(ctx, ProjectIDKey, uuid.MustParse(user.ProjectID))
	ctx = context.WithValue(ctx, PoolIDKey, poolID)

	return ctx, nil, nil
}

func verifyAccessKey(ctx context.Context, verifier TokenVerifier, accessKey *AccessKey) (context.Context, *Claims, error) {
	if accessKey != nil {
		claims, err := verifier.VerifyAccessKey(ctx, accessKey.ID, accessKey.Value)
		if err != nil {
			return ctx, nil, err
		}

		ctx = context.WithValue(ctx, AccountIDKey, uuid.MustParse(claims.AccountID))
		ctx = context.WithValue(ctx, ProjectIDKey, uuid.MustParse(claims.ProjectID))
		ctx = context.WithValue(ctx, PoolIDKey, uuid.MustParse(claims.PoolID))
		ctx = context.WithValue(ctx, ScopesKey, claims.Scopes)
	}

	return ctx, nil, nil
}

func verifyJwtToken(ctx context.Context, keyProvider JWTSignerVerifierProvider, token string) (context.Context, *Claims, error) {
	claims, err := GetTokenClaims(token)
	if err != nil {
		return ctx, nil, err
	}

	verifier, err := keyProvider.GetVerifier(claims.PoolID)
	if err != nil {
		return ctx, nil, err
	}

	_, err = VerifyJWTToken(token, verifier)
	if err != nil {
		return ctx, nil, err
	}

	ctx = context.WithValue(ctx, AccountIDKey, uuid.MustParse(claims.AccountID))
	ctx = context.WithValue(ctx, ProjectIDKey, uuid.MustParse(claims.ProjectID))
	ctx = context.WithValue(ctx, PoolIDKey, uuid.MustParse(claims.PoolID))
	ctx = context.WithValue(ctx, ScopesKey, claims.Scopes)

	return ctx, claims, nil
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
