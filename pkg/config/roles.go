package config

var (
	// project Roles
	projectCreateRole = "project:create"
	projectReadRole   = "project:read"
	projectUpdateRole = "project:update"
	projectDeleteRole = "project:delete"
	// pool Roles
	poolCreateRole       = "pool:create"
	poolReadRole         = "pool:read"
	poolUpdateRole       = "pool:update"
	poolDeleteRole       = "pool:delete"
	poolMemberAddRole    = "pool:member:add"
	poolMemberRemoveRole = "pool:member:remove"
	// client Roles
	clientCreateRole = "client:create"
	clientReadRole   = "client:read"
	clientUpdateRole = "client:update"
	clientDeleteRole = "client:delete"
	// user Roles
	userCreateRole      = "user:create"
	userReadRole        = "user:read"
	userWriteRole       = "user:write"
	userDeleteRole      = "user:delete"
	userGroupAddRole    = "user:group:add"
	userGroupRemoveRole = "user:group:remove"
	// group Roles
	groupCreateRole     = "group:create"
	groupReadRole       = "group:read"
	groupUpdateRole     = "group:update"
	groupDeleteRole     = "group:delete"
	groupRoleAddRole    = "group:role:add"
	groupRoleRemoveRole = "group:role:remove"
	// role Roles
	roleCreateRole = "role:create"
	roleReadRole   = "role:read"
	roleUpdateRole = "role:update"
	roleDeleteRole = "role:delete"
	// internal
	internalRole = "internal" // users with this role can access all auth server resources
)

func ProjectOwnerRoles() Roles {
	return mergeRoles(
		ProjectAdminRoles(),
		newRoleList(
			projectDeleteRole,
		),
	)
}

func ProjectAdminRoles() Roles {
	return mergeRoles(
		ProjectViewerRoles(),
		newRoleList(
			projectCreateRole,
			projectReadRole,
			projectUpdateRole,
		),
	)
}

func ProjectViewerRoles() Roles {
	return newRoleList(
		projectReadRole,
	)
}

func PoolOwnerRoles() Roles {
	return mergeRoles(
		PoolAdminRoles(),
		newRoleList(
			poolDeleteRole,
		),
	)
}

func PoolAdminRoles() Roles {
	return mergeRoles(
		PoolViewerRoles(),
		newRoleList(
			poolReadRole,
			poolUpdateRole,
			poolDeleteRole,
			poolCreateRole,
			poolMemberAddRole,
			poolMemberRemoveRole,
		),
	)
}

func PoolViewerRoles() Roles {
	return newRoleList(poolReadRole, internalRole)
}

func ClientAdminRoles() Roles {
	return mergeRoles(
		ClientViewerRoles(),
		newRoleList(
			clientCreateRole,
			clientUpdateRole,
			clientDeleteRole,
		),
	)
}

func ClientViewerRoles() Roles {
	return newRoleList(clientReadRole, internalRole)
}

func UserAdminRoles() Roles {
	return mergeRoles(
		UserViewerRoles(),
		newRoleList(
			userCreateRole,
			userWriteRole,
			userDeleteRole,
			userGroupAddRole,
			userGroupRemoveRole,
		),
	)
}

func UserViewerRoles() Roles {
	return newRoleList(userReadRole)
}

func GroupAdminRoles() Roles {
	return mergeRoles(
		GroupViewerRoles(),
		newRoleList(
			groupCreateRole,
			groupUpdateRole,
			groupDeleteRole,
			groupRoleAddRole,
			groupRoleRemoveRole,
		),
	)
}

func GroupViewerRoles() Roles {
	return newRoleList(groupReadRole)
}

func RoleAdminRoles() Roles {
	return mergeRoles(
		RoleViewerRoles(),
		newRoleList(
			roleCreateRole,
			roleUpdateRole,
			roleDeleteRole,
		),
	)
}

func RoleViewerRoles() Roles {
	return newRoleList(roleReadRole)
}

type Roles interface {
	roles() []role
	Roles() []string
}

type role struct {
	name string
}

func newRole(name string) *role {
	return &role{
		name: name,
	}
}

func (r *role) roles() []role {
	return []role{*r}
}

func (r *role) Roles() []string {
	return []string{r.name}
}

type roleList []role

func newRoleList(roles ...string) roleList {
	var list roleList
	for _, r := range roles {
		list = append(list, role{name: r})
	}
	return list
}

func (r roleList) roles() []role {
	return r
}

func (r roleList) Roles() []string {
	var roles []string
	for _, role := range r {
		roles = append(roles, role.name)
	}
	return roles
}

func mergeRoles(roles ...Roles) Roles {
	var list roleList
	for _, r := range roles {
		list = append(list, r.roles()...)
	}
	return list
}
