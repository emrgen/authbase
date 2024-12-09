package cmd

import "github.com/spf13/cobra"

var orgCommand = &cobra.Command{
	Use:   "org",
	Short: "org commands",
}

func init() {
	orgCommand.AddCommand(createOrgCommand())
}

func createOrgCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "create",
		Short: "create org",
	}

	return command
}
