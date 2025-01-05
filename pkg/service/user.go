package service

import (
	"context"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/cache"
	"github.com/emrgen/authbase/pkg/model"
	"github.com/emrgen/authbase/pkg/permission"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/emrgen/authbase/x"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ v1.UserServiceServer = new(UserService)

type UserService struct {
	perm  permission.AuthBasePermission
	store store.Provider
	cache *cache.Redis
	v1.UnimplementedUserServiceServer
}

// NewUserService creates a new user service.
func NewUserService(perm permission.AuthBasePermission, store store.Provider, cache *cache.Redis) v1.UserServiceServer {
	return &UserService{perm: perm, store: store, cache: cache}
}

// CreateUser creates a new user.
func (u *UserService) CreateUser(ctx context.Context, request *v1.CreateUserRequest) (*v1.CreateUserResponse, error) {
	var err error
	orgID, err := uuid.Parse(request.GetProjectId())
	if err != nil {
		return nil, err
	}

	err = u.perm.CheckProjectPermission(ctx, orgID, "write")
	if err != nil {
		return nil, err
	}

	// create a new user
	as, err := store.GetProjectStore(ctx, u.store)
	if err != nil {
		return nil, err
	}

	user := model.User{
		ID:        uuid.New().String(),
		Username:  request.GetUsername(),
		Email:     request.GetEmail(),
		ProjectID: orgID.String(),
	}

	password := request.GetPassword()
	if password != "" {
		salt := x.Keygen()
		hashedPassword, err := x.HashPassword(password, salt)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
		user.Salt = salt
	}

	if err := as.CreateUser(ctx, &user); err != nil {
		return nil, err
	}

	return &v1.CreateUserResponse{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}, nil
}

// GetUser gets a user by ID.
func (u *UserService) GetUser(ctx context.Context, request *v1.GetUserRequest) (*v1.GetUserResponse, error) {
	// get the user
	as, err := store.GetProjectStore(ctx, u.store)
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, err
	}

	user, err := as.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = u.perm.CheckProjectPermission(ctx, uuid.MustParse(user.ProjectID), "read")
	if err != nil {
		return nil, err
	}

	return &v1.GetUserResponse{
		User: &v1.User{
			Id:       user.ID,
			Username: user.Username,
		},
	}, nil
}

// ListUsers lists users.
func (u *UserService) ListUsers(ctx context.Context, request *v1.ListUsersRequest) (*v1.ListUsersResponse, error) {
	as, err := store.GetProjectStore(ctx, u.store)
	if err != nil {
		return nil, err
	}

	orgID, err := uuid.Parse(request.GetProjectId())
	if err != nil {
		return nil, err
	}

	err = u.perm.CheckProjectPermission(ctx, orgID, "read")
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	page := x.GetPageFromRequest(request)
	users, total, err := as.ListUsersByOrg(ctx, false, orgID, int(page.Page), int(page.Size))

	var userProtoList []*v1.User
	for _, user := range users {
		userProtoList = append(userProtoList, &v1.User{
			Id:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
			Disabled:  user.Disabled,
			Member:    user.Member,
		})
	}

	return &v1.ListUsersResponse{Users: userProtoList, Meta: &v1.Meta{Total: int32(total)}}, nil
}

// UpdateUser updates a user.
func (u *UserService) UpdateUser(ctx context.Context, request *v1.UpdateUserRequest) (*v1.UpdateUserResponse, error) {
	as, err := store.GetProjectStore(ctx, u.store)
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, err
	}

	err = as.Transaction(func(tx store.AuthBaseStore) error {
		user, err := tx.GetUserByID(ctx, id)
		if err != nil {
			return err
		}

		err = u.perm.CheckProjectPermission(ctx, uuid.MustParse(user.ProjectID), "write")
		if err != nil {
			return err
		}

		if request.GetUsername() != "" {
			user.Username = request.GetUsername()
		}

		if request.GetEmail() != "" {
			user.Email = request.GetEmail()
		}

		err = tx.UpdateUser(ctx, user)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.UpdateUserResponse{
		User: &v1.User{
			Id: id.String(),
		},
	}, nil
}

// DeleteUser deletes a user.
func (u *UserService) DeleteUser(ctx context.Context, request *v1.DeleteUserRequest) (*v1.DeleteUserResponse, error) {
	as, err := store.GetProjectStore(ctx, u.store)
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, err
	}

	user, err := as.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = u.perm.CheckProjectPermission(ctx, uuid.MustParse(user.ProjectID), "write")
	if err != nil {
		return nil, err
	}

	err = as.DeleteUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return &v1.DeleteUserResponse{
		Message: "User deleted successfully.",
	}, nil
}

// ListActiveUsers lists active users.
func (u *UserService) ListActiveUsers(ctx context.Context, request *v1.ListActiveUsersRequest) (*v1.ListActiveUsersResponse, error) {
	as, err := store.GetProjectStore(ctx, u.store)
	if err != nil {
		return nil, err
	}

	orgID, err := uuid.Parse(request.GetProjectId())
	if err != nil {
		return nil, err
	}

	err = u.perm.CheckProjectPermission(ctx, orgID, "read")
	if err != nil {
		return nil, err
	}

	page := x.GetPageFromRequest(request)

	sessions, err := as.ListSessions(ctx, orgID, int(page.Page), int(page.Size))
	if err != nil {
		return nil, err
	}

	var userProtos []*v1.User
	for _, session := range sessions {
		user := session.User
		if user == nil {
			continue
		}
		userProtos = append(userProtos, &v1.User{
			Id:        session.UserID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		})
	}

	return &v1.ListActiveUsersResponse{Users: userProtos}, nil
}

// DisableUser activates a user.
func (u *UserService) DisableUser(ctx context.Context, request *v1.DisableUserRequest) (*v1.DisableUserResponse, error) {
	as, err := store.GetProjectStore(ctx, u.store)
	if err != nil {
		return nil, err
	}
	userID, err := uuid.Parse(request.GetUserId())
	if err != nil {
		return nil, err
	}

	user, err := as.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	err = u.perm.CheckProjectPermission(ctx, uuid.MustParse(user.ProjectID), "write")
	if err != nil {
		return nil, err
	}

	user.Disabled = true

	err = as.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &v1.DisableUserResponse{
		Message: "User disabled successfully.",
	}, nil
}

// EnableUser activates a user.
func (u *UserService) EnableUser(ctx context.Context, request *v1.EnableUserRequest) (*v1.EnableUserResponse, error) {
	as, err := store.GetProjectStore(ctx, u.store)
	if err != nil {
		return nil, err
	}
	userID, err := uuid.Parse(request.GetUserId())
	if err != nil {
		return nil, err
	}

	user, err := as.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	err = u.perm.CheckProjectPermission(ctx, uuid.MustParse(user.ProjectID), "write")
	if err != nil {
		return nil, err
	}

	user.Disabled = false

	err = as.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &v1.EnableUserResponse{
		Message: "User enabled successfully.",
	}, nil
}
