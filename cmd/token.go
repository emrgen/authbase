package cmd

import (
	"context"
	"github.com/emrgen/authbase"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/google/uuid"
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
	var organizationId string
	var email string
	var password string
	var save bool

	command := &cobra.Command{
		Use:   "create",
		Short: "create token",
		Run: func(cmd *cobra.Command, args []string) {
			if organizationId == "" {
				logrus.Error("missing required flags: --organization")
				return
			}
			_, err := uuid.Parse(organizationId)
			if err != nil {
				logrus.Error("organization id must be a valid uuid")
				return
			}

			if email == "" {
				logrus.Error("missing required flags: --email")
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

			token, err := client.CreateToken(context.Background(), &v1.CreateTokenRequest{
				OrganizationId: organizationId,
				Email:          email,
				Password:       password,
				ExpiresIn:      nil,
				Name:           nil,
			})
			if err != nil {
				logrus.Errorf("failed to create token %v", err)
				return
			}

			logrus.Infof("token created successfully %v", token)
			if save {
				logrus.Infof("context updated with token")
				writeContext(Context{
					OrganizationId: organizationId,
					Token:          token.Token,
					Username:       email,
					Password:       password,
				})
			}
		},
	}

	command.Flags().StringVarP(&organizationId, "org-id", "o", "", "organization id")
	command.Flags().StringVarP(&email, "email", "e", "", "user name")
	command.Flags().StringVarP(&password, "password", "p", "", "password")
	command.Flags().BoolVarP(&save, "save", "s", false, "set context")

	return command
}

func listTokenCommand() *cobra.Command {
	var organizationId string
	var userId string
	var password string

	command := &cobra.Command{
		Use:   "list",
		Short: "list token",
		Run: func(cmd *cobra.Command, args []string) {
			if organizationId == "" {
				logrus.Error("missing required flags: --organizationId")
				return
			}

			if Token == "" {
				if userId == "" {
					logrus.Error("missing required flags: --userId")
					return
				}

				if password == "" {
					logrus.Error("missing required flags: --password")
					return
				}
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("error creating client: %v", err)
				return
			}

			res, err := client.ListTokens(context.Background(), &v1.ListTokensRequest{
				OrganizationId: organizationId,
				UserId:         userId,
			})
			if err != nil {
				logrus.Errorf("failed to list tokens %v", err)
				return
			}

			logrus.Infof("user token count: %v", res.Meta.Total)
			for _, token := range res.Tokens {
				logrus.Infof("token %v", token)
			}
		},
	}

	command.Flags().StringVarP(&organizationId, "organizationId", "o", "", "organizationId name")
	command.Flags().StringVarP(&userId, "userId", "u", "", "userId name")
	command.Flags().StringVarP(&password, "password", "p", "", "password")

	return command
}

func deleteTokenCommand() *cobra.Command {
	var tokeID string
	command := &cobra.Command{
		Use:   "delete",
		Short: "delete token",
		Run: func(cmd *cobra.Command, args []string) {
			if Token == "" {
				logrus.Error("missing required flags: --token")
				return
			}

			if tokeID == "" {
				logrus.Error("missing required flags: --token-id")
				return
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("error creating client: %v", err)
				return
			}

			_, err = client.DeleteToken(context.Background(), &v1.DeleteTokenRequest{
				Id: tokeID,
			})
			if err != nil {
				logrus.Errorf("failed to delete token")
				return
			}

			logrus.Infof("token deleted successfully")
		},
	}

	command.Flags().StringVarP(&tokeID, "token-id", "i", "", "token id")
	command.Flags().StringVarP(&Token, "token", "t", "", "token")

	return command
}
