package service

import (
	"context"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/store"
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
	//TODO implement me
	panic("implement me")
}

func (p *PoolService) GetPool(ctx context.Context, request *v1.GetPoolRequest) (*v1.GetPoolResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PoolService) ListPools(ctx context.Context, request *v1.ListPoolsRequest) (*v1.ListPoolsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PoolService) UpdatePool(ctx context.Context, request *v1.UpdatePoolRequest) (*v1.UpdatePoolResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PoolService) DeletePool(ctx context.Context, request *v1.DeletePoolRequest) (*v1.DeletePoolResponse, error) {
	//TODO implement me
	panic("implement me")
}
