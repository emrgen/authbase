package cmd

import (
	"github.com/sirupsen/logrus"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "authbase",
	Short: "authbase is a service for authentication",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}

var Token string
var Organization string

func init() {
	rootCmd.AddCommand(dbCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(orgCommand)
	rootCmd.AddCommand(userCommand)
	rootCmd.AddCommand(memberCommand)
	rootCmd.AddCommand(configCommand)
	rootCmd.AddCommand(tokenCommand)
	rootCmd.AddCommand(strategyCommand)
	rootCmd.AddCommand(contextCommand)

	ctx := readContext()
	if ctx.Organization != "" {
		rootCmd.PersistentFlags().StringVarP(&Organization, "organization", "o", ctx.Organization, "organization name")
	}

	if ctx.Token != "" {
		logrus.Info("token: ", ctx.Token)
		rootCmd.PersistentFlags().StringVarP(&Token, "token", "t", ctx.Token, "token")
	}
}
