package cmd

import (
	"fmt"
	"github.com/emrgen/authbase"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
)

var projectCommand = &cobra.Command{
	Use:   "project",
	Short: "project commands",
}

func init() {
	projectCommand.AddCommand(createProjectCommand())
	projectCommand.AddCommand(listProjectCommand())
	projectCommand.AddCommand(updateProjectCommand())
	projectCommand.AddCommand(deleteProjectCommand())
	projectCommand.AddCommand(getProjectCommand())
}

func createProjectCommand() *cobra.Command {
	var projectName string
	var password string
	var visibleName string
	var email string

	command := &cobra.Command{
		Use:   "create",
		Short: "create org",
		Run: func(cmd *cobra.Command, args []string) {
			loadToken()

			if Token == "" {
				logrus.Errorf("missing required flags: --token")
				return
			}

			if projectName == "" {
				logrus.Errorf("missing required flags: --name")
				return
			}

			if email == "" {
				logrus.Errorf("missing required flags: --email")
				return
			}

			if password == "" {
				logrus.Infof("creating projectName without password")
			}

			if visibleName == "" {
				visibleName = strings.Split(email, "@")[0]
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("error creating client: %v", err)
				return
			}
			defer client.Close()

			res, err := client.CreateProject(tokenContext(), &v1.CreateProjectRequest{
				Name:        projectName,
				VisibleName: visibleName,
				Email:       email,
				Password:    &password,
			})
			if err != nil {
				logrus.Errorf("error creating project: %v", err)
				return
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Name", "Pool ID", "Client ID"})
			table.Append([]string{res.Project.Id, res.Project.Name, res.Project.PoolId, res.Client.Id})
			table.Render()
		},
	}

	command.Flags().StringVarP(&projectName, "projectName-name", "n", "", "name of the project")
	command.Flags().StringVarP(&visibleName, "visibleName", "v", "", "visible name of the project owner")
	command.Flags().StringVarP(&email, "email", "e", "", "email of the projectName")
	command.Flags().StringVarP(&password, "password", "p", "", "password of the project owner")
	return command
}

func getProjectCommand() *cobra.Command {
	var projectID string
	command := &cobra.Command{
		Use:   "get",
		Short: "get project by id",
		Run: func(cmd *cobra.Command, args []string) {
			loadToken()

			if Token == "" {
				logrus.Errorf("missing required flags: --token")
				return
			}

			if projectID == "" {
				logrus.Errorf("missing required flags: --orgName")
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("error creating client: %v", err)
				return
			}
			defer client.Close()

			ctx := tokenContext()
			res, err := client.GetProject(ctx, &v1.GetProjectRequest{
				Id: projectID,
			})
			if err != nil {
				logrus.Errorf("error getting project: %v", err)
				return
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Name", "Owner ID", "Master", "Members", "Users"})
			table.Append([]string{res.Project.Id, res.Project.Name, res.Project.OwnerId,
				strconv.FormatBool(res.Project.Master),
				strconv.FormatInt(int64(res.Members), 10),
				strconv.FormatInt(int64(res.Accounts), 10),
			})
			table.Render()

		},
	}

	command.Flags().StringVarP(&projectID, "project", "r", "", "project id")

	return command
}

func listProjectCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "list",
		Short: "list projects",
		Run: func(cmd *cobra.Command, args []string) {
			loadToken()

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("error creating client: %v", err)
				return
			}
			defer client.Close()

			ctx := tokenContext()
			res, err := client.ListProjects(ctx, &v1.ListProjectsRequest{})
			if err != nil {
				logrus.Errorf("error listing res: %v", err)
				return
			}

			// print the res in a table
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"#", "ID", "Name", "Owner ID", "Master"})
			for i, project := range res.Projects {
				table.Append([]string{strconv.FormatInt(int64(i+1), 10), project.Id, project.Name, project.OwnerId, strconv.FormatBool(project.Master)})
			}
			table.Render()
			fmt.Printf("projects => showing: %v, page: %v, total: %v\n", len(res.Projects), res.Meta.Page, res.Meta.Total)
		},
	}

	return command
}

func updateProjectCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "update",
		Short: "update project",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	return command
}

func deleteProjectCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "delete",
		Short: "delete project",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	return command
}
