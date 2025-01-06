package cmd

import (
	"github.com/emrgen/authbase"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var clientCommand = &cobra.Command{
	Use:   "client",
	Short: "Client commands",
}

func init() {
	clientCommand.AddCommand(clientCreateCmd())
	clientCommand.AddCommand(clientGetCmd())
	clientCommand.AddCommand(clientListCmd())
	clientCommand.AddCommand(clientDeleteCmd())
}

func clientCreateCmd() *cobra.Command {
	var poolID string
	var name string

	command := &cobra.Command{
		Use:   "create",
		Short: "Create a client",
		Run: func(cmd *cobra.Command, args []string) {
			if poolID == "" {
				logrus.Error("missing required flags: --pool-id")
				return
			}

			if name == "" {
				logrus.Error("missing required flags: --name")
				return
			}

			client, err := authbase.NewClient("4000")
			if err != nil {
				logrus.Error(err)
				return
			}
			defer client.Close()

			client.CreateClient(tokenContext(), &v1.CreateClientRequest{
				PoolId: poolID,
				Name:   name,
			})

		},
	}

	command.Flags().StringVarP(&poolID, "pool-id", "p", "", "Pool ID")
	command.Flags().StringVarP(&name, "name", "n", "", "Name")

	return command
}

func clientGetCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "get",
		Short: "Get a client",
		Run: func(cmd *cobra.Command, args []string) {
			//TODO implement me
		},
	}

	return command
}

func clientListCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "list",
		Short: "List clients",
		Run: func(cmd *cobra.Command, args []string) {
			//TODO implement me
		},
	}

	return command
}

func clientDeleteCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "delete",
		Short: "Delete a client",
		Run: func(cmd *cobra.Command, args []string) {
			//TODO implement me
		},
	}

	return command
}
