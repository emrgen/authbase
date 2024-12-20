package cmd

import "github.com/spf13/cobra"

var strategyCommand = &cobra.Command{
	Use:   "strategy",
	Short: "user commands",
}

func init() {
	strategyCommand.AddCommand(createStrategyCommand())
}

func createStrategyCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "create",
		Short: "create strategy",
	}

	return command
}
