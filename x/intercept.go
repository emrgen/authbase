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
	// UserIDKey is the key to store the user id in the context
	UserIDKey = "authbase_user_id"
)

type ProjectID interface {
	GetProjectId() string
}

func InjectPermissionInterceptor(member v1.MemberServiceClient) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		// check if the projectPermission are in the context
		_, err = GetProjectPermission(ctx)
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

		userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
		if !ok {
			return nil, status.Errorf(codes.InvalidArgument, "missing user id")
		}

		md, ok := metadata.FromIncomingContext(ctx)
		token := md.Get("authorization")
		md.Set("Authorization", token[0])
		ctx = metadata.NewOutgoingContext(ctx, md)

		// GetProjectMembership is a gRPC call to the project membership service
		res, err := member.GetMember(ctx, &v1.GetMemberRequest{
			ProjectId: projectUUID.String(),
			MemberId:  userID.String(),
		})
		if err != nil {
			return nil, err
		}

		ctx = context.WithValue(ctx, ProjectPermissionKey, res.Member.Permission)
		ctx = context.WithValue(ctx, ProjectIDKey, projectUUID)

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

func GetAuthbaseUserID(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.New("userID not found in context")
	}

	return userID, nil
}

func GetProjectPermission(ctx context.Context) (v1.Permission, error) {
	permission, ok := ctx.Value(ProjectPermissionKey).(v1.Permission)
	if !ok {
		return v1.Permission_NONE, status.Errorf(codes.InvalidArgument, "missing project permission")
	}

	return permission, nil
}
