package cmd

import (
	"context"
	"github.com/emrgen/authbase"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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

	command := &cobra.Command{
		Use:   "create",
		Short: "create org",
		Run: func(cmd *cobra.Command, args []string) {
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

			organization, err := client.CreateOrganization(context.Background(), &v1.CreateOrganizationRequest{
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
		},
	}

	command.Flags().StringVarP(&organization, "organization", "o", "", "organization of the organization")
	command.Flags().StringVarP(&username, "username", "u", "", "username of the organization")
	command.Flags().StringVarP(&email, "email", "e", "", "email of the organization")
	command.Flags().StringVarP(&password, "password", "p", "", "password of the organization")

	return command
}

func listOrgCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "list",
		Short: "list org",
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
