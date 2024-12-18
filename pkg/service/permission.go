package service

import (
	"context"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/cache"
	"github.com/emrgen/authbase/pkg/model"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/google/uuid"
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
	orgID, err := uuid.Parse(request.GetOrganizationId())
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(request.GetMemberId())
	if err != nil {
		return nil, err
	}

	permission := uint32(0)
	for _, perm := range request.GetPermissions() {
		permission |= uint32(perm)
	}

	permissionModel := model.Permission{
		OrganizationID: orgID.String(),
		UserID:         userID.String(),
		Permission:     permission,
	}

	err = p.store.CreatePermission(ctx, &permissionModel)
	if err != nil {
		return nil, err
	}

	return &v1.CreatePermissionResponse{Message: "Permission created successfully"}, nil
}

func (p *PermissionService) GetPermission(ctx context.Context, request *v1.GetPermissionRequest) (*v1.GetPermissionResponse, error) {
	orgID, err := uuid.Parse(request.GetOrganizationId())
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(request.GetMemberId())
	if err != nil {
		return nil, err
	}

	permission, err := p.store.GetPermissionByID(ctx, orgID, userID)
	if err != nil {
		return nil, err
	}

	permissions := make([]v1.Permission, 0)
	for value, _ := range v1.Permission_name {
		if permission.Permission&uint32(value) > 0 {
			permissions = append(permissions, v1.Permission(value))
		}
	}

	return &v1.GetPermissionResponse{
		Permissions: permissions,
	}, nil
}

func (p *PermissionService) UpdatePermission(ctx context.Context, request *v1.UpdatePermissionRequest) (*v1.UpdatePermissionResponse, error) {
	orgID, err := uuid.Parse(request.GetOrganizationId())
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(request.GetMemberId())
	if err != nil {
		return nil, err
	}

	permission := uint32(0)
	for _, perm := range request.GetPermissions() {
		permission |= uint32(perm)
	}

	err = p.store.Transaction(func(tx store.AuthBaseStore) error {
		perm, err := tx.GetPermissionByID(ctx, orgID, userID)
		if err != nil {
			return err
		}

		// Update the permission
		perm.Permission = permission
		err = tx.UpdatePermission(ctx, perm)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.UpdatePermissionResponse{Message: "Permission updated successfully"}, nil
}

func (p *PermissionService) DeletePermission(ctx context.Context, request *v1.DeletePermissionRequest) (*v1.DeletePermissionResponse, error) {
	orgID, err := uuid.Parse(request.GetOrganizationId())
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(request.GetMemberId())
	if err != nil {
		return nil, err
	}

	err = p.store.DeletePermission(ctx, orgID, userID)
	if err != nil {
		return nil, err
	}

	return &v1.DeletePermissionResponse{Message: "Permission deleted successfully"}, nil
}
