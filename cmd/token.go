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
	tokenCommand.AddCommand(refreshTokenCommand())
}

func refreshTokenCommand() *cobra.Command {
	var refreshToken string

	command := &cobra.Command{
		Use:   "refresh",
		Short: "refresh token",
		Run: func(cmd *cobra.Command, args []string) {
			if refreshToken == "" {
				cmd.Help()
				return
			}

			client, err := authbase.NewClient("4000")
			if err != nil {
				logrus.Error(err)
				return
			}
			defer client.Close()

			res, err := client.Refresh(context.TODO(), &v1.RefreshRequest{
				RefreshToken: refreshToken,
			})
			if err != nil {
				logrus.Error(err)
				return
			}

			cmd.Printf("----------\n")
			cmd.Printf("AccessToken: %v\n", res.Tokens.AccessToken)
			cmd.Printf("RefreshToken: %v\n", res.Tokens.RefreshToken)
		},
	}

	command.Flags().StringVarP(&refreshToken, "refresh-token", "r", "", "refresh token")

	return command
}
