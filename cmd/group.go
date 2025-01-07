package cmd

import (
	"github.com/emrgen/authbase"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var groupCommand = &cobra.Command{
	Use:   "group",
	Short: "Group commands",
}

func init() {
	groupCommand.AddCommand(groupCreateCommand())
	groupCommand.AddCommand(groupListCommand())
	groupCommand.AddCommand(groupUpdateCommand())
	groupCommand.AddCommand(groupAddMemberCommand())
	groupCommand.AddCommand(groupRemoveMemberCommand())
	groupCommand.AddCommand(groupDeleteCommand())
}

func groupCreateCommand() *cobra.Command {
	var name string
	var poolID string
	var scopes []string

	command := &cobra.Command{
		Use:   "create",
		Short: "Create a group",
		Run: func(cmd *cobra.Command, args []string) {
			if poolID == "" {
				logrus.Errorf("missing required flag: --pool-id")
				return
			}

			if name == "" {
				logrus.Errorf("missing required flag: --name")
				return
			}

			client, err := authbase.NewClient("4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}
			defer client.Close()

			res, err := client.CreateGroup(tokenContext(), &v1.CreateGroupRequest{
				PoolId: poolID,
				Name:   name,
				Scopes: scopes,
			})
			if err != nil {
				logrus.Errorf("failed to create group: %v", err)
				return
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Pool ID", "Name", "Scopes"})
			table.Append([]string{res.Group.Id, res.Group.PoolId, res.Group.Name, strings.Join(res.Group.Scopes, ",")})
			table.Render()
		},
	}

	command.Flags().StringVarP(&poolID, "pool-id", "p", "", "Pool ID")
	command.Flags().StringVarP(&name, "name", "n", "", "Group name")
	command.Flags().StringSliceVarP(&scopes, "scopes", "s", []string{}, "Group scopes")

	return command
}

func groupListCommand() *cobra.Command {
	var poolID string
	var accountID string

	command := &cobra.Command{
		Use:   "list",
		Short: "List groups",
		Run: func(cmd *cobra.Command, args []string) {
			if poolID == "" && accountID == "" {
				logrus.Errorf("pool-id or account-id is required")
				return
			}

			client, err := authbase.NewClient("4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}
			defer client.Close()

			req := &v1.ListGroupsRequest{}
			if poolID != "" {
				req.PoolId = &poolID
			}
			if accountID != "" {
				req.AccountId = &accountID
			}

			res, err := client.ListGroups(tokenContext(), req)
			if err != nil {
				logrus.Errorf("failed to list groups: %v", err)
				return
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Pool ID", "Name", "Scopes"})
			for _, group := range res.Groups {
				table.Append([]string{group.Id, group.PoolId, group.Name, strings.Join(group.Scopes, ",")})
			}
			table.Render()
		},
	}

	command.Flags().StringVarP(&poolID, "pool-id", "p", "", "Pool ID")
	command.Flags().StringVarP(&accountID, "account-id", "a", "", "Account ID")

	return command
}

func groupUpdateCommand() *cobra.Command {
	var groupID string
	var name string
	var scopes []string

	command := &cobra.Command{
		Use:   "update",
		Short: "Update a group",
		Run: func(cmd *cobra.Command, args []string) {
			if groupID == "" {
				logrus.Errorf("missing required flag: --group-id")
				return
			}

			client, err := authbase.NewClient("4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}
			defer client.Close()

			req := &v1.UpdateGroupRequest{
				GroupId: groupID,
			}

			if name != "" {
				req.Name = name
			}

			if len(scopes) > 0 {
				req.Scopes = scopes
			}

			_, err = client.UpdateGroup(tokenContext(), req)
			if err != nil {
				logrus.Errorf("failed to update group: %v", err)
				return
			}

			logrus.Infof("group updated")
		},
	}

	command.Flags().StringVarP(&groupID, "group-id", "g", "", "Group ID")
	command.Flags().StringVarP(&name, "name", "n", "", "Group name")
	command.Flags().StringSliceVarP(&scopes, "scopes", "s", []string{}, "Group scopes")

	return command
}

func groupAddMemberCommand() *cobra.Command {
	var groupID string
	var accountID string

	command := &cobra.Command{
		Use:   "add-member",
		Short: "Add a member to a group",
		Run: func(cmd *cobra.Command, args []string) {
			if groupID == "" {
				logrus.Errorf("missing required flag: --group-id")
				return
			}
			if accountID == "" {
				logrus.Errorf("missing required flag: --account-id")
				return
			}

			client, err := authbase.NewClient("4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}
			defer client.Close()

			_, err = client.AddGroupMember(tokenContext(), &v1.AddGroupMemberRequest{
				GroupId:   groupID,
				AccountId: accountID,
			})
			if err != nil {
				logrus.Errorf("failed to add member to group: %v", err)
				return
			}

			cmd.Printf("member added to group")
		},
	}

	command.Flags().StringVarP(&groupID, "group-id", "g", "", "Group ID")
	command.Flags().StringVarP(&accountID, "account-id", "a", "", "Account ID")

	return command
}

func groupRemoveMemberCommand() *cobra.Command {
	var groupID string
	var accountID string

	command := &cobra.Command{
		Use:   "remove-member",
		Short: "Remove a member from a group",
		Run: func(cmd *cobra.Command, args []string) {
			if groupID == "" {
				logrus.Errorf("missing required flag: --group-id")
				return
			}
			if accountID == "" {
				logrus.Errorf("missing required flag: --account-id")
				return
			}

			client, err := authbase.NewClient("4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}
			defer client.Close()

			_, err = client.RemoveGroupMember(tokenContext(), &v1.RemoveGroupMemberRequest{
				GroupId:   groupID,
				AccountId: accountID,
			})
			if err != nil {
				logrus.Errorf("failed to add member to group: %v", err)
				return
			}

			cmd.Printf("member removed from group")
		},
	}

	command.Flags().StringVarP(&groupID, "group-id", "g", "", "Group ID")
	command.Flags().StringVarP(&accountID, "account-id", "a", "", "Account ID")

	return command
}

func groupDeleteCommand() *cobra.Command {
	var groupID string

	command := &cobra.Command{
		Use:   "delete",
		Short: "Delete a group",
		Run: func(cmd *cobra.Command, args []string) {
			if groupID == "" {
				logrus.Errorf("missing required flag: --group-id")
				return
			}

			client, err := authbase.NewClient("4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}
			defer client.Close()

			_, err = client.DeleteGroup(tokenContext(), &v1.DeleteGroupRequest{
				GroupId: groupID,
			})
			if err != nil {
				logrus.Errorf("failed to delete group: %v", err)
				return
			}

			logrus.Infof("group deleted")
		},
	}

	command.Flags().StringVarP(&groupID, "group-id", "g", "", "Group ID")

	return command
}
