package x

import (
	"context"
	"errors"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	// ProjectPermissionKey is the key to store the project permission in the context
	ProjectPermissionKey = "authbase_project_permission"
	// ProjectIDKey is the key to store the project id in the context
	ProjectIDKey = "authbase_project_id"
	// PoolIDKey is the key to store the pool id in the context
	PoolIDKey = "authbase_pool_id"
	// AccountIDKey is the key to store the user id in the context
	AccountIDKey = "authbase_account_id"
	// ScopesKey is the key to store the scopes in the context
	ScopesKey = "authbase_scopes"
	// TokenMissingKey is the key to store the token in the context
	TokenMissingKey = "authbase_token_missing"
	RolesKey        = "authbase_roles"
)

type ProjectID interface {
	GetProjectId() string
}

func InjectPermissionInterceptor(member v1.ProjectMemberServiceClient) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		// check if the projectPermission are in the context
		_, err = GetAuthbaseProjectPermission(ctx)
		if err == nil {
			return nil, status.Errorf(codes.InvalidArgument, "missing project permission")
		}

		pid, ok := req.(ProjectID)
		if !ok {
			return nil, status.Errorf(codes.InvalidArgument, "missing project id")
		}

		projectID := pid.GetProjectId()
		if len(projectID) == 0 {
			return nil, status.Errorf(codes.InvalidArgument, "missing project id")
		}
		logrus.Infof("project id: %s", projectID)

		projectUUID, err := uuid.Parse(projectID)
		if err != nil {
			return nil, err
		}

		userID, ok := ctx.Value(AccountIDKey).(uuid.UUID)
		if !ok {
			return nil, status.Errorf(codes.InvalidArgument, "missing user id")
		}

		md, ok := metadata.FromIncomingContext(ctx)
		token := md.Get("authorization")
		md.Set("Authorization", token[0])
		ctx = metadata.NewOutgoingContext(ctx, md)

		// GetProjectMembership is a gRPC call to the project membership service
		res, err := member.GetProjectMember(ctx, &v1.GetProjectMemberRequest{
			ProjectId: projectUUID.String(),
			MemberId:  userID.String(),
		})
		if err != nil {
			return nil, err
		}

		ctx = context.WithValue(ctx, ProjectPermissionKey, res.ProjectMember.Permission)
		ctx = context.WithValue(ctx, ProjectIDKey, projectUUID)

		return handler(ctx, req)
	}
}

func VerifyTokenInterceptor(keyProvider VerifierProvider, accessKeyService v1.AccessKeyServiceClient) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		token, err := TokenFromHeader(ctx, "Bearer")
		if errors.Is(err, NoAuthHeaderError) {
			// No auth header, skip the interceptor
			ctx = context.WithValue(ctx, TokenMissingKey, true)
			return handler(ctx, req)
		}
		if err != nil {
			return nil, err
		}

		accessKey, err := ParseAccessKey(token)
		if !errors.Is(err, ErrInvalidToken) && err != nil {
			return nil, err
		}

		if accessKey != nil {
			res, err := accessKeyService.GetTokenFromAccessKey(ctx, &v1.GetTokenFromAccessKeyRequest{
				AccessKey: accessKey.String(),
			})
			if err != nil {
				return nil, err
			}

			claims, err := GetTokenClaims(res.GetAccessToken())
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

			claims, err := GetTokenClaims(token)
			if err != nil {
				return nil, err
			}

			verifier, err := keyProvider.GetVerifier(claims.PoolID)
			if err != nil {
				return nil, err
			}

			_, err = VerifyJWTToken(token, verifier)
			if err != nil {
				return nil, err
			}

			ctx = context.WithValue(ctx, AccountIDKey, uuid.MustParse(claims.AccountID))
			ctx = context.WithValue(ctx, ProjectIDKey, uuid.MustParse(claims.ProjectID))
			ctx = context.WithValue(ctx, PoolIDKey, uuid.MustParse(claims.PoolID))
			ctx = context.WithValue(ctx, ScopesKey, claims.Scopes)
		}

		return handler(ctx, req)
	}
}

func GetAuthbaseProjectID(ctx context.Context) (uuid.UUID, error) {
	pid, ok := ctx.Value(ProjectIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, status.Errorf(codes.InvalidArgument, "missing project id")
	}

	return pid, nil
}

func GetAuthbasePoolID(ctx context.Context) (uuid.UUID, error) {
	pid, ok := ctx.Value(PoolIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, status.Errorf(codes.InvalidArgument, "missing pool id")
	}

	return pid, nil
}

func GetAuthbaseAccountID(ctx context.Context) (uuid.UUID, error) {
	accountID, ok := ctx.Value(AccountIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.New("accountID not found in context")
	}

	return accountID, nil
}

func IsAuthbaseTokenMissing(ctx context.Context) bool {
	missing, ok := ctx.Value(TokenMissingKey).(bool)
	if !ok {
		return false
	}

	return missing
}

func GetAuthbaseProjectPermission(ctx context.Context) (v1.Permission, error) {
	permission, ok := ctx.Value(ProjectPermissionKey).(v1.Permission)
	if !ok {
		return v1.Permission_NONE, status.Errorf(codes.InvalidArgument, "missing project permission")
	}

	return permission, nil
}

func GetAuthbaseScopes(ctx context.Context) ([]string, error) {
	scopes, ok := ctx.Value(ScopesKey).([]string)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing scopes")
	}

	return scopes, nil
}
