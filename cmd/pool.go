package cmd

import (
	"github.com/emrgen/authbase"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var poolCommand = &cobra.Command{
	Use:   "pool",
	Short: "pool commands",
}

func init() {
	poolCommand.AddCommand(createPoolCommand())
	poolCommand.AddCommand(listPoolCommand())
	poolCommand.AddCommand(updatePoolCommand())
	poolCommand.AddCommand(deletePoolCommand())
}

func createPoolCommand() *cobra.Command {
	var projectID string
	var name string

	command := &cobra.Command{
		Use:   "create",
		Short: "create pool",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("error creating client: %v", err)
				return
			}
			defer client.Close()

			if projectID == "" {
				account := getAccount(client)
				projectID = account.ProjectId
			}

			if name == "" {
				logrus.Error("missing required flags: --name")
				return
			}

			res, err := client.CreatePool(tokenContext(), &v1.CreatePoolRequest{
				ProjectId: projectID,
				Name:      name,
			})
			if err != nil {
				logrus.Errorf("error creating pool: %v", err)
				return
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Name", "Created At", "Updated At"})
			table.Append([]string{res.Pool.Id, res.Pool.Name, res.Pool.CreatedAt.AsTime().Format("2006-01-02 15:04:05"), res.Pool.UpdatedAt.AsTime().Format("2006-01-02 15:04:05")})
			table.Render()
		},
	}

	command.Flags().StringVarP(&projectID, "project", "r", "", "project id")
	command.Flags().StringVarP(&name, "name", "n", "", "name of the pool")

	return command
}

func listPoolCommand() *cobra.Command {
	var projectID string

	command := &cobra.Command{
		Use:   "list",
		Short: "list pools",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("error creating client: %v", err)
				return
			}
			if projectID == "" {
				account := getAccount(client)
				projectID = account.ProjectId
			}

			res, err := client.ListPools(tokenContext(), &v1.ListPoolsRequest{
				ProjectId: projectID,
			})
			if err != nil {
				logrus.Errorf("error listing pools: %v", err)
				return
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Name", "Created At", "Updated At"})
			for _, pool := range res.Pools {
				table.Append([]string{pool.Id, pool.Name, pool.CreatedAt.AsTime().Format("2006-01-02 15:04:05"), pool.UpdatedAt.AsTime().Format("2006-01-02 15:04:05")})
			}
			table.Render()
		},
	}

	command.Flags().StringVarP(&projectID, "project", "r", "", "project id")

	return command
}

func updatePoolCommand() *cobra.Command {
	var poolName string
	var poolID string

	command := &cobra.Command{
		Use:   "update",
		Short: "update pool",
		Run: func(cmd *cobra.Command, args []string) {
			if poolID == "" {
				logrus.Error("missing required flags: --pool-id")
				return
			}

			if poolName == "" {
				logrus.Error("missing required flags: --name")
				return
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("error creating client: %v", err)
				return
			}
			defer client.Close()

			res, err := client.UpdatePool(tokenContext(), &v1.UpdatePoolRequest{
				PoolId: poolID,
				Name:   poolName,
			})
			if err != nil {
				logrus.Errorf("error updating pool: %v", err)
				return
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Name", "Created At", "Updated At"})
			table.Append([]string{res.Pool.Id, res.Pool.Name, res.Pool.CreatedAt.AsTime().Format("2006-01-02 15:04:05"), res.Pool.UpdatedAt.AsTime().Format("2006-01-02 15:04:05")})
			table.Render()
		},
	}

	command.Flags().StringVarP(&poolID, "pool-id", "p", "", "id of the pool")
	command.Flags().StringVarP(&poolName, "name", "n", "", "name of the pool")

	return command

}

func deletePoolCommand() *cobra.Command {
	var poolID string

	command := &cobra.Command{
		Use:   "delete",
		Short: "delete pool",
		Run: func(cmd *cobra.Command, args []string) {
			if poolID == "" {
				logrus.Error("missing required flags: --pool-id")
				return
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("error creating client: %v", err)
				return
			}
			defer client.Close()

			_, err = client.DeletePool(tokenContext(), &v1.DeletePoolRequest{
				PoolId: poolID,
			})
			if err != nil {
				logrus.Errorf("error deleting pool: %v", err)
				return
			}

			logrus.Info("pool deleted")
		},
	}

	command.Flags().StringVarP(&poolID, "pool-id", "p", "", "id of the pool")

	return command
}
