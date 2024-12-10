package service

import (
	"context"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/cache"
	"github.com/emrgen/authbase/pkg/model"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/google/uuid"
)

var _ v1.MemberServiceServer = new(MemberService)

type MemberService struct {
	store store.AuthBaseStore
	cache *cache.Redis
	v1.UnimplementedMemberServiceServer
}

func NewMemberService(store store.AuthBaseStore, cache *cache.Redis) *MemberService {
	return &MemberService{store: store, cache: cache}
}

func (m *MemberService) CreateMember(ctx context.Context, request *v1.CreateMemberRequest) (*v1.CreateMemberResponse, error) {
	member := model.User{
		ID:             uuid.New().String(),
		OrganizationID: request.GetOrganizationId(),
		Username:       request.GetUsername(),
		Email:          request.GetEmail(),
		Member:         true,
	}

	permission := model.Permission{
		OrganizationID: request.GetOrganizationId(),
		UserID:         member.ID,
	}

	err := m.store.Transaction(func(tx store.AuthBaseStore) error {

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

func (m *MemberService) GetMember(ctx context.Context, request *v1.GetMemberRequest) (*v1.GetMemberResponse, error) {
	id, err := uuid.Parse(request.GetMemberId())
	if err != nil {
		return nil, err
	}

	member, err := m.store.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &v1.GetMemberResponse{
		Member: &v1.Member{
			Id:       member.ID,
			Username: member.Username,
		},
	}, nil
}

func (m *MemberService) ListMember(ctx context.Context, request *v1.ListMemberRequest) (*v1.ListMemberResponse, error) {
	orgID, err := uuid.Parse(request.GetOrganizationId())
	if err != nil {
		return nil, err
	}

	page := v1.Page{
		Page: 0,
		Size: 20,
	}
	if request.GetPage() != nil {
		page.Page = request.GetPage().Page
		page.Size = request.GetPage().Size
	}

	members, total, err := m.store.ListUsersByOrg(ctx, true, orgID, int(page.Page), int(page.Size))
	if err != nil {
		return nil, err
	}

	var memberList []*v1.Member
	for _, member := range members {
		memberList = append(memberList, &v1.Member{
			Id:       member.ID,
			Username: member.Username,
		})
	}

	return &v1.ListMemberResponse{
		Members: memberList,
		Meta:    &v1.Meta{Total: int32(total)},
	}, nil
}

func (m *MemberService) UpdateMember(ctx context.Context, request *v1.UpdateMemberRequest) (*v1.UpdateMemberResponse, error) {
	memberID, err := uuid.Parse(request.GetMemberId())
	if err != nil {
		return nil, err
	}
	orgID, err := uuid.Parse(request.GetOrganizationId())
	if err != nil {
		return nil, err
	}

	// update the member and the permission
	err = m.store.Transaction(func(tx store.AuthBaseStore) error {
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

		permission, err := tx.GetPermissionByID(ctx, orgID, memberID)

		if request.GetPermissions() != nil {
			permission.Permission = 0
			for _, p := range request.GetPermissions() {
				permission.Permission |= uint8(p.Number())
			}
		}

		err = tx.UpdateUser(ctx, member)
		if err != nil {
			return err
		}

		err = tx.UpdatePermission(ctx, permission)
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

func (m *MemberService) DeleteMember(ctx context.Context, request *v1.DeleteMemberRequest) (*v1.DeleteMemberResponse, error) {
	memberID, err := uuid.Parse(request.GetMemberId())
	if err != nil {
		return nil, err
	}

	orgID, err := uuid.Parse(request.GetOrganizationId())
	if err != nil {
		return nil, err
	}

	err = m.store.Transaction(func(tx store.AuthBaseStore) error {
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
