package x

import (
	"context"
	"errors"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("metadata not found")
		}

		switch info.FullMethod {
		case v1.TokenService_CreateToken_FullMethodName:
			request := req.(*v1.CreateTokenRequest)
			user, err := verifier.VerifyEmailPassword(ctx, request.Email, request.Password)
			if err != nil {
				return nil, err
			}

			ctx = context.WithValue(ctx, "userID", uuid.MustParse(user.ID))
			ctx = context.WithValue(ctx, "organizationID", uuid.MustParse(request.OrganizationId))
		}

		logrus.Info(md)

		return handler(ctx, req)
	}
}
