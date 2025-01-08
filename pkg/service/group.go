package service

import (
	"context"
	"errors"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/model"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/emrgen/authbase/x"
	"github.com/google/uuid"
)

func NewGroupService(store store.Provider) *GroupService {
	return &GroupService{
		store: store,
	}
}

var _ v1.GroupServiceServer = (*GroupService)(nil)

type GroupService struct {
	store store.Provider
	v1.UnimplementedGroupServiceServer
}

func (g *GroupService) CreateGroup(ctx context.Context, request *v1.CreateGroupRequest) (*v1.CreateGroupResponse, error) {
	as, err := store.GetProjectStore(ctx, g.store)
	if err != nil {
		return nil, err
	}

	name := request.GetName()
	poolID := request.GetPoolId()
	rolesNames := request.GetRoleNames()
	roles := make([]*model.Role, 0)

	// Create roles if not exist.
	for _, roleName := range rolesNames {
		role := &model.Role{
			Name:   roleName,
			PoolID: poolID,
		}
		err = as.CreateRole(ctx, role)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	group := &model.Group{
		ID:     uuid.New().String(),
		Name:   name,
		PoolID: poolID,
		Roles:  roles,
	}

	err = as.CreateGroup(ctx, group)
	if err != nil {
		return nil, err
	}

	roleProtos := make([]*v1.Role, 0)
	for _, role := range roles {
		roleProtos = append(roleProtos, &v1.Role{
			Name: role.Name,
		})
	}

	return &v1.CreateGroupResponse{
		Group: &v1.Group{
			Id:     group.ID,
			Name:   name,
			PoolId: poolID,
			Roles:  roleProtos,
		},
	}, nil
}

func (g *GroupService) GetGroup(ctx context.Context, request *v1.GetGroupRequest) (*v1.GetGroupResponse, error) {
	as, err := store.GetProjectStore(ctx, g.store)
	if err != nil {
		return nil, err
	}

	groupID := uuid.MustParse(request.GetGroupId())
	group, err := as.GetGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}

	roles := make([]*v1.Role, 0)
	for _, role := range group.Roles {
		roles = append(roles, &v1.Role{
			Name: role.Name,
		})
	}

	return &v1.GetGroupResponse{
		Group: &v1.Group{
			Id:     group.ID,
			Name:   group.Name,
			PoolId: group.PoolID,
			Roles:  roles,
		},
	}, nil
}

// ListGroups lists groups in a pool or groups that an account is a member of.
// when both pool_id and account_id are provided, list groups will be filtered by pool_id.
func (g *GroupService) ListGroups(ctx context.Context, request *v1.ListGroupsRequest) (*v1.ListGroupsResponse, error) {
	var err error
	if request.GetPoolId() == "" && request.GetAccountId() == "" {
		return nil, errors.New("pool_id or account_id is required")
	}

	as, err := store.GetProjectStore(ctx, g.store)
	if err != nil {
		return nil, err
	}

	page := x.GetPageFromRequest(request)
	var groups []*model.Group
	var total int

	// If account_id is provided, list groups that the account is a member of.
	if request.GetAccountId() != "" {
		accountID := uuid.MustParse(request.GetAccountId())
		memberships, err := as.ListGroupMemberByAccount(ctx, accountID)
		if err != nil {
			return nil, err
		}
		for _, membership := range memberships {
			groups = append(groups, membership.Group)
		}
	}

	// If pool_id is provided, list groups in the pool.
	if request.GetPoolId() != "" {
		poolID := uuid.MustParse(request.GetPoolId())
		groups, total, err = as.ListGroups(ctx, poolID, int(page.Page), int(page.Size))
		if err != nil {
			return nil, err
		}
	}

	var responseGroups []*v1.Group
	for _, group := range groups {
		roles := make([]*v1.Role, 0)
		for _, role := range group.Roles {
			roles = append(roles, &v1.Role{
				Name: role.Name,
			})
		}
		responseGroups = append(responseGroups, &v1.Group{
			Id:     group.ID,
			Name:   group.Name,
			PoolId: group.PoolID,
			Roles:  roles,
		})
	}

	return &v1.ListGroupsResponse{
		Groups: responseGroups,
		Meta: &v1.Meta{
			Total: int32(total),
		},
	}, nil

}

func (g *GroupService) UpdateGroup(ctx context.Context, request *v1.UpdateGroupRequest) (*v1.UpdateGroupResponse, error) {
	var err error
	as, err := store.GetProjectStore(ctx, g.store)
	if err != nil {
		return nil, err
	}

	groupID := uuid.MustParse(request.GetGroupId())
	roleNames := request.GetRoleNames()

	err = as.Transaction(func(tx store.AuthBaseStore) error {
		group, err := tx.GetGroup(ctx, groupID)
		if err != nil {
			return err
		}
		poolID, err := uuid.Parse(group.PoolID)
		if err != nil {
			return err
		}

		roles, err := as.ListRolesByNames(ctx, poolID, roleNames)
		if err != nil {
			return err
		}

		if request.GetName() != "" {
			group.Name = request.GetName()
		}
		group.Roles = roles
		err = tx.UpdateGroup(ctx, group)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.UpdateGroupResponse{}, nil
}

func (g *GroupService) DeleteGroup(ctx context.Context, request *v1.DeleteGroupRequest) (*v1.DeleteGroupResponse, error) {
	as, err := store.GetProjectStore(ctx, g.store)
	if err != nil {
		return nil, err
	}

	groupID := uuid.MustParse(request.GetGroupId())
	err = as.DeleteGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}

	return &v1.DeleteGroupResponse{}, nil
}

func (g *GroupService) AddGroupMember(ctx context.Context, request *v1.AddGroupMemberRequest) (*v1.AddGroupMemberResponse, error) {
	as, err := store.GetProjectStore(ctx, g.store)
	if err != nil {
		return nil, err
	}

	groupID := uuid.MustParse(request.GetGroupId())
	accountID := uuid.MustParse(request.GetAccountId())
	groupMember := &model.GroupMember{
		GroupID:   groupID.String(),
		AccountID: accountID.String(),
	}

	err = as.AddGroupMember(ctx, groupMember)
	if err != nil {
		return nil, err
	}

	return &v1.AddGroupMemberResponse{}, nil
}

func (g *GroupService) RemoveGroupMember(ctx context.Context, request *v1.RemoveGroupMemberRequest) (*v1.RemoveGroupMemberResponse, error) {
	as, err := store.GetProjectStore(ctx, g.store)
	if err != nil {
		return nil, err
	}

	groupID := uuid.MustParse(request.GetGroupId())
	accountID := uuid.MustParse(request.GetAccountId())

	err = as.RemoveGroupMember(ctx, groupID, accountID)
	if err != nil {
		return nil, err
	}

	return &v1.RemoveGroupMemberResponse{}, nil
}

func (g *GroupService) ListGroupMembers(ctx context.Context, request *v1.ListGroupMembersRequest) (*v1.ListGroupMembersResponse, error) {
	as, err := store.GetProjectStore(ctx, g.store)
	if err != nil {
		return nil, err
	}

	groupID := uuid.MustParse(request.GetGroupId())
	page := x.GetPageFromRequest(request)

	group, err := as.GetGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}

	groupMembers, total, err := as.ListGroupMembers(ctx, groupID, int(page.Page), int(page.Size))
	if err != nil {
		return nil, err
	}

	var responseGroupMembers []*v1.Account
	for _, member := range groupMembers {
		responseGroupMembers = append(responseGroupMembers, &v1.Account{
			Id:          member.AccountID,
			VisibleName: member.Account.VisibleName,
			Email:       member.Account.Email,
		})
	}

	roles := make([]*v1.Role, 0)
	for _, role := range group.Roles {
		roles = append(roles, &v1.Role{
			Name: role.Name,
		})
	}

	return &v1.ListGroupMembersResponse{
		Members: responseGroupMembers,
		Roles:   roles,
		Meta: &v1.Meta{
			Total: int32(total),
			Page:  page.Page,
			Size:  page.Size,
		},
	}, nil
}

func (g *GroupService) AddRole(ctx context.Context, request *v1.AddRoleRequest) (*v1.AddRoleResponse, error) {
	roleName := request.GetRoleName()
	groupID := uuid.MustParse(request.GetGroupId())
	as, err := store.GetProjectStore(ctx, g.store)
	if err != nil {
		return nil, err
	}

	// Add role to group.
	err = as.Transaction(func(tx store.AuthBaseStore) error {
		group, err := tx.GetGroup(ctx, groupID)
		if err != nil {
			return err
		}

		poolID, err := uuid.Parse(group.PoolID)
		if err != nil {
			return err
		}

		role, err := tx.GetRole(ctx, poolID, roleName)
		if err != nil {
			return err
		}
		roles := make(map[string]*model.Role)
		roles[roleName] = role
		for _, role := range group.Roles {
			roles[role.Name] = role
		}

		group.Roles = make([]*model.Role, 0)
		for _, role := range roles {
			group.Roles = append(group.Roles, role)
		}

		err = tx.UpdateGroup(ctx, group)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.AddRoleResponse{}, nil
}

func (g *GroupService) RemoveRole(ctx context.Context, request *v1.RemoveRoleRequest) (*v1.RemoveRoleResponse, error) {
	roleName := request.GetRoleName()
	groupID := uuid.MustParse(request.GetGroupId())
	as, err := store.GetProjectStore(ctx, g.store)
	if err != nil {
		return nil, err
	}

	// Remove role from group.
	err = as.Transaction(func(tx store.AuthBaseStore) error {
		group, err := tx.GetGroup(ctx, groupID)
		if err != nil {
			return err
		}

		roles := group.Roles

		group.Roles = make([]*model.Role, 0)
		for _, role := range roles {
			if role.Name != roleName {
				group.Roles = append(group.Roles, role)
			}
		}

		err = tx.UpdateGroup(ctx, group)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.RemoveRoleResponse{}, nil
}
