package cmd

import (
	"github.com/emrgen/authbase"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var configCommand = &cobra.Command{
	Use:   "config",
	Short: "config commands",
}

func init() {
	codeCommand := createConfigCommand()
	codeCommand.AddCommand(setCodeMediumCommand())
	configCommand.AddCommand(codeCommand)
}

func createConfigCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "code",
		Short: "configure code medium",
	}

	return command
}

func setCodeMediumCommand() *cobra.Command {
	var medium string
	var via string

	command := &cobra.Command{
		Use:   "set",
		Short: "set code medium",
		Run: func(cmd *cobra.Command, args []string) {
			verifyContext()

			if medium == "" {
				logrus.Error("missing required flags: --medium")
				return
			}

			if via == "" {
				logrus.Error("missing required flags: --via")
				return
			}

			_, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("error creating client: %v", err)
				return
			}

			//client.UpdateOrganizationConfig(context.Background(), &v1.UpdateOrganizationConfigRequest{})

			logrus.Infof("code medium set successfully, %v --> %v", medium, via)
		},
	}

	bindContextFlags(command)
	command.Flags().StringVarP(&medium, "medium", "m", "", "medium")
	command.Flags().StringVarP(&via, "via", "v", "", "via")

	return command
}
