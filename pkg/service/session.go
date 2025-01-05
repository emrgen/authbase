package service

import (
	"context"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/permission"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/google/uuid"
)

func NewSessionService(store store.Provider, perm permission.AuthBasePermission) *SessionService {
	return &SessionService{store: store, perm: perm}
}

var _ v1.SessionServiceServer = (*SessionService)(nil)

// SessionService implements the v1.SessionServiceServer interface.
type SessionService struct {
	store store.Provider
	perm  permission.AuthBasePermission
	v1.UnimplementedSessionServiceServer
}

func (s *SessionService) ListUserSession(ctx context.Context, request *v1.ListUserSessionRequest) (*v1.ListUserSessionResponse, error) {
	//TODO implement me
	panic("implement me")
}

// DeleteSession logs out a user from a session
func (s *SessionService) DeleteSession(ctx context.Context, request *v1.DeleteSessionsRequest) (*v1.DeleteSessionsResponse, error) {
	as, err := store.GetProjectStore(ctx, s.store)
	if err != nil {
		return nil, err
	}

	sessionID, err := uuid.Parse(request.GetSessionId())
	if err != nil {
		return nil, err
	}

	user, err := as.GetUserByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	err = s.perm.CheckProjectPermission(ctx, uuid.MustParse(user.ProjectID), "write")
	if err != nil {
		return nil, err
	}

	err = as.DeleteSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	return &v1.DeleteSessionsResponse{
		Message: "logged out",
	}, nil
}

// DeleteAllSessions logs out all sessions of a user
func (s *SessionService) DeleteAllSessions(ctx context.Context, request *v1.DeleteAllSessionsRequest) (*v1.DeleteAllSessionsResponse, error) {
	var err error

	userID := uuid.MustParse(request.GetUserId())
	as, err := store.GetProjectStore(ctx, s.store)
	if err != nil {
		return nil, err
	}

	// user project id
	user, err := as.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	orgID := uuid.MustParse(user.ProjectID)

	err = s.perm.CheckProjectPermission(ctx, orgID, "write")
	if err != nil {
		return nil, err
	}

	err = as.DeleteSessionByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &v1.DeleteAllSessionsResponse{
		Message: "logged out of all sessions",
	}, nil
}
