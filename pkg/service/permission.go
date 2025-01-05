package service

import (
	"context"

	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/cache"
	"github.com/emrgen/authbase/pkg/model"
	"github.com/emrgen/authbase/pkg/permission"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/google/uuid"
)

var _ v1.PermissionServiceServer = new(PermissionService)

type PermissionService struct {
	perm  permission.AuthBasePermission
	store store.Provider
	cache *cache.Redis
	v1.UnimplementedPermissionServiceServer
}

func NewPermissionService(perm permission.AuthBasePermission, store store.Provider, cache *cache.Redis) *PermissionService {
	return &PermissionService{perm: perm, store: store, cache: cache}
}

// CreatePermission creates a new permission for a member in an project
func (p *PermissionService) CreatePermission(ctx context.Context, request *v1.CreatePermissionRequest) (*v1.CreatePermissionResponse, error) {
	as, err := store.GetProjectStore(ctx, p.store)
	if err != nil {
		return nil, err
	}

	orgID, err := uuid.Parse(request.GetProjectId())
	if err != nil {
		return nil, err
	}

	err = p.perm.CheckProjectPermission(ctx, orgID, "write")
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(request.GetMemberId())
	if err != nil {
		return nil, err
	}

	user, err := as.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	err = p.perm.CheckProjectPermission(ctx, uuid.MustParse(user.ProjectID), "write")
	if err != nil {
		return nil, err
	}

	userPermission := uint32(0)
	for _, perm := range request.GetPermissions() {
		userPermission |= uint32(perm)
	}

	permissionModel := model.ProjectMember{
		ProjectID:  orgID.String(),
		UserID:     userID.String(),
		Permission: userPermission,
	}

	err = as.CreateProjectMember(ctx, &permissionModel)
	if err != nil {
		return nil, err
	}

	return &v1.CreatePermissionResponse{Message: "ProjectMember created successfully"}, nil
}

// GetPermission gets the permission of a member in an project
func (p *PermissionService) GetPermission(ctx context.Context, request *v1.GetPermissionRequest) (*v1.GetPermissionResponse, error) {
	as, err := store.GetProjectStore(ctx, p.store)
	if err != nil {
		return nil, err
	}

	orgID, err := uuid.Parse(request.GetProjectId())
	if err != nil {
		return nil, err
	}

	err = p.perm.CheckProjectPermission(ctx, orgID, "read")
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(request.GetMemberId())
	if err != nil {
		return nil, err
	}

	perm, err := as.GetProjectMemberByID(ctx, orgID, userID)
	if err != nil {
		return nil, err
	}

	permissions := make([]v1.Permission, 0)
	for value := range v1.Permission_name {
		if perm.Permission&uint32(value) > 0 {
			permissions = append(permissions, v1.Permission(value))
		}
	}

	return &v1.GetPermissionResponse{
		Permissions: permissions,
	}, nil
}

// UpdatePermission updates the permission of a member in an project
// 1. check if the caller has project write permission
// 2.
func (p *PermissionService) UpdatePermission(ctx context.Context, request *v1.UpdatePermissionRequest) (*v1.UpdatePermissionResponse, error) {
	as, err := store.GetProjectStore(ctx, p.store)
	if err != nil {
		return nil, err
	}

	orgID, err := uuid.Parse(request.GetProjectId())
	if err != nil {
		return nil, err
	}

	err = p.perm.CheckProjectPermission(ctx, orgID, "write")
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(request.GetMemberId())
	if err != nil {
		return nil, err
	}

	// new permissions overwrite all old permission
	// example:
	// before: permission: Read | Write
	// update request: permission: Read
	// after: permission: Read
	userPermission := uint32(0)
	for _, perm := range request.GetPermissions() {
		userPermission |= uint32(perm)
	}

	err = as.Transaction(func(tx store.AuthBaseStore) error {
		perm, err := tx.GetProjectMemberByID(ctx, orgID, userID)
		if err != nil {
			return err
		}

		// Update the userPermission
		perm.Permission = userPermission
		err = tx.UpdateProjectMember(ctx, perm)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.UpdatePermissionResponse{Message: "ProjectMember updated successfully"}, nil
}

// DeletePermission deletes all permission of a member in an project
func (p *PermissionService) DeletePermission(ctx context.Context, request *v1.DeletePermissionRequest) (*v1.DeletePermissionResponse, error) {
	as, err := store.GetProjectStore(ctx, p.store)
	if err != nil {
		return nil, err
	}

	orgID, err := uuid.Parse(request.GetProjectId())
	if err != nil {
		return nil, err
	}

	err = p.perm.CheckProjectPermission(ctx, orgID, "write")
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(request.GetMemberId())
	if err != nil {
		return nil, err
	}

	err = as.DeleteProjectMember(ctx, orgID, userID)
	if err != nil {
		return nil, err
	}

	return &v1.DeletePermissionResponse{Message: "ProjectMember deleted successfully"}, nil
}
