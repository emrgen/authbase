package cmd

import "github.com/spf13/cobra"

var userCommand = &cobra.Command{
	Use:   "org",
	Short: "org commands",
}

func init() {
	userCommand.AddCommand(createUserCommand())
}

func createUserCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "create",
		Short: "create user",
	}

	return command
}
