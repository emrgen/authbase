package service

import (
	"context"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/model"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/emrgen/authbase/x"
	"github.com/google/uuid"
	"strings"
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
	scopes := request.GetScopes()

	group := &model.Group{
		ID:     uuid.New().String(),
		Name:   name,
		PoolID: poolID,
		Scopes: strings.Join(scopes, ","),
	}

	err = as.CreateGroup(ctx, group)
	if err != nil {
		return nil, err
	}

	return &v1.CreateGroupResponse{
		Group: &v1.Group{
			Id:     group.ID,
			Name:   name,
			PoolId: poolID,
			Scopes: scopes,
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

	return &v1.GetGroupResponse{
		Group: &v1.Group{
			Id:     group.ID,
			Name:   group.Name,
			PoolId: group.PoolID,
			Scopes: strings.Split(group.Scopes, ","),
		},
	}, nil
}

func (g *GroupService) ListGroups(ctx context.Context, request *v1.ListGroupsRequest) (*v1.ListGroupsResponse, error) {
	as, err := store.GetProjectStore(ctx, g.store)
	if err != nil {
		return nil, err
	}

	poolID := uuid.MustParse(request.GetPoolId())
	page := x.GetPageFromRequest(request)

	groups, total, err := as.ListGroups(ctx, poolID, int(page.Page), int(page.Size))
	if err != nil {
		return nil, err
	}

	var responseGroups []*v1.Group
	for _, group := range groups {
		responseGroups = append(responseGroups, &v1.Group{
			Id:     group.ID,
			Name:   group.Name,
			PoolId: group.PoolID,
			Scopes: strings.Split(group.Scopes, ","),
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
	scopes := request.GetScopes()
	err = as.Transaction(func(tx store.AuthBaseStore) error {
		group, err := tx.GetGroup(ctx, groupID)
		if err != nil {
			return err
		}

		if request.GetName() != "" {
			group.Name = request.GetName()
		}
		group.Scopes = strings.Join(scopes, ",")

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

	return &v1.ListGroupMembersResponse{
		Members: responseGroupMembers,
		Scopes:  strings.Split(group.Scopes, ","),
		Meta: &v1.Meta{
			Total: int32(total),
			Page:  page.Page,
			Size:  page.Size,
		},
	}, nil
}
