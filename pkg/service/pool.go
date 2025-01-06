package service

import (
	"context"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/model"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/emrgen/authbase/x"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// NewPoolService creates a new pool service.
func NewPoolService(store store.Provider) v1.PoolServiceServer {
	return &PoolService{
		store: store,
	}
}

var _ v1.PoolServiceServer = new(PoolService)

type PoolService struct {
	store store.Provider
	v1.UnimplementedPoolServiceServer
}

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

	projectID := uuid.MustParse(request.GetProjectId())
	name := request.GetName()

	pool := &model.Pool{
		ID:        uuid.New().String(),
		Name:      name,
		ProjectID: projectID.String(),
	}
 
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

func (p *PoolService) GetPool(ctx context.Context, request *v1.GetPoolRequest) (*v1.GetPoolResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PoolService) ListPools(ctx context.Context, request *v1.ListPoolsRequest) (*v1.ListPoolsResponse, error) {
	projectID := uuid.MustParse(request.GetProjectId())
	as, err := store.GetProjectStore(ctx, p.store)
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

func (p *PoolService) UpdatePool(ctx context.Context, request *v1.UpdatePoolRequest) (*v1.UpdatePoolResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PoolService) DeletePool(ctx context.Context, request *v1.DeletePoolRequest) (*v1.DeletePoolResponse, error) {
	//TODO implement me
	panic("implement me")
}
