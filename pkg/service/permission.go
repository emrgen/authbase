package service

import (
	"context"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/cache"
	"github.com/emrgen/authbase/pkg/store"
)

var _ v1.PermissionServiceServer = new(PermissionService)

type PermissionService struct {
	store store.AuthBaseStore
	cache *cache.Redis
	v1.UnimplementedPermissionServiceServer
}

func NewPermissionService(store store.AuthBaseStore, cache *cache.Redis) *PermissionService {
	return &PermissionService{store: store, cache: cache}
}

func (p *PermissionService) CreatePermission(ctx context.Context, request *v1.CreatePermissionRequest) (*v1.CreatePermissionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PermissionService) GetPermission(ctx context.Context, request *v1.GetPermissionRequest) (*v1.GetPermissionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PermissionService) UpdatePermission(ctx context.Context, request *v1.UpdatePermissionRequest) (*v1.UpdatePermissionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PermissionService) DeletePermission(ctx context.Context, request *v1.DeletePermissionRequest) (*v1.DeletePermissionResponse, error) {
	//TODO implement me
	panic("implement me")
}
