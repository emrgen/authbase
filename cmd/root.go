package cmd

import (
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
var ProjectId string
var Username string
var Password string

func init() {
	rootCmd.AddCommand(dbCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(projectCommand)
	rootCmd.AddCommand(userCommand)
	rootCmd.AddCommand(memberCommand)
	rootCmd.AddCommand(configCommand)
	rootCmd.AddCommand(tokenCommand)
	rootCmd.AddCommand(strategyCommand)
	rootCmd.AddCommand(contextCommand)
	rootCmd.AddCommand(keygenCommand)

	ctx := readContext()
	if ctx.ProjectId != "" {
		rootCmd.PersistentFlags().StringVarP(&ProjectId, "project-id", "p", ctx.ProjectId, "project id")
	}

	if ctx.Token != "" {
		rootCmd.PersistentFlags().StringVarP(&Token, "token", "t", ctx.Token, "token")
	}

	if ctx.Username != "" {
		rootCmd.PersistentFlags().StringVarP(&Username, "username", "u", ctx.Username, "username")
	}

	if ctx.Password != "" {
		rootCmd.PersistentFlags().StringVarP(&Password, "password", "p", ctx.Password, "password")
	}
}
