package x

import (
	"context"
	"errors"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
)

func GetUserID(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value("userID").(uuid.UUID)
	if !ok {
		return uuid.UUID{}, ErrUserNotFoundInContext
	}

	return userID, nil
}

func GetOrganizationID(ctx context.Context) (uuid.UUID, error) {
	organizationID, ok := ctx.Value("organizationID").(uuid.UUID)
	if !ok {
		return uuid.UUID{}, ErrOrganizationNotFoundInContext
	}

	return organizationID, nil
}

func VerifyUserInterceptor(verifier UserVerifier) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// get url path from metadata

		logrus.Info("VerifyUserInterceptor")

		switch info.FullMethod {
		case v1.AuthService_Login_FullMethodName:
			request := req.(*v1.LoginRequest)
			orgID, err := uuid.Parse(request.GetOrganizationId())
			if err != nil {
				return nil, err
			}
			user, err := verifier.VerifyEmailPassword(ctx, orgID, request.Email, request.Password)
			if err != nil {
				return nil, err
			}

			ctx = context.WithValue(ctx, "userID", uuid.MustParse(user.ID))
			ctx = context.WithValue(ctx, "organizationID", uuid.MustParse(request.GetOrganizationId()))
		case v1.TokenService_CreateToken_FullMethodName:
			request := req.(*v1.CreateTokenRequest)
			orgID, err := uuid.Parse(request.GetOrganizationId())
			if err != nil {
				return nil, err
			}
			user, err := verifier.VerifyEmailPassword(ctx, orgID, request.Email, request.Password)
			if err != nil {
				return nil, err
			}

			ctx = context.WithValue(ctx, "userID", uuid.MustParse(user.ID))
			ctx = context.WithValue(ctx, "organizationID", uuid.MustParse(request.GetOrganizationId()))

		default:
			token, err := TokenFromHeader(ctx, "Bearer")
			if err != nil {
				return nil, err
			}
			claims, err := VerifyJWTToken(token)
			if err != nil {
				return nil, err
			}

			ctx = context.WithValue(ctx, "userID", uuid.MustParse(claims.UserID))
			ctx = context.WithValue(ctx, "organizationID", uuid.MustParse(claims.OrganizationID))
		}

		return handler(ctx, req)
	}
}

func TokenFromHeader(ctx context.Context, expectedScheme string) (string, error) {
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
