package cmd

import "github.com/spf13/cobra"

var memberCommand = &cobra.Command{
	Use:   "member",
	Short: "member commands",
}

func init() {
	memberCommand.AddCommand(createMemberCommand())
	memberCommand.AddCommand(listMemberCommand())
	memberCommand.AddCommand(updateMemberCommand())
	memberCommand.AddCommand(deleteMemberCommand())
}

func createMemberCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "create",
		Short: "create user",
	}

	return command
}

func listMemberCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "list",
		Short: "list user",
	}

	return command
}

func updateMemberCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "update",
		Short: "update user",
	}

	return command
}

func deleteMemberCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "delete",
		Short: "delete user",
	}

	return command
}
