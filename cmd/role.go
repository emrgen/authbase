package cmd

import (
	"github.com/emrgen/authbase"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var roleCommand = &cobra.Command{
	Use:   "role",
	Short: "Role commands",
}

func init() {
	roleCommand.AddCommand(roleCreateCommand())
	roleCommand.AddCommand(roleListCommand())
	roleCommand.AddCommand(roleUpdateCommand())
	roleCommand.AddCommand(roleDeleteCommand())
}

func roleCreateCommand() *cobra.Command {
	var poolID string
	var name string

	command := &cobra.Command{
		Use:   "create",
		Short: "Create a role",
		Run: func(cmd *cobra.Command, args []string) {
			if poolID == "" {
				logrus.Errorf("missing required flag: --pool-id")
				return
			}

			if name == "" {
				logrus.Errorf("missing required flag: --name")
				return
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}

			res, err := client.CreateRole(tokenContext(), &v1.CreateRoleRequest{
				Name: name,
			})
			if err != nil {
				logrus.Errorf("failed to create role: %v", err)
				return
			}

			logrus.Infof("role created successfully %v", res.GetRole().GetName())
		},
	}

	bindContextFlags(command)
	command.Flags().StringVarP(&poolID, "pool-id", "p", "", "pool id")
	command.Flags().StringVarP(&name, "name", "n", "", "name")

	return command
}

func roleListCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "list",
		Short: "List roles",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}

			res, err := client.ListRoles(tokenContext(), &v1.ListRolesRequest{})
			if err != nil {
				logrus.Errorf("failed to list roles: %v", err)
				return
			}

			for _, role := range res.GetRoles() {
				logrus.Infof("role: %v", role.GetName())
			}
		},
	}

	bindContextFlags(command)

	return command
}

func roleUpdateCommand() *cobra.Command {
	var name string

	command := &cobra.Command{
		Use:   "update",
		Short: "Update a role",
		Run: func(cmd *cobra.Command, args []string) {
			if name == "" {
				logrus.Errorf("missing required flag: --name")
				return
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}

			res, err := client.UpdateRole(tokenContext(), &v1.UpdateRoleRequest{
				Name: name,
			})
			if err != nil {
				logrus.Errorf("failed to update role: %v", err)
				return
			}

			logrus.Infof("role updated successfully %v", res.GetRole().GetName())
		},
	}

	bindContextFlags(command)
	command.Flags().StringVarP(&name, "name", "n", "", "name")

	return command
}

func roleDeleteCommand() *cobra.Command {
	var name string

	command := &cobra.Command{
		Use:   "delete",
		Short: "Delete a role",
		Run: func(cmd *cobra.Command, args []string) {
			if name == "" {
				logrus.Errorf("missing required flag: --name")
				return
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}

			res, err := client.DeleteRole(tokenContext(), &v1.DeleteRoleRequest{
				RoleName: name,
			})
			if err != nil {
				logrus.Errorf("failed to delete role: %v", err)
				return
			}

			logrus.Infof("role deleted successfully %v", res.GetMessage())
		},
	}

	bindContextFlags(command)
	command.Flags().StringVarP(&name, "name", "n", "", "name")

	return command
}
