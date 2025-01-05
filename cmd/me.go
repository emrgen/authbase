package cmd

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(whoamiCommand())
}

func whoamiCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "whoami",
		Short: "config information",
		Run: func(cmd *cobra.Command, args []string) {
			loadToken()

			if Token == "" {
				logrus.Errorf("token is not configured")
				return
			}

			// decode token
			token, _, err := jwt.NewParser().ParseUnverified(Token, jwt.MapClaims{})
			if err != nil {
				logrus.Errorf("failed to parse token: %v", err)
				return
			}

			claim := token.Claims.(jwt.MapClaims)

			// print claims
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Project ID", "User ID", "Email", "Username"})
			table.Append([]string{claim["project_id"].(string), claim["user_id"].(string), claim["email"].(string), claim["username"].(string)})
			table.Render()
		},
	}

	return command
}
