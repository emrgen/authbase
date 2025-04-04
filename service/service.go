package service

// Package service provides a service registry for gRPC and HTTP services

// Config holds the configuration for the service registry
type Config struct {
	grpc string
	http string
}

// Register registers the service with the given configuration
func Register(cfg *Config) error {
	return nil
}
