package cmd

import (
	"context"
	"fmt"
	"github.com/emrgen/authbase"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var accessKeyCommand = &cobra.Command{
	Use:   "key",
	Short: "key commands",
}

func init() {
	accessKeyCommand.AddCommand(createTokenCommand())
	accessKeyCommand.AddCommand(listTokenCommand())
	accessKeyCommand.AddCommand(deleteTokenCommand())
	accessKeyCommand.AddCommand(refreshTokenCommand())
}

func createTokenCommand() *cobra.Command {
	var clientID string
	var clientSecret string
	var email string
	var password string
	var save bool

	// when client id is not provided, it will try to the token from the context
	command := &cobra.Command{
		Use:   "create",
		Short: "create access key",
		Run: func(cmd *cobra.Command, args []string) {
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

			req := &v1.CreateAccessKeyRequest{
				ClientId:     clientID,
				ClientSecret: clientSecret,
				Email:        email,
				Password:     password,
				ExpiresIn:    nil,
				Name:         nil,
			}

			ctx := context.Background()
			if clientSecret == "" && clientID == "" {
				ctx = tokenContext()
			}

			token, err := client.CreateAccessKey(ctx, req)
			if err != nil {
				logrus.Errorf("failed to create token %v", err)
				return
			}

			logrus.Infof("token created successfully")
			fmt.Printf("Token: %v\n", token.Token.AccessKey)

			if save {
				logrus.Infof("context updated with token")
				writeContext(Context{
					Token: token.Token.AccessKey,
				})
			}
		},
	}

	command.Flags().StringVarP(&clientID, "client-id", "c", "", "client id")
	command.Flags().StringVarP(&clientSecret, "client-secret", "s", "", "client secret")
	command.Flags().StringVarP(&email, "email", "e", "", "user name")
	command.Flags().StringVarP(&password, "password", "p", "", "password")
	command.Flags().BoolVarP(&save, "save", "v", false, "set context")

	return command
}

func listTokenCommand() *cobra.Command {
	var projectID string
	var userId string
	var password string

	command := &cobra.Command{
		Use:   "list",
		Short: "list token",
		Run: func(cmd *cobra.Command, args []string) {
			loadToken()
			if Token == "" {
				logrus.Error("missing required flags: --token")
				return
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("error creating client: %v", err)
				return
			}

			req := &v1.ListAccessKeysRequest{}

			if projectID != "" {
				req.ProjectId = &projectID
			}

			if userId != "" {
				req.AccountId = &userId
			}

			res, err := client.ListAccessKeys(tokenContext(), req)
			if err != nil {
				logrus.Errorf("failed to list tokens %v", err)
				return
			}

			for _, token := range res.Tokens {
				fmt.Printf("AccountID: %v\n", token.AccountId)
				fmt.Printf("PoolId: %v\n", token.PoolId)
				fmt.Printf("ID: %v\n", token.Id)
				fmt.Printf("Token: %v\n", token.AccessKey)
				fmt.Printf("------\n")
			}
		},
	}

	command.Flags().StringVarP(&projectID, "projectID", "r", "", "project id")
	command.Flags().StringVarP(&userId, "userId", "u", "", "account id")
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

			_, err = client.DeleteAccessKey(context.Background(), &v1.DeleteAccessKeyRequest{
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
