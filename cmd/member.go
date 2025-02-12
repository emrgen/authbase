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

var memberCommand = &cobra.Command{
	Use:   "member",
	Short: "member commands",
}

func init() {
	memberCommand.AddCommand(addMemberCommand())
	memberCommand.AddCommand(removeMemberCommand())
	memberCommand.AddCommand(createMemberCommand())
	memberCommand.AddCommand(listMemberCommand())
	memberCommand.AddCommand(updateMemberCommand())
	memberCommand.AddCommand(deleteMemberCommand())
}

func createMemberCommand() *cobra.Command {
	var username string
	var email string
	var projectID string

	command := &cobra.Command{
		Use:   "create",
		Short: "create member",
		Run: func(cmd *cobra.Command, args []string) {
			loadToken()
			assertToken()

			if projectID == "" {
				logrus.Errorf("missing required flag: --project")
				return
			}

			if username == "" {
				logrus.Errorf("missing required flag: --username")
				return
			}

			if email == "" {
				logrus.Errorf("missing required flag: --email")
				return
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}

			member, err := client.CreateProjectMember(tokenContext(), &v1.CreateProjectMemberRequest{
				ProjectId:  projectID,
				Username:   username,
				Email:      email,
				Permission: v1.Permission_VIEWER,
			})
			if err != nil {
				logrus.Errorf("failed to create member: %v", err)
				return
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID"})
			table.Append([]string{member.Id})
			table.Render()
		},
	}

	bindContextFlags(command)

	command.Flags().StringVarP(&projectID, "project-id", "r", "", "project id")
	command.Flags().StringVarP(&username, "username", "u", "", "username")
	command.Flags().StringVarP(&email, "email", "e", "", "email")

	return command
}

func addMemberCommand() *cobra.Command {
	var projectID string
	var userID string
	var permission uint32

	command := &cobra.Command{
		Use:   "add",
		Short: "add member",
		Run: func(cmd *cobra.Command, args []string) {
			loadToken()

			if projectID == "" {
				logrus.Errorf("missing required flag: --project")
				return
			}

			if userID == "" {
				logrus.Errorf("missing required flag: --user-id")
				return
			}

			_, ok := v1.Permission_name[int32(permission)]
			if !ok {
				logrus.Errorf("invalid permission: %v", permission)
				return
			}
			if permission == 0 {
				permission = uint32(v1.Permission_NONE)
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}
			defer client.Close()

			ctx := tokenContext()
			res, err := client.AddProjectMember(ctx, &v1.AddProjectMemberRequest{
				MemberId:   userID,
				ProjectId:  projectID,
				Permission: v1.Permission(permission),
			})
			if err != nil {
				logrus.Errorf("failed to add member: %v", err)
				return
			}

			logrus.Infof("member added successfully %v", res)
		},
	}

	bindContextFlags(command)
	command.Flags().StringVarP(&userID, "user-id", "u", "", "user id")
	command.Flags().Uint32VarP(&permission, "permission", "p", 0, "permission")

	return command
}

func listMemberCommand() *cobra.Command {
	var projectID string

	command := &cobra.Command{
		Use:   "list",
		Short: "list members",
		Run: func(cmd *cobra.Command, args []string) {
			loadToken()
			assertToken()

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}
			defer client.Close()

			res, err := client.ListProjectMember(tokenContext(), &v1.ListProjectMemberRequest{
				ProjectId: projectID,
			})
			if err != nil {
				logrus.Errorf("failed to list member: %v", err)
				return
			}

			logrus.Infof("member list for project: %v", projectID)
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"#", "ID", "Username", "Email", "Permission"})
			for i, member := range res.Members {
				perm, ok := v1.Permission_name[int32(member.Permission)]
				if !ok {
					logrus.Errorf("failed to get permission name: %v", member.Permission)
					perm = "ERROR"
				}

				table.Append([]string{
					strconv.Itoa(i + 1),
					member.Id, member.Username, member.Email, perm,
				})
			}
			table.Render()
		},
	}

	command.Flags().StringVarP(&projectID, "project-id", "r", "", "project id")

	return command
}

func updateMemberCommand() *cobra.Command {
	var userID string
	var permission uint32

	command := &cobra.Command{
		Use:   "update",
		Short: "update user",
		Run: func(cmd *cobra.Command, args []string) {
			loadToken()
			if Token == "" {
				logrus.Errorf("missing required flag: --token")
				return
			}

			if userID == "" {
				logrus.Errorf("missing required flag: --user-id")
				return
			}

			_, ok := v1.Permission_name[int32(permission)]
			if !ok {
				logrus.Errorf("invalid permission: %v", permission)
				return
			}
			if permission == 0 {
				permission = uint32(v1.Permission_NONE)
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}
			defer client.Close()

			_, err = client.UpdateProjectMember(tokenContext(), &v1.UpdateProjectMemberRequest{
				MemberId:   userID,
				Permission: v1.Permission(permission),
			})
			if err != nil {
				logrus.Errorf("failed to update member: %v", err)
				return
			}

			logrus.Infof("member updated successfully %v", "")
		},
	}

	bindContextFlags(command)
	command.Flags().StringVarP(&userID, "user-id", "u", "", "user id")
	command.Flags().Uint32VarP(&permission, "permission", "p", 0, "permission")

	return command
}

func deleteMemberCommand() *cobra.Command {
	var projectID string
	var userID string

	command := &cobra.Command{
		Use:   "delete",
		Short: "delete user",
		Run: func(cmd *cobra.Command, args []string) {
			loadToken()

			if projectID == "" {
				logrus.Errorf("missing required flag: --project")
				return
			}

			if userID == "" {
				logrus.Errorf("missing required flag: --username")
				return
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}
			defer client.Close()

			_, err = client.RemoveProjectMember(tokenContext(), &v1.RemoveProjectMemberRequest{
				ProjectId: projectID,
				MemberId:  userID,
			})

			logrus.Infof("org member removed successfully")
		},
	}

	bindContextFlags(command)
	command.Flags().StringVarP(&userID, "user-id", "u", "", "user id")

	return command
}

func removeMemberCommand() *cobra.Command {
	var projectID string
	var userID string

	command := &cobra.Command{
		Use:   "remove",
		Short: "remove member",
		Run: func(cmd *cobra.Command, args []string) {
			loadToken()

			if projectID == "" {
				logrus.Errorf("missing required flag: --project")
				return
			}

			if userID == "" {
				logrus.Errorf("missing required flag: --user-id")
				return
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}

			ctx := tokenContext()
			res, err := client.RemoveProjectMember(ctx, &v1.RemoveProjectMemberRequest{MemberId: userID, ProjectId: projectID})
			if err != nil {
				logrus.Errorf("failed to remove member: %v", err)
				return
			}

			logrus.Infof("member removed successfully %v", res)
		},
	}

	bindContextFlags(command)
	command.Flags().StringVarP(&userID, "user-id", "u", "", "user id")

	return command
}
