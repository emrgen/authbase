package cmd

import (
	"github.com/emrgen/authbase"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var idpCommand = &cobra.Command{
	Use:   "idp",
	Short: "idp commands",
}

func init() {
	idpCommand.AddCommand(idpCreateCommand())
	idpCommand.AddCommand(idpListCommand())
	idpCommand.AddCommand(idpUpdateCommand())
	idpCommand.AddCommand(idpDeleteCommand())
}

func idpCreateCommand() *cobra.Command {
	var poolID string
	var name string
	var clientID string
	var clientSecret string

	command := &cobra.Command{
		Use:   "create",
		Short: "Create an IDP",
		Run: func(cmd *cobra.Command, args []string) {

			if name == "" {
				logrus.Infof("missing required flag: --name")
				return
			}

			if clientID == "" {
				logrus.Infof("missing required flag: --client-id")
				return
			}

			if clientSecret == "" {
				logrus.Infof("missing required flag: --client-secret")
				return
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}

			if poolID == "" {
				account := getAccount(client)
				poolID = account.PoolId
			}

			_, err = client.AddOauthProvider(tokenContext(), &v1.AddOauthProviderRequest{
				PoolId: poolID,
				Provider: &v1.OAuthProvider{
					Provider:     v1.Idp_IDP_GOOGLE,
					ClientId:     clientID,
					ClientSecret: clientSecret,
					RedirectUris: nil,
				},
			})
			if err != nil {
				logrus.Errorf("failed to create idp: %v", err)
				return
			}

			logrus.Infof("idp created successfully %v", name)
		},
	}

	command.Flags().StringVarP(&name, "name", "n", "", "name")

	return command
}

func idpListCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "list",
		Short: "List IDPs",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}

			res, err := client.ListOauthProviders(tokenContext(), &v1.ListOauthProvidersRequest{})
			if err != nil {
				logrus.Errorf("failed to list idps: %v", err)
				return
			}

			table := tablewriter.NewWriter(cmd.OutOrStdout())
			table.SetHeader([]string{"ID", "Provider", "Pool ID"})
			for _, idp := range res.Providers {
				table.Append([]string{idp.Id, idp.Provider.String(), idp.PoolId})
			}
			table.Render()
		},
	}

	return command
}

func idpUpdateCommand() *cobra.Command {
	var provider string
	var clientID string
	var clientSecret string

	command := &cobra.Command{
		Use:   "update",
		Short: "Update an IDP",
		Run: func(cmd *cobra.Command, args []string) {

			if provider == "" {
				logrus.Infof("missing required flag: --provider")
				return
			}

			if clientID == "" {
				logrus.Infof("missing required flag: --client-id")
				return
			}

			if clientSecret == "" {
				logrus.Infof("missing required flag: --client-secret")
				return
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}

			_, err = client.UpdateOauthProvider(tokenContext(), &v1.UpdateOauthProviderRequest{
				Provider: &v1.OAuthProvider{
					Provider:     v1.Idp_IDP_GOOGLE,
					ClientId:     clientID,
					ClientSecret: clientSecret,
				},
			})
			if err != nil {
				logrus.Errorf("failed to update idp: %v", err)
				return
			}

			logrus.Infof("idp updated successfully %v", provider)
		},
	}

	command.Flags().StringVarP(&provider, "provider", "p", "", "idp provider name")

	return command
}

func idpDeleteCommand() *cobra.Command {
	var clientID string
	var provider string

	command := &cobra.Command{
		Use:   "delete",
		Short: "Delete an IDP",
		Run: func(cmd *cobra.Command, args []string) {
			if provider == "" {
				logrus.Infof("missing required flag: --idp-provider-name")
				return
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}

			_, err = client.DeleteOauthProvider(tokenContext(), &v1.DeleteOauthProviderRequest{
				Provider: v1.Idp_IDP_GOOGLE,
			})
			if err != nil {
				logrus.Errorf("failed to delete idp: %v", err)
				return
			}

			logrus.Infof("idp deleted successfully %v", provider)
		},
	}

	command.Flags().StringVarP(&provider, "provider", "p", "", "idp provider name")
	command.Flags().StringVarP(&clientID, "client-id", "c", "", "client id")

	return command
}
