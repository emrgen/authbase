package cmd

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc/metadata"
)

var contextCommand = &cobra.Command{
	Use:   "context",
	Short: "context commands",
}

func init() {
	contextCommand.AddCommand(setContextCommand())
	contextCommand.AddCommand(currentContextCommand())
	contextCommand.AddCommand(resetContextCommand())
}

type Context struct {
	OrganizationId string `json:"organization_id"`
	Token          string `json:"token"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	ExpireAt       int64  `json:"expire_at"`
}

func loadToken() {
	context := readContext()
	Token = context.Token
}

func readContext() Context {
	var context Context
	viper.SetConfigName("authbase")
	viper.AddConfigPath("./.tmp")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("error reading config file: ", err)
	}

	if err := viper.UnmarshalKey("context", &context); err != nil {
		fmt.Println("error unmarshalling config file: ", err)
	}

	return context
}

func writeContext(context Context) {
	viper.SetConfigName("authbase")
	viper.AddConfigPath("./.tmp")
	viper.SetConfigType("yml")
	viper.Set("context", context)

	if err := viper.WriteConfig(); err != nil {
		fmt.Println("error writing config file: ", err)
	}
}

// saves the context info to the config file in ~/.config/authbase
func setContextCommand() *cobra.Command {
	var token string
	var organization string
	command := &cobra.Command{
		Use:   "set",
		Short: "set context",
		Run: func(cmd *cobra.Command, args []string) {
			if token == "" || organization == "" {
				logrus.Infof(`missing required flags: --token, --organization`)
				return
			}

			// save the context info to the config file
			viper.SetConfigName("authbase")
			viper.AddConfigPath("./.tmp")
			viper.SetConfigType("yml")
			viper.Set("context", Context{
				OrganizationId: organization,
				Token:          token,
			})

			if err := viper.WriteConfig(); err != nil {
				fmt.Println("error writing config file: ", err)
			} else {
				fmt.Println("context saved")
			}
		},
	}

	command.Flags().StringVarP(&token, "token", "t", "", "token")
	command.Flags().StringVarP(&organization, "organization", "o", "", "organization")

	return command
}

func currentContextCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "current",
		Short: "current context",
	}

	return command
}

func resetContextCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "reset",
		Short: "reset context",
	}

	return command
}

func verifyContext() {
	if Token == "" {
		logrus.Error("missing required flags: --token")
		return
	}

	if OrganizationId == "" {
		logrus.Error("missing required flags: --organization")
		return
	}
}

func bindContextFlags(command *cobra.Command) {
	command.Flags().StringVarP(&Token, "token", "t", "", "token")
	command.Flags().StringVarP(&OrganizationId, "organization", "o", "", "organization")
}

func tokenContext(token string) context.Context {
	cfg := readContext()
	if token == "" {
		Token = token
	} else {
		Token = cfg.Token
	}

	md := metadata.New(map[string]string{"Authorization": "Bearer " + Token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	return ctx
}
