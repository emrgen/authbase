package service

import (
	"context"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/model"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/emrgen/authbase/x"
	"github.com/google/uuid"
)

// NewPoolMemberService creates a new pool member service.
func NewPoolMemberService(store store.Provider) v1.PoolMemberServiceServer {
	return &PoolMemberService{
		store: store,
	}
}

var _ v1.PoolMemberServiceServer = new(PoolMemberService)

type PoolMemberService struct {
	store store.Provider
	v1.UnimplementedPoolMemberServiceServer
}

func (p *PoolMemberService) CreatePoolMember(ctx context.Context, request *v1.CreatePoolMemberRequest) (*v1.CreatePoolMemberResponse, error) {
	poolID := uuid.MustParse(request.GetPoolId())
	accountID := uuid.MustParse(request.GetAccountId())
	permission := request.GetPermission()

	member := model.PoolMember{
		AccountID:  accountID.String(),
		PoolID:     poolID.String(),
		Permission: uint32(permission),
	}

	as, err := store.GetProjectStore(ctx, p.store)
	if err != nil {
		return nil, err
	}

	err = as.AddPoolMember(ctx, &member)
	if err != nil {
		return nil, err
	}

	return &v1.CreatePoolMemberResponse{
		PoolMember: &v1.PoolMember{
			PoolId:     member.PoolID,
			AccountId:  member.AccountID,
			Permission: permission,
		},
	}, nil
}

func (p *PoolMemberService) GetPoolMember(ctx context.Context, request *v1.GetPoolMemberRequest) (*v1.GetPoolMemberResponse, error) {
	poolID := uuid.MustParse(request.GetPoolId())
	accountID := uuid.MustParse(request.GetAccountId())

	as, err := store.GetProjectStore(ctx, p.store)
	if err != nil {
		return nil, err
	}

	member, err := as.GetPoolMember(ctx, poolID, accountID)
	if err != nil {
		return nil, err
	}

	return &v1.GetPoolMemberResponse{
		PoolMember: &v1.PoolMember{
			PoolId:    member.PoolID,
			AccountId: member.AccountID,
		},
	}, nil
}

func (p *PoolMemberService) ListPoolMembers(ctx context.Context, request *v1.ListPoolMembersRequest) (*v1.ListPoolMembersResponse, error) {
	poolID := uuid.MustParse(request.GetPoolId())
	page := x.GetPageFromRequest(request)
	as, err := store.GetProjectStore(ctx, p.store)
	if err != nil {
		return nil, err
	}
	members, total, err := as.ListPoolMembers(ctx, poolID, int(page.Page), int(page.Size))
	if err != nil {
		return nil, err
	}

	var pbMembers []*v1.PoolMember
	for _, member := range members {
		pbMembers = append(pbMembers, &v1.PoolMember{
			PoolId:    member.PoolID,
			AccountId: member.AccountID,
		})
	}

	return &v1.ListPoolMembersResponse{
		PoolMembers: pbMembers,
		Meta: &v1.Meta{
			Total: int32(total),
			Page:  page.Page,
			Size:  page.Size,
		},
	}, nil

}

func (p *PoolMemberService) UpdatePoolMember(ctx context.Context, request *v1.UpdatePoolMemberRequest) (*v1.UpdatePoolMemberResponse, error) {
	poolID := uuid.MustParse(request.GetPoolId())
	accountID := uuid.MustParse(request.GetAccountId())
	permission := request.GetPermission()

	as, err := store.GetProjectStore(ctx, p.store)
	if err != nil {
		return nil, err
	}

	err = as.Transaction(func(tx store.AuthBaseStore) error {
		member, err := tx.GetPoolMember(ctx, poolID, accountID)
		if err != nil {
			return err
		}

		member.Permission = uint32(permission)
		return tx.UpdatePoolMember(ctx, member)
	})
	if err != nil {
		return nil, err
	}

	return &v1.UpdatePoolMemberResponse{}, nil
}

func (p *PoolMemberService) DeletePoolMember(ctx context.Context, request *v1.DeletePoolMemberRequest) (*v1.DeletePoolMemberResponse, error) {
	poolID := uuid.MustParse(request.GetPoolId())
	accountID := uuid.MustParse(request.GetAccountId())
	as, err := store.GetProjectStore(ctx, p.store)
	if err != nil {
		return nil, err
	}

	err = as.RemovePoolMember(ctx, poolID, accountID)
	if err != nil {
		return nil, err
	}

	return &v1.DeletePoolMemberResponse{}, nil
}
