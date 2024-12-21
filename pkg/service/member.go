package service

import (
	"context"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/cache"
	"github.com/emrgen/authbase/pkg/model"
	"github.com/emrgen/authbase/pkg/permission"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/emrgen/authbase/x"
	"github.com/google/uuid"
)

var _ v1.MemberServiceServer = new(MemberService)

type MemberService struct {
	perm  permission.AuthBasePermission
	store store.Provider
	cache *cache.Redis
	v1.UnimplementedMemberServiceServer
}

func NewMemberService(store store.Provider, cache *cache.Redis) *MemberService {
	return &MemberService{store: store, cache: cache}
}

// CreateMember creates a member of an organization
func (m *MemberService) CreateMember(ctx context.Context, request *v1.CreateMemberRequest) (*v1.CreateMemberResponse, error) {
	as, err := store.GetProjectStore(ctx, m.store)
	if err != nil {
		return nil, err
	}

	member := model.User{
		ID:             uuid.New().String(),
		OrganizationID: request.GetOrganizationId(),
		Username:       request.GetUsername(),
		Email:          request.GetEmail(),
		Member:         true,
	}

	// create a permission for the new member
	permission := model.Permission{
		OrganizationID: request.GetOrganizationId(),
		UserID:         member.ID,
	}

	// if the user already exists, return an error
	err = as.Transaction(func(tx store.AuthBaseStore) error {
		if err := tx.CreateUser(ctx, &member); err != nil {
			return err
		}

		if err := tx.CreatePermission(ctx, &permission); err != nil {
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

// GetMember gets a member by ID of an organization
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

	orgID, err := uuid.Parse(member.OrganizationID)
	if err != nil {
		return nil, err
	}

	// check if the user has the read permission
	err = m.perm.CheckOrganizationPermission(ctx, orgID, "read")
	if err != nil {
		return nil, err
	}

	perm, err := as.GetPermissionByID(ctx, orgID, id)
	permissions := make([]v1.Permission, 0)
	for value, _ := range v1.Permission_name {
		if perm.Permission&uint32(value) == 1 {
			permissions = append(permissions, v1.Permission(value))
		}
	}

	return &v1.GetMemberResponse{
		Member: &v1.Member{
			Id:          member.ID,
			Username:    member.Username,
			Permissions: permissions,
		},
	}, nil
}

// ListMember lists members of an organization
func (m *MemberService) ListMember(ctx context.Context, request *v1.ListMemberRequest) (*v1.ListMemberResponse, error) {
	var err error
	orgID, err := uuid.Parse(request.GetOrganizationId())
	if err != nil {
		return nil, err
	}

	err = m.perm.CheckOrganizationPermission(ctx, orgID, "read")
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

	permissions, err := as.ListPermissionsByUsers(ctx, orgID, userIDs)
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

		permissions := make([]v1.Permission, 0)
		for value, _ := range v1.Permission_name {
			if perm > uint32(value) {
				permissions = append(permissions, v1.Permission(value))
			}
		}

		memberList = append(memberList, &v1.Member{
			Id:          member.ID,
			Username:    member.Username,
			Permissions: permissions,
		})
	}

	return &v1.ListMemberResponse{
		Members: memberList,
		Meta:    &v1.Meta{Total: int32(total)},
	}, nil
}

// UpdateMember updates a member of an organization
func (m *MemberService) UpdateMember(ctx context.Context, request *v1.UpdateMemberRequest) (*v1.UpdateMemberResponse, error) {
	orgID, err := uuid.Parse(request.GetOrganizationId())
	if err != nil {
		return nil, err
	}

	err = m.perm.CheckOrganizationPermission(ctx, orgID, "write")
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

	// update the member and the permission
	err = as.Transaction(func(tx store.AuthBaseStore) error {
		member, err := tx.GetUserByID(ctx, memberID)
		if err != nil {
			return err
		}

		if request.GetUsername() != "" {
			member.Username = request.GetUsername()
		}

		if request.GetEmail() != "" {
			member.Email = request.GetEmail()
		}

		perm, err := tx.GetPermissionByID(ctx, orgID, memberID)

		if request.GetPermissions() != nil {
			perm.Permission = 0
			for _, p := range request.GetPermissions() {
				perm.Permission |= uint32(p.Number())
			}
		}

		err = tx.UpdateUser(ctx, member)
		if err != nil {
			return err
		}

		err = tx.UpdatePermission(ctx, perm)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.UpdateMemberResponse{
		Id:      memberID.String(),
		Message: "Member updated successfully.",
	}, nil
}

// DeleteMember deletes a member of an organization
func (m *MemberService) DeleteMember(ctx context.Context, request *v1.DeleteMemberRequest) (*v1.DeleteMemberResponse, error) {
	orgID, err := uuid.Parse(request.GetOrganizationId())
	if err != nil {
		return nil, err
	}

	err = m.perm.CheckOrganizationPermission(ctx, orgID, "write")
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

	err = as.Transaction(func(tx store.AuthBaseStore) error {
		err := tx.DeleteUser(ctx, memberID)
		if err != nil {
			return err
		}

		err = tx.DeletePermission(ctx, orgID, memberID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.DeleteMemberResponse{
		Message: "Member deleted successfully.",
	}, nil
}
