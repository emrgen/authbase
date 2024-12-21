package permission

import (
	"context"
	v1 "github.com/emrgen/authbase/apis/v1"
)

var _ v1.UserServiceServer = new(CheckedUserService)

type CheckedUserService struct {
	perm AuthBasePermission
	base v1.UserServiceServer
	v1.UnimplementedUserServiceServer
}

func NewCheckedUserService(base v1.UserServiceServer, perm AuthBasePermission) v1.UserServiceServer {
	return &CheckedUserService{perm: perm, base: base}
}

func (u *CheckedUserService) CreateUser(ctx context.Context, request *v1.CreateUserRequest) (*v1.CreateUserResponse, error) {
	// permission check
	//userID, err := x.GetUserID(ctx)
	//if err != nil {
	//	return nil, err
	//}
	//
	//orgID, err := uuid.Parse(request.GetOrganizationId())
	//if err != nil {
	//	return nil, err
	//}
	//
	//ok, err := u.perm.CheckUserPermission(ctx, userID, orgID, "organization", "create_user")
	//if err != nil {
	//	return nil, err
	//}
	//if !ok {
	//	return nil, x.ErrForbidden
	//}

	return u.base.CreateUser(ctx, request)
}

func (u *CheckedUserService) GetUser(ctx context.Context, request *v1.GetUserRequest) (*v1.GetUserResponse, error) {
	//userID, err := x.GetUserID(ctx)
	//if err != nil {
	//	return nil, err
	//}
	//
	//orgID, err := uuid.Parse(request.GetOrganizationId())
	//if err != nil {
	//	return nil, err
	//}
	//
	//ok, err := u.perm.CheckUserPermission(ctx, userID, orgID, "organization", "viewer")
	//if err != nil {
	//	return nil, err
	//}
	//if !ok {
	//	return nil, x.ErrForbidden
	//}

	return u.base.GetUser(ctx, request)
}

func (u *CheckedUserService) ListUsers(ctx context.Context, request *v1.ListUsersRequest) (*v1.ListUsersResponse, error) {
	return u.base.ListUsers(ctx, request)
}

func (u *CheckedUserService) UpdateUser(ctx context.Context, request *v1.UpdateUserRequest) (*v1.UpdateUserResponse, error) {
	return u.base.UpdateUser(ctx, request)
}

func (u *CheckedUserService) DeleteUser(ctx context.Context, request *v1.DeleteUserRequest) (*v1.DeleteUserResponse, error) {
	return u.base.DeleteUser(ctx, request)
}

func (u *CheckedUserService) ActiveUsers(ctx context.Context, request *v1.ActiveUsersRequest) (*v1.ActiveUsersResponse, error) {
	return u.base.ActiveUsers(ctx, request)
}

func (u *CheckedUserService) DeactivateUser(ctx context.Context, request *v1.DeactivateUserRequest) (*v1.DeactivateUserResponse, error) {
	return u.base.DeactivateUser(ctx, request)
}
