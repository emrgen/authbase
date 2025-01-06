package store

import (
	"context"
	"errors"
	"github.com/emrgen/authbase/pkg/model"
	"github.com/google/uuid"
)

var (
	ErrProjectExists           = errors.New("project already exists")
	ErrProjectNotFound         = errors.New("project not found")
	ErrPermissionNotFound      = errors.New("permission not found")
	ErrPermissionAlreadyExists = errors.New("permission already exists")
)

// AuthBaseStore is the interface for interacting with the database.
type AuthBaseStore interface {
	AccountStore
	SessionStore
	ProjectStore
	ProjectMemberStore
	ProviderStore
	RefreshTokenStore
	AccessKeyStore
	VerificationCodeStore
	ClientStore
	Migrate() error
	Transaction(func(AuthBaseStore) error) error
}

// ProjectStore is the interface for interacting with the project database.
type ProjectStore interface {
	// CreateProject creates a new project in the database.
	CreateProject(ctx context.Context, org *model.Project) error
	// GetProjectByName retrieves an project by its name.
	GetProjectByName(ctx context.Context, name string) (*model.Project, error)
	// GetProjectByID retrieves an project by its ID.
	GetProjectByID(ctx context.Context, id uuid.UUID) (*model.Project, error)
	// GetMasterProject retrieves the master project.
	GetMasterProject(ctx context.Context) (*model.Project, error)
	// ListProjects retrieves a list of projects.
	ListProjects(ctx context.Context, page, perPage int) ([]*model.Project, int, error)
	// UpdateProject updates an project in the database.
	UpdateProject(ctx context.Context, org *model.Project) error
	// DeleteProject deletes an project from the database.
	DeleteProject(ctx context.Context, id uuid.UUID) error
	// CreateKeypair creates a new keypair in the database.
	CreateKeypair(ctx context.Context, keypair *model.Keypair) error
	// GetKeypair retrieves a keypair by its ID.
	GetKeypair(ctx context.Context, id uuid.UUID) (*model.Keypair, error)
}

// AccountStore is the interface for interacting with the user database.
type AccountStore interface {
	// CreateAccount creates a new user in the database.
	CreateAccount(ctx context.Context, user *model.Account) error
	// GetAccountByEmail retrieves a user by their email address.
	GetAccountByEmail(ctx context.Context, orgID uuid.UUID, email string) (*model.Account, error)
	// GetAccountByID retrieves a user by their ID.
	GetAccountByID(ctx context.Context, id uuid.UUID) (*model.Account, error)
	// UpdateAccount updates a user in the database.
	UpdateAccount(ctx context.Context, user *model.Account) error
	// DeleteAccount deletes a user from the database.
	DeleteAccount(ctx context.Context, id uuid.UUID) error
	// ListAccountsByOrg retrieves a list of users by project.
	ListAccountsByOrg(ctx context.Context, member bool, projectID uuid.UUID, page, perPage int) ([]*model.Account, int, error)
	// DisableAccount disables a user in the database.
	DisableAccount(ctx context.Context, id uuid.UUID) error
	// EnableAccount enables a user in the database.
	EnableAccount(ctx context.Context, id uuid.UUID) error
	// VerifyAccount verifies a user in the database.
	VerifyAccount(ctx context.Context, id uuid.UUID) error
	// AccountExists checks if a user exists in the database.
	AccountExists(ctx context.Context, projectID uuid.UUID, username, email string) ([]*model.Account, error)
	// GetAccountCount retrieves the number of users in a project.
	GetAccountCount(ctx context.Context, projectID uuid.UUID) (uint32, error)
}

// SessionStore is the interface for interacting with the session database.
type SessionStore interface {
	// CreateSession creates a new session in the database.
	CreateSession(ctx context.Context, session *model.Session) error
	// ListSessions retrieves a list of sessions.
	ListSessions(ctx context.Context, orgID uuid.UUID, page, perPage int) ([]*model.Session, error)
	// ListActiveSessions retrieves a list of active sessions.
	ListActiveSessions(ctx context.Context, userID uuid.UUID) ([]*model.Session, error)
	// DeleteSession deletes a session from the database.
	DeleteSession(ctx context.Context, sessionID uuid.UUID) error
	// DeleteSessionByAccountID deletes a session from the database by user ID.
	DeleteSessionByAccountID(ctx context.Context, userID uuid.UUID) error
}

// ProjectMemberStore is the interface for interacting with the permission database.
type ProjectMemberStore interface {
	// CreateProjectMember creates a new permission in the database.
	CreateProjectMember(ctx context.Context, permission *model.ProjectMember) error
	// GetProjectMemberByID retrieves a permission by its ID.
	GetProjectMemberByID(ctx context.Context, orgID, userID uuid.UUID) (*model.ProjectMember, error)
	// ListProjectMembers retrieves a list of permissions.
	ListProjectMembers(ctx context.Context, projectID uuid.UUID, page, perPage int) ([]*model.ProjectMember, error)
	// ListProjectMembersByAccountIDs retrieves a list of permissions by project.
	ListProjectMembersByAccountIDs(ctx context.Context, projectID uuid.UUID, userIDs []uuid.UUID) ([]*model.ProjectMember, error)
	// UpdateProjectMember updates a permission in the database.
	UpdateProjectMember(ctx context.Context, permission *model.ProjectMember) error
	// DeleteProjectMember deletes a permission from the database.
	DeleteProjectMember(ctx context.Context, orgID, userID uuid.UUID) error
	// GetMemberCount retrieves the number of members in a project.
	GetMemberCount(ctx context.Context, projectID uuid.UUID) (uint32, error)
}

// ProviderStore is the interface for interacting with the provider database.
type ProviderStore interface {
	// CreateOauthProvider CreateProvider creates a new provider in the database.
	CreateOauthProvider(ctx context.Context, provider *model.OauthProvider) error
	// GetOauthProviderByID retrieves a provider by its ID.
	GetOauthProviderByID(ctx context.Context, id uuid.UUID) (*model.OauthProvider, error)
	// GetOauthProviderByName retrieves a provider by its name.
	GetOauthProviderByName(ctx context.Context, orgID uuid.UUID, provider string) (*model.OauthProvider, error)
	// ListOauthProviders retrieves a list of providers.
	ListOauthProviders(ctx context.Context, orgID uuid.UUID, page, perPage int) ([]*model.OauthProvider, uint32, error)
	// UpdateOauthProvider updates a provider in the database.
	UpdateOauthProvider(ctx context.Context, provider *model.OauthProvider) error
	// DeleteOauthProvider deletes a provider from the database.
	DeleteOauthProvider(ctx context.Context, id uuid.UUID) error
}

// RefreshTokenStore is the interface for interacting with the refresh token database.
type RefreshTokenStore interface {
	// CreateRefreshToken creates a new refresh token in the database.
	CreateRefreshToken(ctx context.Context, token *model.RefreshToken) error
	// GetRefreshTokenByID retrieves a
	GetRefreshTokenByID(ctx context.Context, token string) (*model.RefreshToken, error)
	// ListRefreshTokens retrieves a list of refresh tokens.
	ListRefreshTokens(ctx context.Context, page, perPage int) ([]*model.RefreshToken, error)
	// UpdateRefreshToken updates a refresh token in the database.
	UpdateRefreshToken(ctx context.Context, token *model.RefreshToken) error
	// DeleteRefreshToken deletes a refresh token from the database.
	DeleteRefreshToken(ctx context.Context, token string) error
}

// AccessKeyStore is the interface for interacting with the token database.
type AccessKeyStore interface {
	// CreateAccessKey creates a new token in the database.
	CreateAccessKey(ctx context.Context, token *model.AccessKey) error
	// GetAccessKeyByID retrieves a token by its ID.
	GetAccessKeyByID(ctx context.Context, id uuid.UUID) (*model.AccessKey, error)
	// ListAccountAccessKeys retrieves a list of tokens by user.
	ListAccountAccessKeys(ctx context.Context, orgID, userID uuid.UUID, page, perPage int) ([]*model.AccessKey, int, error)
	// DeleteAccessKey updates a token in the database.
	DeleteAccessKey(ctx context.Context, id uuid.UUID) error
}

// VerificationCodeStore is the interface for interacting with the verification code database.
type VerificationCodeStore interface {
	// CreateVerificationCode creates a new code
	CreateVerificationCode(ctx context.Context, code *model.VerificationCode) error
	// GetVerificationCode retrieves a code by hash id
	GetVerificationCode(ctx context.Context, code string) (*model.VerificationCode, error)
	// DeleteVerificationCode deletes the code
	DeleteVerificationCode(ctx context.Context, code string) error
}

// ClientStore is the interface for interacting with the client database.
type ClientStore interface {
	// CreateClient creates a new client in the database.
	CreateClient(ctx context.Context, client *model.Client) error
	// GetClientByID retrieves a client by its ID.
	GetClientByID(ctx context.Context, id uuid.UUID) (*model.Client, error)
	// ListClients retrieves a list of clients.
	ListClients(ctx context.Context, projectID uuid.UUID, page, perPage int) ([]*model.Client, int, error)
	// UpdateClient updates a client in the database.
	UpdateClient(ctx context.Context, client *model.Client) error
	// DeleteClient deletes a client from the database.
	DeleteClient(ctx context.Context, id uuid.UUID) error
}
