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
	ErrRoleNotFound            = errors.New("role not found")
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
	PoolStore
	PoolMemberStore
	GroupStore
	RoleStore
	ApplicationStore
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
	GetAccountByEmail(ctx context.Context, poolID uuid.UUID, email string) (*model.Account, error)
	// GetAccountByID retrieves a user by their ID.
	GetAccountByID(ctx context.Context, id uuid.UUID) (*model.Account, error)
	// UpdateAccount updates a user in the database.
	UpdateAccount(ctx context.Context, user *model.Account) error
	// DeleteAccount deletes a user from the database.
	DeleteAccount(ctx context.Context, id uuid.UUID) error
	// ListProjectAccounts retrieves a list of users by project.
	ListProjectAccounts(ctx context.Context, member bool, projectID uuid.UUID, page, perPage int) ([]*model.Account, int, error)
	// ListPoolAccounts retrieves a list of users by pool.
	ListPoolAccounts(ctx context.Context, member bool, poolID uuid.UUID, page, perPage int) ([]*model.Account, int, error)
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
	// ListActiveAccounts retrieves a list of sessions.
	ListActiveAccounts(ctx context.Context, poolID uuid.UUID, page, perPage int) ([]*model.Session, error)
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

type PoolStore interface {
	// CreatePool creates a new pool in the database.
	CreatePool(ctx context.Context, pool *model.Pool) error
	//GetMasterPool retrieves the master pool.
	GetMasterPool(ctx context.Context, projectID uuid.UUID) (*model.Pool, error)
	// GetPoolByID retrieves a pool by its name.
	GetPoolByID(ctx context.Context, id uuid.UUID) (*model.Pool, error)
	// ListPools retrieves a list of pools.
	ListPools(ctx context.Context, projectID uuid.UUID, page, perPage int) ([]*model.Pool, int, error)
	// UpdatePool updates a pool in the database.
	UpdatePool(ctx context.Context, pool *model.Pool) error
	// DeletePool deletes a pool from the database.
	DeletePool(ctx context.Context, id uuid.UUID) error
}

type PoolMemberStore interface {
	// AddPoolMember creates a new pool member in the database.
	AddPoolMember(ctx context.Context, member *model.PoolMember) error
	// GetPoolMember retrieves a pool member by its ID.
	GetPoolMember(ctx context.Context, poolID, accountID uuid.UUID) (*model.PoolMember, error)
	// ListPoolMembers retrieves a list of pool members.
	ListPoolMembers(ctx context.Context, poolID uuid.UUID, page, perPage int) ([]*model.PoolMember, int, error)
	// UpdatePoolMember updates a pool member in the database.
	UpdatePoolMember(ctx context.Context, member *model.PoolMember) error
	// RemovePoolMember deletes a pool member from the database.
	RemovePoolMember(ctx context.Context, poolID, accountID uuid.UUID) error
}

type GroupStore interface {
	// CreateGroup creates a new group in the database.
	CreateGroup(ctx context.Context, group *model.Group) error
	// GetGroup retrieves a group by its ID.
	GetGroup(ctx context.Context, id uuid.UUID) (*model.Group, error)
	// ListGroups retrieves a list of groups.
	ListGroups(ctx context.Context, projectID uuid.UUID, page, perPage int) ([]*model.Group, int, error)
	// UpdateGroup updates a group in the database.
	UpdateGroup(ctx context.Context, group *model.Group) error
	// DeleteGroup deletes a group from the database.
	DeleteGroup(ctx context.Context, id uuid.UUID) error
	// AddGroupMember creates a new group member in the database.
	AddGroupMember(ctx context.Context, member *model.GroupMember) error
	// ListGroupMemberByAccount retrieves a group by its account ID.
	ListGroupMemberByAccount(ctx context.Context, accountID uuid.UUID) ([]*model.GroupMember, error)
	// RemoveGroupMember deletes a group member from the database.
	RemoveGroupMember(ctx context.Context, groupID, accountID uuid.UUID) error
	// ListGroupMembers retrieves a list of group members.
	ListGroupMembers(ctx context.Context, groupID uuid.UUID, page, perPage int) ([]*model.GroupMember, int, error)
}

type RoleStore interface {
	// CreateRole creates a new role in the database.
	CreateRole(ctx context.Context, role *model.Role) error
	// GetRole retrieves a role by its ID.
	GetRole(ctx context.Context, poolID uuid.UUID, name string) (*model.Role, error)
	// ListRolesByNames retrieves a list of roles by names.
	ListRolesByNames(ctx context.Context, poolID uuid.UUID, names []string) ([]*model.Role, error)
	// ListRoles retrieves a list of roles.
	ListRoles(ctx context.Context, poolID uuid.UUID, page, perPage int) ([]*model.Role, int, error)
	// UpdateRole updates a role in the database.
	UpdateRole(ctx context.Context, role *model.Role) error
	// DeleteRole deletes a role from the database.
	DeleteRole(ctx context.Context, poolID uuid.UUID, name string) error
}

type ApplicationStore interface {
	// CreateApplication creates a new application in the database.
	CreateApplication(ctx context.Context, app *model.Application) error
	// GetApplication retrieves an application by its ID.
	GetApplication(ctx context.Context, id uuid.UUID) (*model.Application, error)
	// ListApplications retrieves a list of applications.
	ListApplications(ctx context.Context, projectID uuid.UUID, page, perPage int) ([]*model.Application, int, error)
	// UpdateApplication updates an application in the database.
	UpdateApplication(ctx context.Context, app *model.Application) error
	// DeleteApplication deletes an application from the database.
	DeleteApplication(ctx context.Context, id uuid.UUID) error
}
