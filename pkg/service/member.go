package service

import (
	"context"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/cache"
	"github.com/emrgen/authbase/pkg/store"
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
	//TODO implement me
	panic("implement me")
}

func (m *MemberService) GetMember(ctx context.Context, request *v1.GetMemberRequest) (*v1.GetMemberResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MemberService) ListMember(ctx context.Context, request *v1.ListMemberRequest) (*v1.ListMemberResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MemberService) UpdateMember(ctx context.Context, request *v1.UpdateMemberRequest) (*v1.UpdateMemberResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MemberService) DeleteMember(ctx context.Context, request *v1.DeleteMemberRequest) (*v1.DeleteMemberResponse, error) {
	//TODO implement me
	panic("implement me")
}
