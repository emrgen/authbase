package cmd

import (
	"github.com/emrgen/authbase"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

func init() {
	rootCmd.AddCommand(whoamiCommand())
}

func whoamiCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "whoami",
		Short: "config information",
		Run: func(cmd *cobra.Command, args []string) {

			client, err := authbase.NewClient("4000")
			if err != nil {
				panic(err)
			}
			defer client.Close()

			res, err := client.GetAccessKeyAccount(tokenContext(), &v1.GetAccessKeyAccountRequest{})
			if err != nil {
				logrus.Errorf("failed to get account from token: %s", err)
				return
			}

			// print claims
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Project ID", "User ID", "Email", "Username", "Member"})
			table.Append([]string{res.Account.ProjectId, res.Account.Id, res.Account.Email, res.Account.VisibleName, strconv.FormatBool(res.Account.Member)})
			table.Render()
		},
	}

	return command
}
