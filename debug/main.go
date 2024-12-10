package main

import (
	"github.com/emrgen/authbase/pkg/server"
	"os"
)

// main starts the gRPC and HTTP servers.
// this is intended for local debugging only.
func main() {
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "4000"
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "4001"
	}

	svr := server.NewServerFromEnv()

	err := svr.Start(grpcPort, httpPort)
	if err != nil {
		return
	}
}
