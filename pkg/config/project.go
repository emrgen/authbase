package config

type ProjectConfig struct {
	DatabaseConn              string
	AccessTokenExpireInterval int
}

type AdminProjectConfig struct {
	OrgName      string
	VisibleName  string
	Email        string
	Password     string
	ClientId     string
	ClientSecret string
}

func (a AdminProjectConfig) Valid() bool {
	if a.OrgName == "" || a.VisibleName == "" || a.Email == "" || a.Password == "" {
		return false
	}
	return true
}
