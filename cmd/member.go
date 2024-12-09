package cmd

import "github.com/spf13/cobra"

var memberCommand = &cobra.Command{
	Use:   "org",
	Short: "org commands",
}

func init() {
	memberCommand.AddCommand(createMemberCommand())
}

func createMemberCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "create",
		Short: "create user",
	}

	return command
}
