package cmd

import (
	"github.com/emrgen/authbase"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strconv"
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
	var project string
	var password string
	var username string
	var email string
	var master bool

	command := &cobra.Command{
		Use:   "create",
		Short: "create org",
		Run: func(cmd *cobra.Command, args []string) {
			loadToken()

			if Token == "" {
				logrus.Errorf("missing required flags: --token")
				return
			}

			if project == "" {
				logrus.Errorf("missing required flags: --project")
				return
			}

			if email == "" {
				logrus.Errorf("missing required flags: --email")
				return
			}

			if username == "" {
				logrus.Errorf("missing required flags: --username")
				return
			}

			if password == "" {
				logrus.Infof("creating project without password")
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("error creating client: %v", err)
				return
			}
			defer client.Close()

			ctx := tokenContext()

			// create the master project
			if master {
				logrus.Infof("creating master project")
				project, err := client.CreateAdminProject(ctx, &v1.CreateAdminProjectRequest{
					Name:     project,
					Username: username,
					Email:    email,
					Password: &password,
				})
				if err != nil {
					logrus.Errorf("error creating master project: %v", err)
					return
				}
				logrus.Infof("master project created: %v", project)
			} else {
				project, err := client.CreateProject(ctx, &v1.CreateProjectRequest{
					Name:     project,
					Username: username,
					Email:    email,
					Password: &password,
				})
				if err != nil {
					logrus.Errorf("error creating project: %v", err)
					return
				}
				logrus.Infof("project created: %v", project)
			}
		},
	}

	command.Flags().StringVarP(&project, "project", "r", "", "project of the project")
	command.Flags().StringVarP(&username, "username", "u", "", "username of the project")
	command.Flags().StringVarP(&email, "email", "e", "", "email of the project")
	command.Flags().StringVarP(&password, "password", "p", "", "password of the project")
	command.Flags().BoolVarP(&master, "master", "m", false, "master project")

	return command
}

func getProjectCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "get",
		Short: "get project by id",
		Run: func(cmd *cobra.Command, args []string) {
			loadToken()
			if Token == "" {
				logrus.Errorf("missing required flags: --token")
				return
			}

			if ProjectId == "" {
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
				Id: ProjectId,
			})
			if err != nil {
				logrus.Errorf("error getting project: %v", err)
				return
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Name", "Owner ID", "Master"})
			table.Append([]string{res.Project.Id, res.Project.Name, res.Project.OwnerId, strconv.FormatBool(res.Project.Master)})
			table.Render()

		},
	}

	bindContextFlags(command)

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

			ctx := tokenContext()
			projects, err := client.ListProjects(ctx, &v1.ListProjectsRequest{})
			if err != nil {
				logrus.Errorf("error listing projects: %v", err)
				return
			}

			// print the projects in a table
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"#", "ID", "Name", "Owner ID", "Master"})
			for i, project := range projects.Projects {
				table.Append([]string{strconv.FormatInt(int64(i+1), 10), project.Id, project.Name, project.OwnerId, strconv.FormatBool(project.Master)})
			}

			table.Render()
		},
	}

	return command
}

func updateProjectCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "update",
		Short: "update org",
	}

	return command
}

func deleteProjectCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "delete",
		Short: "delete org",
	}

	return command
}
