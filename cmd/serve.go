package cmd

import (
	"github.com/emrgen/authbase/pkg/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start the server",
	Long:  `starts the http and grpc server for firstime service`,
	Run: func(cmd *cobra.Command, args []string) {
		//get grpc port
		grpcPort, err := cmd.Flags().GetString("gp")
		if err != nil {
			logrus.Errorf("error getting grpc port: %v", err)
		}

		//get http port
		httpPort, err := cmd.Flags().GetString("hp")
		if err != nil {
			logrus.Errorf("error getting http port: %v", err)
		}

		logrus.Infof("grpc port: %s, http port: %s", grpcPort, httpPort)
		svr := server.NewServerFromEnv()
		err = svr.Start(grpcPort, httpPort)
		if err != nil {
			logrus.Errorf("error starting service: %v", err)
			return
		}
	},
}

// init function to add flags to serveCmd
func init() {
	var grpcPort, httpPort string
	serveCmd.Flags().StringVar(&grpcPort, "gp", "4000", "Port to run grpc server on")
	serveCmd.Flags().StringVar(&httpPort, "hp", "4001", "Port to run http server on")
}
