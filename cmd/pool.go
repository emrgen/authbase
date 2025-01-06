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
	poolCommand.AddCommand(deletePoolCommand())
}

func createPoolCommand() *cobra.Command {
	var projectID string
	var name string

	command := &cobra.Command{
		Use:   "create",
		Short: "create pool",
		Run: func(cmd *cobra.Command, args []string) {
			if projectID == "" {
				logrus.Error("missing required flags: --project")
				return
			}

			if name == "" {
				logrus.Error("missing required flags: --name")
				return
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("error creating client: %v", err)
				return
			}
			defer client.Close()

			res, err := client.CreatePool(tokenContext(), &v1.CreatePoolRequest{})
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
			if projectID == "" {
				logrus.Error("missing required flags: --project")
				return
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("error creating client: %v", err)
				return
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

func deletePoolCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "delete",
		Short: "delete pool",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	return command
}
