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

func init() {
	rootCmd.AddCommand(dbCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(projectCommand)
	rootCmd.AddCommand(userCommand)
	rootCmd.AddCommand(memberCommand)
	rootCmd.AddCommand(configCommand)
	rootCmd.AddCommand(accessKeyCommand)
	rootCmd.AddCommand(strategyCommand)
	rootCmd.AddCommand(contextCommand)
	rootCmd.AddCommand(keygenCommand)
	rootCmd.AddCommand(poolCommand)
	rootCmd.AddCommand(clientCommand)
	rootCmd.AddCommand(groupCommand)
	rootCmd.AddCommand(roleCommand)
	rootCmd.AddCommand(tokenCommand)
	rootCmd.AddCommand(idpCommand)

	ctx := readContext()
	if ctx.Token != "" {
		rootCmd.PersistentFlags().StringVarP(&Token, "token", "t", ctx.Token, "token")
	}
}
