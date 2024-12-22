package cmd

import (
	"context"
	"github.com/emrgen/authbase"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var orgCommand = &cobra.Command{
	Use:   "org",
	Short: "org commands",
}

func init() {
	orgCommand.AddCommand(createOrgCommand())
	orgCommand.AddCommand(listOrgCommand())
	orgCommand.AddCommand(updateOrgCommand())
	orgCommand.AddCommand(deleteOrgCommand())
}

func createOrgCommand() *cobra.Command {
	var organization string
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

			if organization == "" {
				logrus.Errorf("missing required flags: --organization")
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
				logrus.Infof("creating organization without password")
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("error creating client: %v", err)
				return
			}
			defer client.Close()

			ctx := tokenContext(Token)

			// create the master organization
			if master {
				logrus.Infof("creating master organization")
				organization, err := client.CreateAdminOrganization(ctx, &v1.CreateAdminOrganizationRequest{
					Name:     organization,
					Username: username,
					Email:    email,
					Password: &password,
				})
				if err != nil {
					logrus.Errorf("error creating master organization: %v", err)
					return
				}
				logrus.Infof("master organization created: %v", organization)
			} else {
				organization, err := client.CreateOrganization(ctx, &v1.CreateOrganizationRequest{
					Name:     organization,
					Username: username,
					Email:    email,
					Password: &password,
				})
				if err != nil {
					logrus.Errorf("error creating organization: %v", err)
					return
				}
				logrus.Infof("organization created: %v", organization)
			}
		},
	}

	command.Flags().StringVarP(&organization, "organization", "o", "", "organization of the organization")
	command.Flags().StringVarP(&username, "username", "u", "", "username of the organization")
	command.Flags().StringVarP(&email, "email", "e", "", "email of the organization")
	command.Flags().StringVarP(&password, "password", "p", "", "password of the organization")
	command.Flags().BoolVarP(&master, "master", "m", false, "master organization")

	return command
}

func listOrgCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "list",
		Short: "list org",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("error creating client: %v", err)
				return
			}

			organizations, err := client.ListOrganizations(context.Background(), &v1.ListOrganizationsRequest{})
			if err != nil {
				logrus.Errorf("error listing organizations: %v", err)
				return
			}

			// print the organizations in a table
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"#", "ID", "Name", "Owner ID", "Master"})
			for i, organization := range organizations.Organizations {
				table.Append([]string{strconv.FormatInt(int64(i+1), 10), organization.Id, organization.Name, organization.OwnerId, strconv.FormatBool(organization.Master)})
			}

			table.Render()
		},
	}

	return command
}

func updateOrgCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "update",
		Short: "update org",
	}

	return command
}

func deleteOrgCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "delete",
		Short: "delete org",
	}

	return command
}
