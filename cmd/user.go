package cmd

import "github.com/spf13/cobra"

var userCommand = &cobra.Command{
	Use:   "user",
	Short: "user commands",
}

func init() {
	userCommand.AddCommand(createUserCommand())
	userCommand.AddCommand(listUserCommand())
	userCommand.AddCommand(updateUserCommand())
	userCommand.AddCommand(deleteUserCommand())
}

func createUserCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "create",
		Short: "create user",
	}

	return command
}

func listUserCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "list",
		Short: "list user",
	}

	return command
}

func updateUserCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "update",
		Short: "update user",
	}

	return command
}

func deleteUserCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "delete",
		Short: "delete user",
	}

	return command
}
