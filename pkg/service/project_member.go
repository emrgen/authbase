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

var _ v1.ProjectMemberServiceServer = new(ProjectMemberService)

type ProjectMemberService struct {
	perm  permission.AuthBasePermission
	store store.Provider
	cache *cache.Redis
	v1.UnimplementedProjectMemberServiceServer
}

func NewProjectMemberService(perm permission.AuthBasePermission, store store.Provider, cache *cache.Redis) *ProjectMemberService {
	return &ProjectMemberService{perm: perm, store: store, cache: cache}
}

// CreateProjectMember creates a member of an project
func (m *ProjectMemberService) CreateProjectMember(ctx context.Context, request *v1.CreateProjectMemberRequest) (*v1.CreateProjectMemberResponse, error) {
	as, err := store.GetProjectStore(ctx, m.store)
	if err != nil {
		return nil, err
	}

	member := model.Account{
		ID:            uuid.New().String(),
		ProjectID:     request.GetProjectId(),
		Username:      request.GetUsername(),
		VisibleName:   request.GetVisibleName(),
		Email:         request.GetEmail(),
		ProjectMember: true,
	}

	// create a perm for the new member
	perm := model.ProjectMember{
		ProjectID:  request.GetProjectId(),
		AccountID:  member.ID,
		Permission: uint32(request.GetPermission()),
	}

	// if the user already exists, return an error
	err = as.Transaction(func(tx store.AuthBaseStore) error {
		if err := tx.CreateAccount(ctx, &member); err != nil {
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

	return &v1.CreateProjectMemberResponse{
		Id: member.ID,
	}, nil
}

// GetProjectMember gets a member by ID of an project
func (m *ProjectMemberService) GetProjectMember(ctx context.Context, request *v1.GetProjectMemberRequest) (*v1.GetProjectMemberResponse, error) {
	as, err := store.GetProjectStore(ctx, m.store)
	if err != nil {
		return nil, err
	}
	id, err := uuid.Parse(request.GetMemberId())
	if err != nil {
		return nil, err
	}

	member, err := as.GetAccountByID(ctx, id)
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
	return &v1.GetProjectMemberResponse{
		ProjectMember: &v1.ProjectMember{
			Id:         member.ID,
			Username:   member.Username,
			Permission: v1.Permission(perm.Permission),
		},
	}, nil
}

// ListProjectMember lists members of an project
func (m *ProjectMemberService) ListProjectMember(ctx context.Context, request *v1.ListProjectMemberRequest) (*v1.ListProjectMemberResponse, error) {
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
	members, total, err := as.ListProjectAccounts(ctx, true, orgID, int(page.Page), int(page.Size))
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

	permissions, err := as.ListProjectMembersByAccountIDs(ctx, orgID, userIDs)
	permissionMap := make(map[string]uint32)
	for _, perm := range permissions {
		permissionMap[perm.AccountID] = perm.Permission
	}

	var memberList []*v1.ProjectMember
	for _, member := range members {
		perm := permissionMap[member.ID]
		if perm == 0 {
			continue
		}

		memberList = append(memberList, &v1.ProjectMember{
			Id:         member.ID,
			Username:   member.Username,
			Email:      member.Email,
			Permission: v1.Permission(perm),
		})
	}

	return &v1.ListProjectMemberResponse{
		Members: memberList,
		Meta:    &v1.Meta{Total: int32(total)},
	}, nil
}

// UpdateProjectMember updates a member of an project
func (m *ProjectMemberService) UpdateProjectMember(ctx context.Context, request *v1.UpdateProjectMemberRequest) (*v1.UpdateProjectMemberResponse, error) {
	var err error
	userID, err := uuid.Parse(request.GetMemberId())
	if err != nil {
		return nil, err
	}

	as, err := store.GetProjectStore(ctx, m.store)
	if err != nil {
		return nil, err
	}

	member, err := as.GetAccountByID(ctx, userID)
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

		err = tx.UpdateAccount(ctx, member)
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

	return &v1.UpdateProjectMemberResponse{
		Id:      member.ID,
		Message: "ProjectMember updated successfully.",
	}, nil
}

// AddProjectMember makes a user a member of an organization.
func (m *ProjectMemberService) AddProjectMember(ctx context.Context, request *v1.AddProjectMemberRequest) (*v1.AddProjectMemberResponse, error) {
	logrus.Info("AddProjectMember")
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

	user, err := as.GetAccountByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	user.ProjectMember = true

	pool, err := as.GetPoolByID(ctx, uuid.MustParse(user.PoolID))
	if err != nil {
		return nil, err
	}

	if !pool.Default || pool.ProjectID != orgID.String() {
		return nil, errors.New("user does not belong to the project default pool")
	}

	perm, err := as.GetProjectMemberByID(ctx, orgID, userID)
	if errors.Is(err, store.ErrPermissionNotFound) {
		perm = &model.ProjectMember{
			ProjectID: orgID.String(),
			AccountID: userID.String(),
		}
	} else if err != nil {
		return nil, nil
	}

	permValue := request.GetPermission()
	if permValue != v1.Permission_UNKNOWN {
		perm.Permission = uint32(permValue)
	}

	err = as.Transaction(func(tx store.AuthBaseStore) error {
		err = tx.UpdateAccount(ctx, user)
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

	return &v1.AddProjectMemberResponse{
		Message: "ProjectMember added successfully",
	}, nil
}

// RemoveProjectMember removes a member from an organization
func (m *ProjectMemberService) RemoveProjectMember(ctx context.Context, request *v1.RemoveProjectMemberRequest) (*v1.RemoveProjectMemberResponse, error) {
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

	err = as.Transaction(func(tx store.AuthBaseStore) error {
		// get the user member
		member, err := tx.GetProjectMemberByID(ctx, orgID, memberID)
		if member.Permission == uint32(v1.Permission_OWNER) {
			return errors.New("cannot remove an owner")
		}

		user, err := as.GetAccountByID(ctx, memberID)
		if err != nil {
			return err
		}
		user.ProjectMember = false

		err = tx.UpdateAccount(ctx, user)
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

	return &v1.RemoveProjectMemberResponse{
		Message: "user member removed successfully",
	}, nil
}
