package cmd

import "github.com/spf13/cobra"

var tokenCommand = &cobra.Command{
	Use:   "token",
	Short: "token commands",
}

func init() {
	tokenCommand.AddCommand(createTokenCommand())
	tokenCommand.AddCommand(listTokenCommand())
	tokenCommand.AddCommand(deleteTokenCommand())
}

func createTokenCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "create",
		Short: "create token",
	}

	return command
}

func listTokenCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "list",
		Short: "list token",
	}

	return command
}

func deleteTokenCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "delete",
		Short: "delete token",
	}

	return command
}
