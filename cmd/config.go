package cmd

import "github.com/spf13/cobra"

var configCommand = &cobra.Command{
	Use:   "config",
	Short: "config commands",
}

func init() {
	configCommand.AddCommand(createConfigCommand())
}

func createConfigCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "code",
		Short: "configure code medium",
	}

	return command
}
