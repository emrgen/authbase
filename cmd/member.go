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

	command := &cobra.Command{
		Use:   "create",
		Short: "create member",
		Run: func(cmd *cobra.Command, args []string) {
			loadToken()
			verifyToken()

			if username == "" {
				logrus.Errorf("missing required flag: --username")
				return
			}

			if email == "" {
				logrus.Errorf("missing required flag: --email")
				return
			}

			_, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}

			logrus.Infof("org member created successfully %v", "")
		},
	}

	bindContextFlags(command)

	command.Flags().StringVarP(&username, "username", "u", "", "username")
	command.Flags().StringVarP(&email, "email", "e", "", "email")

	return command
}

func addMemberCommand() *cobra.Command {
	var userID string
	var permission uint32

	command := &cobra.Command{
		Use:   "add",
		Short: "add member",
		Run: func(cmd *cobra.Command, args []string) {
			loadToken()

			if ProjectId == "" {
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
			res, err := client.AddMember(ctx, &v1.AddMemberRequest{
				MemberId:   userID,
				ProjectId:  ProjectId,
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
	command := &cobra.Command{
		Use:   "list",
		Short: "list user",
		Run: func(cmd *cobra.Command, args []string) {
			loadToken()
			verifyToken()

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}
			defer client.Close()

			res, err := client.ListMember(tokenContext(), &v1.ListMemberRequest{
				ProjectId: ProjectId,
			})
			if err != nil {
				logrus.Errorf("failed to list member: %v", err)
				return
			}

			logrus.Infof("member list for project: %v", ProjectId)
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"#", "ID", "Username", "Email", "Permission"})
			for i, member := range res.Members {
				perm, ok := v1.Permission_name[int32(member.Permission)]
				if !ok {
					logrus.Errorf("failed to get permission name: %v", member.Permission)
					perm = "ERROR"
				}

				v1.Permission_NONE.String()
				table.Append([]string{
					strconv.Itoa(i + 1),
					member.Id, member.Username, member.Email, perm,
				})
			}
			table.Render()
		},
	}

	bindContextFlags(command)
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

			_, err = client.UpdateMember(tokenContext(), &v1.UpdateMemberRequest{
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
	var userID string

	command := &cobra.Command{
		Use:   "delete",
		Short: "delete user",
		Run: func(cmd *cobra.Command, args []string) {
			loadToken()
			verifyToken()

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

			_, err = client.RemoveMember(tokenContext(), &v1.RemoveMemberRequest{
				ProjectId: ProjectId,
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
	var userID string

	command := &cobra.Command{
		Use:   "remove",
		Short: "remove member",
		Run: func(cmd *cobra.Command, args []string) {
			loadToken()

			if ProjectId == "" {
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
			res, err := client.RemoveMember(ctx, &v1.RemoveMemberRequest{MemberId: userID, ProjectId: ProjectId})
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
