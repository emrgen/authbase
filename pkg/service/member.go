package service

import (
	"context"
	"errors"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/cache"
	"github.com/emrgen/authbase/pkg/model"
	"github.com/emrgen/authbase/pkg/permission"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/emrgen/authbase/x"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var _ v1.MemberServiceServer = new(MemberService)

type MemberService struct {
	perm  permission.AuthBasePermission
	store store.Provider
	cache *cache.Redis
	v1.UnimplementedMemberServiceServer
}

func NewMemberService(perm permission.AuthBasePermission, store store.Provider, cache *cache.Redis) *MemberService {
	return &MemberService{perm: perm, store: store, cache: cache}
}

// CreateMember creates a member of an project
func (m *MemberService) CreateMember(ctx context.Context, request *v1.CreateMemberRequest) (*v1.CreateMemberResponse, error) {
	as, err := store.GetProjectStore(ctx, m.store)
	if err != nil {
		return nil, err
	}

	member := model.User{
		ID:        uuid.New().String(),
		ProjectID: request.GetProjectId(),
		Username:  request.GetUsername(),
		Email:     request.GetEmail(),
		Member:    true,
	}

	// create a perm for the new member
	perm := model.ProjectMember{
		ProjectID:  request.GetProjectId(),
		UserID:     member.ID,
		Permission: uint32(request.GetPermission()),
	}

	// if the user already exists, return an error
	err = as.Transaction(func(tx store.AuthBaseStore) error {
		if err := tx.CreateUser(ctx, &member); err != nil {
			return err
		}

		if err := tx.CreateProjectMember(ctx, &perm); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.CreateMemberResponse{
		Id: member.ID,
	}, nil
}

// GetMember gets a member by ID of an project
func (m *MemberService) GetMember(ctx context.Context, request *v1.GetMemberRequest) (*v1.GetMemberResponse, error) {
	as, err := store.GetProjectStore(ctx, m.store)
	if err != nil {
		return nil, err
	}
	id, err := uuid.Parse(request.GetMemberId())
	if err != nil {
		return nil, err
	}

	member, err := as.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	orgID, err := uuid.Parse(member.ProjectID)
	if err != nil {
		return nil, err
	}

	// check if the user has the read permission
	err = m.perm.CheckProjectPermission(ctx, orgID, "read")
	if err != nil {
		return nil, err
	}

	perm, err := as.GetProjectMemberByID(ctx, orgID, id)
	return &v1.GetMemberResponse{
		Member: &v1.Member{
			Id:         member.ID,
			Username:   member.Username,
			Permission: v1.Permission(perm.Permission),
		},
	}, nil
}

// ListMember lists members of an project
func (m *MemberService) ListMember(ctx context.Context, request *v1.ListMemberRequest) (*v1.ListMemberResponse, error) {
	var err error
	orgID, err := uuid.Parse(request.GetProjectId())
	if err != nil {
		return nil, err
	}

	err = m.perm.CheckProjectPermission(ctx, orgID, "read")
	if err != nil {
		return nil, err
	}

	as, err := store.GetProjectStore(ctx, m.store)
	if err != nil {
		return nil, err
	}

	page := x.GetPageFromRequest(request)
	members, total, err := as.ListUsersByOrg(ctx, true, orgID, int(page.Page), int(page.Size))
	if err != nil {
		return nil, err
	}

	userIDs := make([]uuid.UUID, 0)
	for _, member := range members {
		id, err := uuid.Parse(member.ID)
		if err != nil {
			return nil, err
		}
		userIDs = append(userIDs, id)
	}

	permissions, err := as.ListProjectMembersUsers(ctx, orgID, userIDs)
	permissionMap := make(map[string]uint32)
	for _, perm := range permissions {
		permissionMap[perm.UserID] = perm.Permission
	}

	var memberList []*v1.Member
	for _, member := range members {
		perm := permissionMap[member.ID]
		if perm == 0 {
			continue
		}

		memberList = append(memberList, &v1.Member{
			Id:         member.ID,
			Username:   member.Username,
			Email:      member.Email,
			Permission: v1.Permission(perm),
		})
	}

	return &v1.ListMemberResponse{
		Members: memberList,
		Meta:    &v1.Meta{Total: int32(total)},
	}, nil
}

// UpdateMember updates a member of an project
func (m *MemberService) UpdateMember(ctx context.Context, request *v1.UpdateMemberRequest) (*v1.UpdateMemberResponse, error) {
	var err error
	userID, err := uuid.Parse(request.GetMemberId())
	if err != nil {
		return nil, err
	}

	as, err := store.GetProjectStore(ctx, m.store)
	if err != nil {
		return nil, err
	}

	member, err := as.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	orgID, err := uuid.Parse(member.ProjectID)
	if err != nil {
		return nil, err
	}

	err = m.perm.CheckProjectPermission(ctx, orgID, "write")
	if err != nil {
		return nil, err
	}

	// update the member and the permission
	err = as.Transaction(func(tx store.AuthBaseStore) error {
		perm, err := tx.GetProjectMemberByID(ctx, orgID, userID)
		if err != nil {
			return err
		}

		if request.GetPermission() != v1.Permission_UNKNOWN {
			perm.Permission = uint32(request.GetPermission())
		}

		err = tx.UpdateUser(ctx, member)
		if err != nil {
			return err
		}

		err = tx.UpdateProjectMember(ctx, perm)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.UpdateMemberResponse{
		Id:      member.ID,
		Message: "Member updated successfully.",
	}, nil
}

// AddMember makes a user a member of an organization.
func (m *MemberService) AddMember(ctx context.Context, request *v1.AddMemberRequest) (*v1.AddMemberResponse, error) {
	logrus.Info("AddMember")
	as, err := store.GetProjectStore(ctx, m.store)
	if err != nil {
		return nil, err
	}

	orgID, err := uuid.Parse(request.GetProjectId())
	if err != nil {
		return nil, err
	}

	err = m.perm.CheckProjectPermission(ctx, orgID, "write")
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
	user.Member = true

	err = m.perm.CheckProjectPermission(ctx, uuid.MustParse(user.ProjectID), "write")
	if err != nil {
		return nil, err
	}

	perm, err := as.GetProjectMemberByID(ctx, orgID, userID)
	if errors.Is(err, store.ErrPermissionNotFound) {
		perm = &model.ProjectMember{
			ProjectID: orgID.String(),
			UserID:    userID.String(),
		}
	} else if err != nil {
		return nil, nil
	}

	permValue := request.GetPermission()
	if permValue != v1.Permission_UNKNOWN {
		perm.Permission = uint32(permValue)
	}

	err = as.Transaction(func(tx store.AuthBaseStore) error {
		err = tx.UpdateUser(ctx, user)
		if err != nil {
			return err
		}

		err = tx.CreateProjectMember(ctx, perm)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.AddMemberResponse{
		Message: "Member added successfully",
	}, nil
}

// RemoveMember removes a member from an organization
func (m *MemberService) RemoveMember(ctx context.Context, request *v1.RemoveMemberRequest) (*v1.RemoveMemberResponse, error) {
	orgID, err := uuid.Parse(request.GetProjectId())
	if err != nil {
		return nil, err
	}

	err = m.perm.CheckProjectPermission(ctx, orgID, "write")
	if err != nil {
		return nil, err
	}

	as, err := store.GetProjectStore(ctx, m.store)
	if err != nil {
		return nil, err
	}

	memberID, err := uuid.Parse(request.GetMemberId())
	if err != nil {
		return nil, err
	}

	user, err := as.GetUserByID(ctx, memberID)
	if err != nil {
		return nil, err
	}
	user.Member = false

	err = as.Transaction(func(tx store.AuthBaseStore) error {
		err := tx.UpdateUser(ctx, user)
		if err != nil {
			return err
		}

		// delete the permission of the user
		err = tx.DeleteProjectMember(ctx, orgID, memberID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.RemoveMemberResponse{
		Message: "user member removed successfully",
	}, nil
}
