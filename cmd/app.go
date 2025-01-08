package cmd

import (
	"github.com/emrgen/authbase"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var applicationCommand = &cobra.Command{
	Use:   "app",
	Short: "Application commands",
}

func init() {
	applicationCommand.AddCommand(applicationCreateCommand())
	applicationCommand.AddCommand(applicationListCommand())
	applicationCommand.AddCommand(applicationDeleteCommand())
	applicationCommand.AddCommand(applicationUpdateCommand())
}

func applicationCreateCommand() *cobra.Command {
	var poolID string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an application",
		Run: func(cmd *cobra.Command, args []string) {
			if poolID == "" {
				logrus.Error("missing required flags: --pool-id")
				return
			}

			client, err := authbase.NewClient("4000")
			if err != nil {
				logrus.Error(err)
				return
			}
			defer client.Close()

			res, err := client.CreateApplication(tokenContext(), &v1.CreateApplicationRequest{
				PoolId: "",
				Name:   "",
			})
			if err != nil {
				logrus.Error(err)
				return
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Name", "Pool ID"})
			table.Append([]string{res.Application.Id, res.Application.Name, res.Application.PoolId})
			table.Render()
		},
	}

	cmd.Flags().StringVar(&poolID, "pool-id", "", "The pool id")

	return cmd
}

func applicationListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List applications",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := authbase.NewClient("4000")
			if err != nil {
				logrus.Error(err)
				return
			}
			defer client.Close()

			res, err := client.ListApplications(tokenContext(), &v1.ListApplicationsRequest{})
			if err != nil {
				logrus.Error(err)
				return
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Name", "Pool ID"})
			for _, app := range res.Applications {
				table.Append([]string{app.Id, app.Name, app.PoolId})
			}
			table.Render()
		},
	}

	return cmd
}

func applicationUpdateCommand() *cobra.Command {
	var appID string
	var name string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update an application",
		Run: func(cmd *cobra.Command, args []string) {
			if appID == "" || name == "" {
				logrus.Error("missing required flags: --app-id, --name")
				return
			}

			client, err := authbase.NewClient("4000")
			if err != nil {
				logrus.Error(err)
				return
			}
			defer client.Close()

			_, err = client.UpdateApplication(tokenContext(), &v1.UpdateApplicationRequest{
				Id:   appID,
				Name: name,
			})
			if err != nil {
				logrus.Error(err)
				return
			}

			logrus.Info("Application updated")
		},
	}

	cmd.Flags().StringVar(&appID, "app-id", "", "The application id")
	cmd.Flags().StringVar(&name, "name", "", "The application name")

	return cmd
}

func applicationDeleteCommand() *cobra.Command {
	var appID string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an application",
		Run: func(cmd *cobra.Command, args []string) {
			if appID == "" {
				logrus.Error("missing required flags: --app-id")
				return
			}

			client, err := authbase.NewClient("4000")
			if err != nil {
				logrus.Error(err)
				return
			}
			defer client.Close()

			_, err = client.DeleteApplication(tokenContext(), &v1.DeleteApplicationRequest{Id: appID})
			if err != nil {
				logrus.Error(err)
				return
			}

			logrus.Info("Application deleted")
		},
	}

	cmd.Flags().StringVar(&appID, "app-id", "", "The application id")

	return cmd
}
