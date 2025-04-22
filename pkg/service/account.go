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
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
)

var _ v1.AccountServiceServer = new(AccountService)

type AccountService struct {
	perm  permission.AuthBasePermission
	store store.Provider
	cache *cache.Redis
	v1.UnimplementedAccountServiceServer
}

// NewAccountService creates a new user service.
func NewAccountService(perm permission.AuthBasePermission, store store.Provider, cache *cache.Redis) v1.AccountServiceServer {
	return &AccountService{perm: perm, store: store, cache: cache}
}

// CreateAccount creates a new user.
func (u *AccountService) CreateAccount(ctx context.Context, request *v1.CreateAccountRequest) (*v1.CreateAccountResponse, error) {
	var err error
	if request.GetPoolId() == "" {
		return nil, status.Error(codes.InvalidArgument, "client_id is required")
	}

	// create a new user
	as, err := store.GetProjectStore(ctx, u.store)
	if err != nil {
		return nil, err
	}

	poolID, err := uuid.Parse(request.GetPoolId())
	if err != nil {
		return nil, err
	}

	pool, err := as.GetPoolByID(ctx, poolID)
	if err != nil {
		return nil, err
	}

	projectID, err := uuid.Parse(pool.ProjectID)
	if err != nil {
		return nil, err
	}

	visibleName := request.GetVisibleName()
	userName := request.GetUsername()

	// if visibleName is empty, use the username or the email prefix as the visible name
	if visibleName == "" && userName != "" {
		visibleName = userName
	}
	if visibleName == "" {
		visibleName = strings.Split(request.GetEmail(), "@")[0]
	}

	err = u.perm.CheckProjectPermission(ctx, projectID, "write")
	if err != nil {
		return nil, err
	}

	user := model.Account{
		ID:          uuid.New().String(),
		Username:    request.GetUsername(),
		Email:       request.GetEmail(),
		VisibleName: request.GetVisibleName(),
		ProjectID:   projectID.String(),
		PoolID:      poolID.String(),
	}

	password := request.GetPassword()
	if password != "" {
		salt := x.GenerateSalt()
		hashedPassword := x.HashPassword(password, salt)
		user.PasswordHash = string(hashedPassword)
		user.Salt = salt
	}

	if err := as.CreateAccount(ctx, &user); err != nil {
		return nil, err
	}

	return &v1.CreateAccountResponse{
		Account: &v1.Account{
			Id:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			PoolId:    poolID.String(),
			CreatedAt: timestamppb.New(user.CreatedAt),
		},
	}, nil
}

// GetCurrentAccount gets the current user.
func (u *AccountService) GetCurrentAccount(ctx context.Context, request *v1.GetCurrentAccountRequest) (*v1.GetCurrentAccountResponse, error) {
	as, err := store.GetProjectStore(ctx, u.store)
	if err != nil {
		return nil, err
	}

	accountID, err := x.GetAuthbaseAccountID(ctx)
	if err != nil {
		return nil, err
	}

	user, err := as.GetAccountByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	err = u.perm.CheckProjectPermission(ctx, uuid.MustParse(user.ProjectID), "read")
	if err != nil {
		return nil, err
	}

	return &v1.GetCurrentAccountResponse{
		Account: &v1.Account{
			Id:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			Disabled:  user.Disabled,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
			Member:    user.ProjectMember,
		},
	}, nil
}

// GetAccount gets a user by ID.
func (u *AccountService) GetAccount(ctx context.Context, request *v1.GetAccountRequest) (*v1.GetAccountResponse, error) {
	// get the user
	as, err := store.GetProjectStore(ctx, u.store)
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, err
	}

	user, err := as.GetAccountByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = u.perm.CheckProjectPermission(ctx, uuid.MustParse(user.ProjectID), "read")
	if err != nil {
		return nil, err
	}

	return &v1.GetAccountResponse{
		Account: &v1.Account{
			Id:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			Disabled:  user.Disabled,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
			Member:    user.ProjectMember,
		},
	}, nil
}

// ListAccounts lists users.
func (u *AccountService) ListAccounts(ctx context.Context, request *v1.ListAccountsRequest) (*v1.ListAccountsResponse, error) {
	var err error

	if request.GetProjectId() == "" && request.GetPoolId() == "" {
		return nil, status.Error(codes.InvalidArgument, "project_id or pool_id is required")
	}
	as, err := store.GetProjectStore(ctx, u.store)
	if err != nil {
		return nil, err
	}

	var projectID uuid.UUID
	var poolID uuid.UUID
	if request.GetPoolId() != "" {
		poolID, err = uuid.Parse(request.GetPoolId())
		if err != nil {
			return nil, err
		}
		pool, err := as.GetPoolByID(ctx, poolID)
		if err != nil {
			return nil, err
		}
		projectID = uuid.MustParse(pool.ProjectID)
	} else {
		projectID, err = uuid.Parse(request.GetProjectId())
		if err != nil {
			return nil, err
		}
	}

	err = u.perm.CheckProjectPermission(ctx, projectID, "read")
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	page := x.GetPageFromRequest(request)
	var users []*model.Account
	var total int

	if poolID != uuid.Nil {
		users, total, err = as.ListPoolAccounts(ctx, false, poolID, int(page.Page), int(page.Size))
	} else {
		users, total, err = as.ListProjectAccounts(ctx, false, projectID, int(page.Page), int(page.Size))
	}

	var userProtoList []*v1.Account
	for _, user := range users {
		userProtoList = append(userProtoList, &v1.Account{
			Id:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			VisibleName: user.VisibleName,
			CreatedAt:   timestamppb.New(user.CreatedAt),
			UpdatedAt:   timestamppb.New(user.UpdatedAt),
			Disabled:    user.Disabled,
			Member:      user.ProjectMember,
		})
	}

	return &v1.ListAccountsResponse{Accounts: userProtoList, Meta: &v1.Meta{Total: int32(total)}}, nil
}

// UpdateAccount updates a user.
func (u *AccountService) UpdateAccount(ctx context.Context, request *v1.UpdateAccountRequest) (*v1.UpdateAccountResponse, error) {
	as, err := store.GetProjectStore(ctx, u.store)
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, err
	}

	err = as.Transaction(func(tx store.AuthBaseStore) error {
		user, err := tx.GetAccountByID(ctx, id)
		if err != nil {
			return err
		}

		err = u.perm.CheckProjectPermission(ctx, uuid.MustParse(user.ProjectID), "write")
		if err != nil {
			return err
		}

		user.VisibleName = request.GetVisibleName()

		err = tx.UpdateAccount(ctx, user)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.UpdateAccountResponse{
		Account: &v1.Account{
			Id: id.String(),
		},
	}, nil
}

// DeleteAccount deletes a user.
func (u *AccountService) DeleteAccount(ctx context.Context, request *v1.DeleteAccountRequest) (*v1.DeleteAccountResponse, error) {
	as, err := store.GetProjectStore(ctx, u.store)
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, err
	}

	user, err := as.GetAccountByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = u.perm.CheckProjectPermission(ctx, uuid.MustParse(user.ProjectID), "write")
	if err != nil {
		return nil, err
	}

	err = as.DeleteAccount(ctx, id)
	if err != nil {
		return nil, err
	}

	return &v1.DeleteAccountResponse{
		Message: "Account deleted successfully.",
	}, nil
}

// ListActiveAccounts lists active users. Users that have logged in recently and the session is still valid.
func (u *AccountService) ListActiveAccounts(ctx context.Context, request *v1.ListActiveAccountsRequest) (*v1.ListActiveAccountsResponse, error) {
	as, err := store.GetProjectStore(ctx, u.store)
	if err != nil {
		return nil, err
	}

	poolID, err := uuid.Parse(request.GetPoolId())
	if err != nil {
		return nil, err
	}
	pool, err := as.GetPoolByID(ctx, poolID)
	if err != nil {
		return nil, err
	}
	projectID := uuid.MustParse(pool.ProjectID)

	err = u.perm.CheckProjectPermission(ctx, projectID, "read")
	if err != nil {
		return nil, err
	}

	page := x.GetPageFromRequest(request)

	logrus.Infof("pool id: %s", poolID.String())

	sessions, err := as.ListActiveAccounts(ctx, poolID, int(page.Page), int(page.Size))
	if err != nil {
		return nil, err
	}

	var userProtos []*v1.Account
	for _, session := range sessions {
		user := session.Account
		if user == nil {
			continue
		}
		userProtos = append(userProtos, &v1.Account{
			Id:          session.AccountID,
			PoolId:      user.PoolID,
			ProjectId:   user.ProjectID,
			Username:    user.Username,
			VisibleName: user.VisibleName,
			Email:       user.Email,
			CreatedAt:   timestamppb.New(user.CreatedAt),
			UpdatedAt:   timestamppb.New(user.UpdatedAt),
		})
	}

	return &v1.ListActiveAccountsResponse{Accounts: userProtos, Meta: &v1.Meta{Total: 0, Page: 0, Size: 0}}, nil
}

// ListInactiveAccounts lists inactive users. Users that have not logged in for a while and the session is expired.
func (u *AccountService) ListInactiveAccounts(ctx context.Context, request *v1.ListInactiveAccountsRequest) (*v1.ListInactiveAccountsResponse, error) {
	//as, err := store.GetProjectStore(ctx, u.store)
	//if err != nil {
	//	return nil, err
	//}
	//
	//orgID, err := uuid.Parse(request.GetProjectId())
	//if err != nil {
	//	return nil, err
	//}
	//
	//err = u.perm.CheckProjectPermission(ctx, orgID, "read")
	//if err != nil {
	//	return nil, err
	//}
	//
	//page := x.GetPageFromRequest(request)
	//
	//sessions, err := as.ListInactiveSessions(ctx, orgID, int(page.Page), int(page.Size))
	//if err != nil {
	//	return nil, err
	//}
	//
	//var userProtos []*v1.Account
	//for _, session := range sessions {
	//	user := session.Account
	//	if user == nil {
	//		continue
	//	}
	//	userProtos = append(userProtos, &v1.Account{
	//		Id:          session.AccountID,
	//		Accountname: user.Accountname,
	//		Email:       user.Email,
	//		CreatedAt:   timestamppb.New(user.CreatedAt),
	//		UpdatedAt:   timestamppb.New(user.UpdatedAt),
	//	})
	//}

	return &v1.ListInactiveAccountsResponse{}, nil
}

// DisableAccount activates a user. Disabled users cannot log in.
func (u *AccountService) DisableAccount(ctx context.Context, request *v1.DisableAccountRequest) (*v1.DisableAccountResponse, error) {
	as, err := store.GetProjectStore(ctx, u.store)
	if err != nil {
		return nil, err
	}
	userID, err := uuid.Parse(request.GetAccountId())
	if err != nil {
		return nil, err
	}

	user, err := as.GetAccountByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	err = u.perm.CheckProjectPermission(ctx, uuid.MustParse(user.ProjectID), "write")
	if err != nil {
		return nil, err
	}

	user.Disabled = true

	err = as.UpdateAccount(ctx, user)
	if err != nil {
		return nil, err
	}

	return &v1.DisableAccountResponse{
		Message: "Account disabled successfully.",
	}, nil
}

// EnableAccount activates a user. Enabled users can log in.
func (u *AccountService) EnableAccount(ctx context.Context, request *v1.EnableAccountRequest) (*v1.EnableAccountResponse, error) {
	as, err := store.GetProjectStore(ctx, u.store)
	if err != nil {
		return nil, err
	}
	userID, err := uuid.Parse(request.GetAccountId())
	if err != nil {
		return nil, err
	}

	user, err := as.GetAccountByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	err = u.perm.CheckProjectPermission(ctx, uuid.MustParse(user.ProjectID), "write")
	if err != nil {
		return nil, err
	}

	user.Disabled = false

	err = as.UpdateAccount(ctx, user)
	if err != nil {
		return nil, err
	}

	return &v1.EnableAccountResponse{
		Message: "Account enabled successfully.",
	}, nil
}
