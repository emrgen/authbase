package config

type ProjectConfig struct {
	DatabaseConn              string
	AccessTokenExpireInterval int
}

// AdminProjectConfig is used to create a initial master project.
// The master project is used to manage all other projects.
type AdminProjectConfig struct {
	OrgName      string
	VisibleName  string
	Email        string
	Password     string
	ClientId     string
	ClientSecret string
}

// Valid checks if the AdminProjectConfig is valid.
func (a AdminProjectConfig) Valid() bool {
	if a.OrgName == "" || a.VisibleName == "" || a.Email == "" || a.Password == "" {
		return false
	}
	return true
}
