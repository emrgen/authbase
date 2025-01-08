package config

import (
	"os"
)

// config package is used to load the configuration from the environment variables
// the logic behind having separate config package is usability
// we can use the same package in other services as well to load the configuration

type AppMode string

const (
	ModeMultiStore  AppMode = "multistore"
	ModeSingleStore AppMode = "singlestore"
)

type Environment string

const (
	Development Environment = "development"
	Production  Environment = "production"
	Staging     Environment = "staging"
	Testing     Environment = "testing"
)

type Config struct {
	Environment Environment
	DB          *DBConfig
	AdminOrg    *AdminProjectConfig
	Mode        AppMode
	AppKey      string
}

type DBConfig struct {
	Type             string
	ConnectionString string
	FilePath         string
}

// FromEnv loads the configuration from the environment variables
func FromEnv() (*Config, error) {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	dbConfig := &DBConfig{}
	dbConfig.Type = os.Getenv("DB_TYPE")
	if dbConfig.Type == "" || dbConfig.Type == "sqlite3" {
		dbConfig.Type = "sqlite3"
		dbConfig.FilePath = "./.tmp/db/authbase.db"
	} else {
		dbConfig.ConnectionString = os.Getenv("DB_CONNECTION_STRING")
	}

	appKey := os.Getenv("APP_KEY")

	adminOrgConfig := &AdminProjectConfig{}
	adminOrgConfig.OrgName = os.Getenv("ADMIN_ORGANIZATION_NAME")
	adminOrgConfig.VisibleName = os.Getenv("SUPER_ADMIN_VISIBLE_NAME")
	adminOrgConfig.Email = os.Getenv("SUPER_ADMIN_EMAIL")
	adminOrgConfig.Password = os.Getenv("SUPER_ADMIN_PASSWORD")
	adminOrgConfig.ClientId = os.Getenv("SUPER_ADMIN_CLIENT_ID")
	adminOrgConfig.ClientSecret = os.Getenv("SUPER_ADMIN_CLIENT_SECRET")

	mode := os.Getenv("APP_MODE")
	if mode == "" {
		mode = "singlestore"
	}

	config := &Config{
		Environment: Environment(env),
		DB:          dbConfig,
		AppKey:      appKey,
		AdminOrg:    adminOrgConfig,
		Mode:        AppMode(mode),
	}

	return config, nil
}

var config *Config

// GetConfig returns the configuration
func GetConfig() *Config {
	return config
}

func init() {
	cfg, err := FromEnv()
	if err != nil {
		panic(err)
	}

	config = cfg
}
