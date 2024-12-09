package cmd

import "github.com/spf13/cobra"

var orgCommand = &cobra.Command{
	Use:   "org",
	Short: "org commands",
}

func init() {
	orgCommand.AddCommand(createOrgCommand())
	orgCommand.AddCommand(listOrgCommand())
	orgCommand.AddCommand(updateOrgCommand())
	orgCommand.AddCommand(deleteOrgCommand())
}

func createOrgCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "create",
		Short: "create org",
	}

	return command
}

func listOrgCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "list",
		Short: "list org",
	}

	return command
}

func updateOrgCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "update",
		Short: "update org",
	}

	return command
}

func deleteOrgCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "delete",
		Short: "delete org",
	}

	return command
}
