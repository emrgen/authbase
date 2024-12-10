package config

import "os"

// config package is used to load the configuration from the environment variables
// the logic behind having separate config package is usability
// we can use the same package in other services as well to load the configuration

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
	if dbConfig.Type == "" || dbConfig.Type == "sqlite" {
		dbConfig.Type = "sqlite"
		dbConfig.FilePath = "./.tmp/db/authbase.db"
	} else {
		dbConfig.ConnectionString = os.Getenv("DB_CONNECTION_STRING")
	}

	appKey := os.Getenv("APP_KEY")

	config := &Config{
		Environment: Environment(env),
		DB:          dbConfig,
		AppKey:      appKey,
	}

	return config, nil
}
