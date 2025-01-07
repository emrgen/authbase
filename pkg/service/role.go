package service

import (
	"context"
	"encoding/json"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/model"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/emrgen/authbase/x"
	"github.com/google/uuid"
)

// NewRoleService creates a new role service.
func NewRoleService(store store.Provider) v1.RoleServiceServer {
	return &RoleService{store: store}
}

var _ v1.RoleServiceServer = (*RoleService)(nil)

// RoleService is the service for managing roles.
type RoleService struct {
	store store.Provider
	v1.UnimplementedRoleServiceServer
}

func (r *RoleService) CreateRole(ctx context.Context, request *v1.CreateRoleRequest) (*v1.CreateRoleResponse, error) {
	as, err := store.GetProjectStore(ctx, r.store)
	if err != nil {
		return nil, err
	}
	attributes := request.GetAttributes()
	attrJSON, err := json.Marshal(attributes)
	if err != nil {
		return nil, err
	}

	role := &model.Role{
		Name:       request.GetName(),
		PoolID:     request.GetPoolId(),
		Attributes: string(attrJSON),
	}

	if err := as.CreateRole(ctx, role); err != nil {
		return nil, err
	}

	return &v1.CreateRoleResponse{
		Role: &v1.Role{
			Name: request.Name,
		},
	}, nil
}

func (r *RoleService) GetRole(ctx context.Context, request *v1.GetRoleRequest) (*v1.GetRoleResponse, error) {
	as, err := store.GetProjectStore(ctx, r.store)
	if err != nil {
		return nil, err
	}

	poolID, err := uuid.Parse(request.GetPoolId())
	if err != nil {
		return nil, err
	}

	role, err := as.GetRole(ctx, poolID, request.GetRoleName())
	if err != nil {
		return nil, err
	}

	return &v1.GetRoleResponse{
		Role: &v1.Role{
			Name: role.Name,
		},
	}, nil
}

func (r *RoleService) ListRoles(ctx context.Context, request *v1.ListRolesRequest) (*v1.ListRolesResponse, error) {
	as, err := store.GetProjectStore(ctx, r.store)
	if err != nil {
		return nil, err
	}

	poolID, err := uuid.Parse(request.GetPoolId())
	if err != nil {
		return nil, err
	}

	page := x.GetPageFromRequest(request)

	roles, total, err := as.ListRoles(ctx, poolID, int(page.Page), int(page.Size))
	if err != nil {
		return nil, err
	}

	var listRoles []*v1.Role
	for _, role := range roles {
		listRoles = append(listRoles, &v1.Role{
			Name: role.Name,
		})
	}

	return &v1.ListRolesResponse{
		Roles: listRoles,
		Meta: &v1.Meta{
			Total: int32(total),
			Page:  int32(page.Page),
			Size:  int32(page.Size),
		},
	}, nil

}

func (r *RoleService) UpdateRole(ctx context.Context, request *v1.UpdateRoleRequest) (*v1.UpdateRoleResponse, error) {
	as, err := store.GetProjectStore(ctx, r.store)
	if err != nil {
		return nil, err
	}

	poolID, err := uuid.Parse(request.GetPoolId())
	if err != nil {
		return nil, err
	}

	err = as.Transaction(func(tx store.AuthBaseStore) error {
		role, err := tx.GetRole(ctx, poolID, request.GetRoleName())
		if err != nil {
			return err
		}

		if request.Attributes != nil {
			attributes := request.GetAttributes()
			attrJSON, err := json.Marshal(attributes)
			if err != nil {
				return err
			}
			role.Attributes = string(attrJSON)
		}

		if err := tx.UpdateRole(ctx, role); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.UpdateRoleResponse{
		Role: &v1.Role{
			Name: request.GetRoleName(),
		},
	}, nil
}

func (r *RoleService) DeleteRole(ctx context.Context, request *v1.DeleteRoleRequest) (*v1.DeleteRoleResponse, error) {
	as, err := store.GetProjectStore(ctx, r.store)
	if err != nil {
		return nil, err
	}

	poolID, err := uuid.Parse(request.GetPoolId())
	if err != nil {
		return nil, err
	}

	if err := as.DeleteRole(ctx, poolID, request.GetRoleName()); err != nil {
		return nil, err
	}

	return &v1.DeleteRoleResponse{}, nil
}
