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
}

type DBConfig struct {
	Type             string
	ConnectionString string
}

// FromEnv loads the configuration from the environment variables
func FromEnv() (*Config, error) {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = "sqlite3"
	}

	dbConnectionString := os.Getenv("DB_CONNECTION_STRING")
	if dbConnectionString == "" {
		dbConnectionString = "authbase.db"
	}

	dbConfig := &DBConfig{
		Type:             dbType,
		ConnectionString: dbConnectionString,
	}

	config := &Config{
		Environment: Environment(env),
		DB:          dbConfig,
	}

	return config, nil
}
