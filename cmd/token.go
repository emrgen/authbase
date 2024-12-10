package cmd

import (
	"context"
	"github.com/emrgen/authbase"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var tokenCommand = &cobra.Command{
	Use:   "token",
	Short: "token commands",
}

func init() {
	tokenCommand.AddCommand(createTokenCommand())
	tokenCommand.AddCommand(listTokenCommand())
	tokenCommand.AddCommand(deleteTokenCommand())
}

func createTokenCommand() *cobra.Command {
	var organization string
	var user string
	var password string
	var save bool

	command := &cobra.Command{
		Use:   "create",
		Short: "create token",
		Run: func(cmd *cobra.Command, args []string) {
			if organization == "" {
				logrus.Error("missing required flags: --organization")
				return
			}

			if user == "" {
				logrus.Error("missing required flags: --user")
				return
			}

			if password == "" {
				logrus.Error("missing required flags: --password")
				return
			}

			_, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("error creating client: %v", err)
				return
			}

			logrus.Infof("token created successfully")
			if save {
				writeContext(Context{
					Organization: organization,
					Token:        "token",
					ExpireAt:     0,
				})
			}
		},
	}

	command.Flags().StringVarP(&organization, "organization", "o", "", "organization name")
	command.Flags().StringVarP(&user, "user", "u", "", "user name")
	command.Flags().StringVarP(&password, "password", "p", "", "password")
	command.Flags().BoolVarP(&save, "save", "s", false, "set context")

	return command
}

func listTokenCommand() *cobra.Command {
	var organization string
	var user string
	var password string

	command := &cobra.Command{
		Use:   "list",
		Short: "list token",
		Run: func(cmd *cobra.Command, args []string) {
			if organization == "" {
				logrus.Error("missing required flags: --organization")
				return
			}

			if user == "" {
				logrus.Error("missing required flags: --user")
				return
			}

			if password == "" {
				logrus.Error("missing required flags: --password")
				return
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("error creating client: %v", err)
				return
			}

			res, err := client.ListTokens(context.Background(), &v1.ListTokensRequest{
				OrganizationId: "",
			})
			if err != nil {
				return
			}

			logrus.Infof("token list")
			for _, token := range res.Tokens {
				logrus.Infof("token %v", token)
			}
		},
	}

	command.Flags().StringVarP(&organization, "organization", "o", "", "organization name")
	command.Flags().StringVarP(&user, "user", "u", "", "user name")
	command.Flags().StringVarP(&password, "password", "p", "", "password")

	return command
}

func deleteTokenCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "delete",
		Short: "delete token",
		Run: func(cmd *cobra.Command, args []string) {
			verifyContext()

			_, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("error creating client: %v", err)
				return
			}

			logrus.Infof("token deleted successfully")
		},
	}

	command.Flags().StringVarP(&Organization, "organization", "o", "", "organization name")
	command.Flags().StringVarP(&Token, "token", "t", "", "token")

	return command
}
