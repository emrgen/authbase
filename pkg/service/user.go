package service

import (
	"context"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/cache"
	"github.com/emrgen/authbase/pkg/model"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/emrgen/authbase/x"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ v1.UserServiceServer = new(UserService)

type UserService struct {
	store store.AuthBaseStore
	cache *cache.Redis
	v1.UnimplementedUserServiceServer
}

// NewUserService creates a new user service.
func NewUserService(store store.AuthBaseStore, cache *cache.Redis) *UserService {
	return &UserService{store: store, cache: cache}
}

// CreateUser creates a new user.
func (u *UserService) CreateUser(ctx context.Context, request *v1.CreateUserRequest) (*v1.CreateUserResponse, error) {
	user := model.User{
		ID:             uuid.New().String(),
		Username:       request.GetUsername(),
		Email:          request.GetEmail(),
		OrganizationID: request.GetOrganizationId(),
	}

	if err := u.store.CreateUser(ctx, &user); err != nil {
		return nil, err
	}

	return &v1.CreateUserResponse{
		Id: user.ID,
	}, nil
}

// GetUser gets a user by ID.
func (u *UserService) GetUser(ctx context.Context, request *v1.GetUserRequest) (*v1.GetUserResponse, error) {
	id, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, err
	}

	user, err := u.store.GetUserByID(ctx, id)
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
	orgID, err := uuid.Parse(request.GetOrganizationId())
	if err != nil {
		return nil, err
	}

	page := x.GetPageFromRequest(request)
	users, total, err := u.store.ListUsersByOrg(ctx, false, orgID, int(page.Page), int(page.Size))

	var userProtos []*v1.User
	for _, user := range users {
		userProtos = append(userProtos, &v1.User{
			Id:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		})
	}

	return &v1.ListUsersResponse{Users: userProtos, Meta: &v1.Meta{Total: int32(total)}}, nil
}

// UpdateUser updates a user.
func (u *UserService) UpdateUser(ctx context.Context, request *v1.UpdateUserRequest) (*v1.UpdateUserResponse, error) {
	id, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, err
	}

	err = u.store.Transaction(func(tx store.AuthBaseStore) error {
		user, err := tx.GetUserByID(ctx, id)
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
	id, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, err
	}

	err = u.store.DeleteUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return &v1.DeleteUserResponse{
		Message: "User deleted successfully.",
	}, nil
}
