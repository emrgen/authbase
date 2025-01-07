package server

type roles interface {
	roles() []role
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

func mergeRoles(roles ...roles) roles {
	var list roleList
	for _, r := range roles {
		list = append(list, r.roles()...)
	}
	return list
}

var (
	// project roles
	projectCreateRole = "project:create"
	projectReadRole   = "project:read"
	projectUpdateRole = "project:update"
	projectDeleteRole = "project:delete"
	// pool roles
	poolCreateRole       = "pool:create"
	poolReadRole         = "pool:read"
	poolUpdateRole       = "pool:update"
	poolDeleteRole       = "pool:delete"
	poolMemberAddRole    = "pool:member:add"
	poolMemberRemoveRole = "pool:member:remove"
	// client roles
	clientCreateRole = "client:create"
	clientReadRole   = "client:read"
	clientUpdateRole = "client:update"
	clientDeleteRole = "client:delete"
	// user roles
	userCreateRole      = "user:create"
	userReadRole        = "user:read"
	userWriteRole       = "user:write"
	userDeleteRole      = "user:delete"
	userGroupAddRole    = "user:group:add"
	userGroupRemoveRole = "user:group:remove"
	// group roles
	groupCreateRole     = "group:create"
	groupReadRole       = "group:read"
	groupUpdateRole     = "group:update"
	groupDeleteRole     = "group:delete"
	groupRoleAddRole    = "group:role:add"
	groupRoleRemoveRole = "group:role:remove"
	// role roles
	roleCreateRole = "role:create"
	roleReadRole   = "role:read"
	roleUpdateRole = "role:update"
	roleDeleteRole = "role:delete"
)

func projectOwnerRoles() roles {
	return mergeRoles(
		projectAdminRoles(),
		newRoleList(
			projectDeleteRole,
		),
	)
}

func projectAdminRoles() roles {
	return mergeRoles(
		projectViewerRoles(),
		newRoleList(
			projectCreateRole,
			projectReadRole,
			projectUpdateRole,
		),
	)
}

func projectViewerRoles() roles {
	return newRoleList(
		projectReadRole,
	)
}

func poolOwnerRoles() roles {
	return mergeRoles(
		poolAdminRoles(),
		newRoleList(
			poolDeleteRole,
		),
	)
}

func poolAdminRoles() roles {
	return mergeRoles(
		poolViewerRoles(),
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

func poolViewerRoles() roles {
	return mergeRoles(
		newRole(poolReadRole),
	)
}

func clientAdminRoles() roles {
	return mergeRoles(
		clientViewerRoles(),
		newRoleList(
			clientCreateRole,
			clientUpdateRole,
			clientDeleteRole,
		),
	)
}

func clientViewerRoles() roles {
	return newRole(clientReadRole)
}

func userAdminRoles() roles {
	return mergeRoles(
		userViewerRoles(),
		newRoleList(
			userCreateRole,
			userWriteRole,
			userDeleteRole,
			userGroupAddRole,
			userGroupRemoveRole,
		),
	)
}

func userViewerRoles() roles {
	return newRole(userReadRole)
}

func groupAdminRoles() roles {
	return mergeRoles(
		groupViewerRoles(),
		newRoleList(
			groupCreateRole,
			groupUpdateRole,
			groupDeleteRole,
			groupRoleAddRole,
			groupRoleRemoveRole,
		),
	)
}

func groupViewerRoles() roles {
	return newRole(groupReadRole)
}

func roleAdminRoles() roles {
	return mergeRoles(
		roleViewerRoles(),
		newRoleList(
			roleCreateRole,
			roleUpdateRole,
			roleDeleteRole,
		),
	)
}

func roleViewerRoles() roles {
	return newRole(roleReadRole)
}
