package service

import (
	"context"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/model"
	"github.com/emrgen/authbase/pkg/permission"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/emrgen/authbase/x"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// NewPoolService creates a new account pool service.
func NewPoolService(store store.Provider, perm permission.MemberPermission) v1.PoolServiceServer {
	return &PoolService{
		store: store,
		perm:  perm,
	}
}

var _ v1.PoolServiceServer = new(PoolService)

// PoolService is the service for managing account pools.
type PoolService struct {
	store store.Provider
	perm  permission.MemberPermission
	v1.UnimplementedPoolServiceServer
}

// CreatePool creates a new pool for the given project. A pool is a group of accounts that can be managed together.
func (p *PoolService) CreatePool(ctx context.Context, request *v1.CreatePoolRequest) (*v1.CreatePoolResponse, error) {
	var err error
	as, err := store.GetProjectStore(ctx, p.store)
	if err != nil {
		return nil, err
	}
	accountID, err := x.GetAuthbaseAccountID(ctx)
	if err != nil {
		return nil, err
	}

	// check if the user has permission to create a pool for the project
	// TODO: move the permission checks to interceptor
	err = p.perm.CheckProjectPermission(ctx, accountID, permission.ProjectPermissionWrite)
	if err != nil {
		return nil, err
	}

	projectID := uuid.MustParse(request.GetProjectId())
	name := request.GetName()

	pool := &model.Pool{
		ID:        uuid.New().String(),
		Name:      name,
		ProjectID: projectID.String(),
	}

	// notice that the account is not part of the pool
	member := &model.PoolMember{
		PoolID:     pool.ID,
		AccountID:  accountID.String(),
		Permission: uint32(v1.Permission_OWNER),
	}

	err = as.Transaction(func(tx store.AuthBaseStore) error {
		err := tx.CreatePool(ctx, pool)
		if err != nil {
			return err
		}

		err = tx.AddPoolMember(ctx, member)
		if err != nil {
			return err
		}

		logrus.Infof("pool created: %v", request.Client)

		// NOTE: once the pool is created it cant be used to create new accounts
		// first the owner needs to create a client for the pool
		// if the request is to create a default client for the newly created pool
		if request.GetClient() {
			// we are saving the secret in the database,
			// so that the user can check it later as client config
			// TODO: move the secret to a secure vault, to avoid storing it in the database
			client := model.Client{
				ID:          uuid.New().String(),
				PoolID:      pool.ID,
				Name:        "default",
				CreatedByID: accountID.String(),
			}
			err = tx.CreateClient(ctx, &client)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.CreatePoolResponse{
		Pool: &v1.Pool{
			Id:   pool.ID,
			Name: pool.Name,
		},
	}, nil
}

// GetPool gets the pool by ID.
func (p *PoolService) GetPool(ctx context.Context, request *v1.GetPoolRequest) (*v1.GetPoolResponse, error) {
	poolID := uuid.MustParse(request.GetPoolId())
	as, err := store.GetProjectStore(ctx, p.store)
	if err != nil {
		return nil, err
	}

	pool, err := as.GetPoolByID(ctx, poolID)
	if err != nil {
		return nil, err
	}

	// check if the user has permission to get the pool
	projectID := uuid.MustParse(pool.ProjectID)
	err = p.perm.CheckProjectPermission(ctx, projectID, permission.ProjectPermissionRead)
	if err != nil {
		return nil, err
	}

	return &v1.GetPoolResponse{
		Pool: &v1.Pool{
			Id:        pool.ID,
			Name:      pool.Name,
			CreatedAt: timestamppb.New(pool.CreatedAt),
			UpdatedAt: timestamppb.New(pool.UpdatedAt),
		},
	}, nil
}

// ListPools lists all pools for the given project.
func (p *PoolService) ListPools(ctx context.Context, request *v1.ListPoolsRequest) (*v1.ListPoolsResponse, error) {
	projectID := uuid.MustParse(request.GetProjectId())
	as, err := store.GetProjectStore(ctx, p.store)
	if err != nil {
		return nil, err
	}

	// check if the user has permission to list pools for the project
	err = p.perm.CheckProjectPermission(ctx, projectID, permission.ProjectPermissionRead)
	if err != nil {
		return nil, err
	}

	page := x.GetPageFromRequest(request)

	pools, total, err := as.ListPools(ctx, projectID, int(page.Page), int(page.Size))
	if err != nil {
		return nil, err
	}

	var poolProtos []*v1.Pool
	for _, pool := range pools {
		poolProtos = append(poolProtos, &v1.Pool{
			Id:        pool.ID,
			Name:      pool.Name,
			ProjectId: projectID.String(),
			CreatedAt: timestamppb.New(pool.CreatedAt),
			UpdatedAt: timestamppb.New(pool.UpdatedAt),
		})
	}

	return &v1.ListPoolsResponse{
		Pools: poolProtos,
		Meta: &v1.Meta{
			Total: int32(total),
			Page:  page.Page,
			Size:  page.Size,
		},
	}, nil
}

// UpdatePool updates the pool with the given ID.
func (p *PoolService) UpdatePool(ctx context.Context, request *v1.UpdatePoolRequest) (*v1.UpdatePoolResponse, error) {
	poolID := uuid.MustParse(request.GetPoolId())
	as, err := store.GetProjectStore(ctx, p.store)
	if err != nil {
		return nil, err
	}

	// check if the user has permission to update the pool
	accountID, err := x.GetAuthbaseAccountID(ctx)
	if err != nil {
		return nil, err
	}

	err = p.perm.CheckProjectPermission(ctx, accountID, permission.ProjectPermissionWrite)
	if err != nil {
		return nil, err
	}

	err = as.Transaction(func(tx store.AuthBaseStore) error {
		pool, err := tx.GetPoolByID(ctx, poolID)
		if err != nil {
			return err
		}

		// check if the user has permission to update the pool
		projectID := uuid.MustParse(pool.ProjectID)
		err = p.perm.CheckProjectPermission(ctx, projectID, permission.ProjectPermissionWrite)
		if err != nil {
			return err
		}

		pool.Name = request.GetName()
		err = tx.UpdatePool(ctx, pool)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.UpdatePoolResponse{
		Pool: &v1.Pool{
			Id:   poolID.String(),
			Name: request.GetName(),
		},
	}, nil
}

// DeletePool deletes the pool with the given ID.
func (p *PoolService) DeletePool(ctx context.Context, request *v1.DeletePoolRequest) (*v1.DeletePoolResponse, error) {
	as, err := store.GetProjectStore(ctx, p.store)
	if err != nil {
		return nil, err
	}

	poolID := uuid.MustParse(request.GetPoolId())
	err = as.Transaction(func(tx store.AuthBaseStore) error {
		pool, err := tx.GetPoolByID(ctx, poolID)
		if err != nil {
			return err
		}

		// Check if the pool is default
		if pool.Default {
			return status.Error(codes.FailedPrecondition, "cannot delete default pool")
		}

		err = tx.DeletePool(ctx, poolID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.DeletePoolResponse{}, nil
}
