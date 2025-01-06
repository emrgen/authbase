package config

type OrgConfig struct {
	DatabaseConn              string
	AccessTokenExpireInterval int
}

type AdminOrgConfig struct {
	OrgName     string
	VisibleName string
	Email       string
	Password    string
}

func (a AdminOrgConfig) Valid() bool {
	if a.OrgName == "" || a.VisibleName == "" || a.Email == "" || a.Password == "" {
		return false
	}
	return true
}
