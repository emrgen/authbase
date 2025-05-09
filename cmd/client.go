package cmd

import (
	"github.com/emrgen/authbase"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/olekukonko/tablewriter"
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
			client, err := authbase.NewClient("4000")
			if err != nil {
				logrus.Error(err)
				return
			}
			defer client.Close()

			if poolID == "" {
				poolID = getAccount(client).PoolId
			}

			if name == "" {
				logrus.Error("missing required flags: --name")
				return
			}

			res, err := client.CreateClient(tokenContext(), &v1.CreateClientRequest{
				PoolId: poolID,
				Name:   name,
			})
			if err != nil {
				logrus.Error(err)
				return
			}

			table := tablewriter.NewWriter(cmd.OutOrStdout())
			table.SetHeader([]string{"Pool ID", "Client ID", "Name", "Created By"})
			table.Append([]string{res.Client.PoolId, res.Client.Id, res.Client.Name, res.Client.CreatedByUser.VisibleName})
			table.Render()

		},
	}

	command.Flags().StringVarP(&poolID, "pool-id", "p", "", "Pool ID")
	command.Flags().StringVarP(&name, "name", "n", "", "Name")

	return command
}

func clientGetCmd() *cobra.Command {
	var clientID string
	command := &cobra.Command{
		Use:   "get",
		Short: "Get a client",
		Run: func(cmd *cobra.Command, args []string) {
			if clientID == "" {
				logrus.Error("missing required flags: --client-id")
				return
			}

			client, err := authbase.NewClient("4000")
			if err != nil {
				logrus.Error(err)
				return
			}
			defer client.Close()

			res, err := client.GetClient(tokenContext(), &v1.GetClientRequest{
				ClientId: clientID,
			})
			if err != nil {
				logrus.Error(err)
				return
			}

			table := tablewriter.NewWriter(cmd.OutOrStdout())
			table.SetHeader([]string{"ID", "Pool ID", "Name", "Created By"})
			table.Append([]string{res.Client.Id, res.Client.PoolId, res.Client.Name, res.Client.CreatedByUser.VisibleName})
			table.Render()
		},
	}

	command.Flags().StringVarP(&clientID, "client-id", "c", "", "Client ID")

	return command
}

func clientListCmd() *cobra.Command {
	var poolID string

	command := &cobra.Command{
		Use:   "list",
		Short: "List clients",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := authbase.NewClient("4000")
			if err != nil {
				logrus.Error(err)
				return
			}
			defer client.Close()

			if poolID == "" {
				account := getAccount(client)
				poolID = account.PoolId
			}

			res, err := client.ListClients(tokenContext(), &v1.ListClientsRequest{
				PoolId: poolID,
			})
			if err != nil {
				logrus.Error(err)
				return
			}

			table := tablewriter.NewWriter(cmd.OutOrStdout())
			table.SetHeader([]string{"ID", "Name"})
			for _, client := range res.Clients {
				table.Append([]string{client.Id, client.Name})
			}
			table.Render()
		},
	}

	command.Flags().StringVarP(&poolID, "pool-id", "p", "", "Pool ID")
	return command
}

func clientDeleteCmd() *cobra.Command {
	var clientID string

	command := &cobra.Command{
		Use:   "delete",
		Short: "Delete a client",
		Run: func(cmd *cobra.Command, args []string) {
			if clientID == "" {
				logrus.Error("missing required flags: --client-id")
				return
			}

			client, err := authbase.NewClient("4000")
			if err != nil {
				logrus.Error(err)
				return
			}
			defer client.Close()

			_, err = client.DeleteClient(tokenContext(), &v1.DeleteClientRequest{
				ClientId: clientID,
			})
			if err != nil {
				logrus.Error(err)
				return
			}

			logrus.Infof("Client %s deleted", clientID)
		},
	}

	command.Flags().StringVarP(&clientID, "client-id", "c", "", "Client ID")

	return command
}
