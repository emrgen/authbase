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
			verifyContext()
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

			logrus.Infof("user created successfully %v", "")
		},
	}

	bindContextFlags(command)

	command.Flags().StringVarP(&username, "username", "u", "", "username")
	command.Flags().StringVarP(&email, "email", "e", "", "email")

	return command
}

func listMemberCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "list",
		Short: "list user",
		Run: func(cmd *cobra.Command, args []string) {
			verifyContext()

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}
			defer client.Close()

			res, err := client.ListMember(tokenContext(), &v1.ListMemberRequest{
				OrganizationId: OrganizationId,
			})
			if err != nil {
				logrus.Errorf("failed to list member: %v", err)
				return
			}

			logrus.Infof("member list for organization %v", OrganizationId)
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
	var username string
	var email string

	command := &cobra.Command{
		Use:   "update",
		Short: "update user",
		Run: func(cmd *cobra.Command, args []string) {
			verifyContext()

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

			logrus.Infof("user updated successfully %v", "")
		},
	}

	bindContextFlags(command)
	command.Flags().StringVarP(&username, "username", "u", "", "username")
	command.Flags().StringVarP(&email, "email", "e", "", "email")

	return command
}

func deleteMemberCommand() *cobra.Command {
	var username string
	command := &cobra.Command{
		Use:   "delete",
		Short: "delete user",
		Run: func(cmd *cobra.Command, args []string) {
			verifyContext()
			if username == "" {
				logrus.Errorf("missing required flag: --username")
				return
			}

			_, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}

			logrus.Infof("user deleted successfully %v", "")
		},
	}

	bindContextFlags(command)
	command.Flags().StringVarP(&username, "username", "u", "", "username")

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

			if OrganizationId == "" {
				logrus.Errorf("missing required flag: --organization")
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
			defer client.Close()

			ctx := tokenContext()
			res, err := client.AddMember(ctx, &v1.AddMemberRequest{MemberId: userID, OrganizationId: OrganizationId})
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

func removeMemberCommand() *cobra.Command {
	var userID string

	command := &cobra.Command{
		Use:   "remove",
		Short: "remove member",
		Run: func(cmd *cobra.Command, args []string) {
			loadToken()

			if OrganizationId == "" {
				logrus.Errorf("missing required flag: --organization")
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
			res, err := client.RemoveMember(ctx, &v1.RemoveMemberRequest{MemberId: userID, OrganizationId: OrganizationId})
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
