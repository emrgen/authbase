package cmd

import (
	"fmt"
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
var groupMemberCommand = &cobra.Command{
	Use:   "member",
	Short: "Group member commands",
}

var groupRoleCommand = &cobra.Command{
	Use:   "role",
	Short: "Group role commands",
}

func init() {
	groupCommand.AddCommand(groupCreateCommand())
	groupCommand.AddCommand(groupGetCommand())
	groupCommand.AddCommand(groupListCommand())
	groupCommand.AddCommand(groupUpdateCommand())
	groupCommand.AddCommand(groupDeleteCommand())
	groupCommand.AddCommand(groupAddMemberCommand())
	groupCommand.AddCommand(groupRemoveMemberCommand())

	groupCommand.AddCommand(groupMemberCommand)
	groupMemberCommand.AddCommand(groupAddMemberCommand())
	groupMemberCommand.AddCommand(groupRemoveMemberCommand())
	groupMemberCommand.AddCommand(groupListMemberCommand())

	groupCommand.AddCommand(groupRoleCommand)
	groupRoleCommand.AddCommand(groupAddRoleCommand())
	groupRoleCommand.AddCommand(groupRemoveRoleCommand())
}

func groupCreateCommand() *cobra.Command {
	var name string
	var poolID string
	var scopes []string

	command := &cobra.Command{
		Use:   "create",
		Short: "Create a group",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := authbase.NewClient("4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}
			defer client.Close()

			if poolID == "" {
				poolID = getAccountPoolID(client)
			}

			if name == "" {
				logrus.Errorf("missing required flag: --name")
				return
			}

			res, err := client.CreateGroup(tokenContext(), &v1.CreateGroupRequest{
				PoolId:    poolID,
				Name:      name,
				RoleNames: scopes,
			})
			if err != nil {
				logrus.Errorf("failed to create group: %v", err)
				return
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Pool ID", "Name", "Scopes"})
			roleNames := make([]string, 0)
			for _, role := range res.Group.Roles {
				roleNames = append(roleNames, role.Name)
			}
			table.Append([]string{res.Group.Id, res.Group.PoolId, res.Group.Name, strings.Join(roleNames, ",")})
			table.Render()
		},
	}

	command.Flags().StringVarP(&poolID, "pool-id", "p", "", "Pool ID")
	command.Flags().StringVarP(&name, "name", "n", "", "Group name")
	command.Flags().StringSliceVarP(&scopes, "scopes", "s", []string{}, "Group scopes")

	return command
}

func groupGetCommand() *cobra.Command {
	var groupID string

	command := &cobra.Command{
		Use:   "get",
		Short: "Get a group",
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

			res, err := client.GetGroup(tokenContext(), &v1.GetGroupRequest{
				GroupId: groupID,
			})
			if err != nil {
				logrus.Errorf("failed to get group: %v", err)
				return
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Pool ID", "Name", "Scopes"})
			roleNames := make([]string, 0)
			for _, role := range res.Group.Roles {
				roleNames = append(roleNames, role.Name)
			}
			table.Append([]string{res.Group.Id, res.Group.PoolId, res.Group.Name, strings.Join(roleNames, ",")})
			table.Render()
		},
	}

	command.Flags().StringVarP(&groupID, "group-id", "g", "", "Group ID")

	return command
}

func groupListCommand() *cobra.Command {
	var poolID string
	var accountID string

	command := &cobra.Command{
		Use:   "list",
		Short: "List groups",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := authbase.NewClient("4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}
			defer client.Close()
			if poolID == "" {
				poolID = getAccountPoolID(client)
			}

			req := &v1.ListGroupsRequest{}

			if accountID != "" {
				req.AccountId = &accountID
			}

			if accountID == "" && poolID != "" {
				req.PoolId = &poolID
			}

			if req.PoolId == nil && req.AccountId == nil {
				logrus.Errorf("missing required flag: --pool-id or --account-id")
				return
			}

			res, err := client.ListGroups(tokenContext(), req)
			if err != nil {
				logrus.Errorf("failed to list groups: %v", err)
				return
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Pool ID", "Name", "Scopes"})
			for _, group := range res.Groups {
				roleNames := make([]string, 0)
				for _, role := range group.Roles {
					roleNames = append(roleNames, role.Name)
				}
				table.Append([]string{group.Id, group.PoolId, group.Name, strings.Join(roleNames, ",")})
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
				req.RoleNames = scopes
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

func groupAddMemberCommand() *cobra.Command {
	var groupID string
	var accountID string

	command := &cobra.Command{
		Use:   "add",
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

			fmt.Print("member added to group\n")
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
		Use:   "remove",
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

			fmt.Print("member removed from group\n")
		},
	}

	command.Flags().StringVarP(&groupID, "group-id", "g", "", "Group ID")
	command.Flags().StringVarP(&accountID, "account-id", "a", "", "Account ID")

	return command
}

func groupListMemberCommand() *cobra.Command {
	var groupID string

	command := &cobra.Command{
		Use:   "list",
		Short: "List group members",
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

			res, err := client.ListGroupMembers(tokenContext(), &v1.ListGroupMembersRequest{
				GroupId: groupID,
			})
			if err != nil {
				logrus.Errorf("failed to list group members: %v", err)
				return
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Account ID", "Name", "Email", "Scopes"})
			for _, member := range res.Members {
				roleNames := make([]string, 0)
				for _, role := range res.Roles {
					roleNames = append(roleNames, role.Name)
				}
				table.Append([]string{member.Id, member.VisibleName, member.Email, strings.Join(roleNames, ",")})
			}
			table.Render()
		},
	}

	command.Flags().StringVarP(&groupID, "group-id", "g", "", "Group ID")

	return command
}

func groupAddRoleCommand() *cobra.Command {
	var groupID string
	var roleName string

	command := &cobra.Command{
		Use:   "add",
		Short: "Add a role to a group",
		Run: func(cmd *cobra.Command, args []string) {
			if groupID == "" {
				logrus.Errorf("missing required flag: --group-id")
				return
			}
			if roleName == "" {
				logrus.Errorf("missing required flag: --role-name")
				return
			}

			client, err := authbase.NewClient("4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}
			defer client.Close()

			_, err = client.AddRole(tokenContext(), &v1.AddRoleRequest{
				GroupId:  groupID,
				RoleName: roleName,
			})
			if err != nil {
				logrus.Errorf("failed to add role to group: %v", err)
				return
			}

			fmt.Print("role added to group\n")
		},
	}

	command.Flags().StringVarP(&groupID, "group-id", "g", "", "Group ID")
	command.Flags().StringVarP(&roleName, "role-name", "r", "", "Role Name")

	return command
}

func groupRemoveRoleCommand() *cobra.Command {
	var groupID string
	var roleName string

	command := &cobra.Command{
		Use:   "remove",
		Short: "Remove a role from a group",
		Run: func(cmd *cobra.Command, args []string) {
			if groupID == "" {
				logrus.Errorf("missing required flag: --group-id")
				return
			}
			if roleName == "" {
				logrus.Errorf("missing required flag: --role-name")
				return
			}

			client, err := authbase.NewClient("4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}
			defer client.Close()

			_, err = client.RemoveRole(tokenContext(), &v1.RemoveRoleRequest{
				GroupId:  groupID,
				RoleName: roleName,
			})
			if err != nil {
				logrus.Errorf("failed to remove role from group: %v", err)
				return
			}

			fmt.Print("role removed from group\n")
		},
	}

	command.Flags().StringVarP(&groupID, "group-id", "g", "", "Group ID")
	command.Flags().StringVarP(&roleName, "role-name", "r", "", "Role Name")

	return command
}
