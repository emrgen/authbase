package cmd

import (
	"context"
	"fmt"
	"github.com/emrgen/authbase"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/golang-jwt/jwt/v5"
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
	var projectId string
	var email string
	var password string
	var save bool

	command := &cobra.Command{
		Use:   "create",
		Short: "create token",
		Run: func(cmd *cobra.Command, args []string) {
			if projectId == "" {
				logrus.Error("missing required flags: --project")
				return
			}
			_, err := uuid.Parse(projectId)
			if err != nil {
				logrus.Error("project id must be a valid uuid")
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
				ProjectId: projectId,
				Email:     email,
				Password:  password,
				ExpiresIn: nil,
				Name:      nil,
			})
			if err != nil {
				logrus.Errorf("failed to create token %v", err)
				return
			}

			logrus.Infof("token created successfully")
			fmt.Printf("Token: %v\n", token.Token)

			if save {
				logrus.Infof("context updated with token")
				writeContext(Context{
					ProjectId: projectId,
					Token:     token.Token,
					Username:  email,
					Password:  password,
				})
			}
		},
	}

	command.Flags().StringVarP(&projectId, "project-id", "p", "", "project id")
	command.Flags().StringVarP(&email, "email", "e", "", "user name")
	command.Flags().StringVarP(&password, "password", "w", "", "password")
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
			loadToken()
			if Token == "" {
				logrus.Error("missing required flags: --token")
				return
			}

			if organizationId == "" {
				logrus.Error("missing required flags: --organizationId")
				return
			}

			if userId == "" {
				token, _, err := jwt.NewParser().ParseUnverified(Token, jwt.MapClaims{})
				if err != nil {
					logrus.Errorf("failed to parse token: %v", err)
					return
				}
				userId = token.Claims.(jwt.MapClaims)["user_id"].(string)
			}

			if userId == "" {
				logrus.Error("missing required flags: --userId")
				return
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("error creating client: %v", err)
				return
			}

			ctx := tokenContext()
			res, err := client.ListTokens(ctx, &v1.ListTokensRequest{
				ProjectId: organizationId,
				UserId:    userId,
			})
			if err != nil {
				logrus.Errorf("failed to list tokens %v", err)
				return
			}

			for _, token := range res.Tokens {
				fmt.Printf("UserID: %v\n", userId)
				fmt.Printf("ID: %v\n", token.Id)
				fmt.Printf("Token: %v\n", token.Token)
				fmt.Printf("------\n")
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
