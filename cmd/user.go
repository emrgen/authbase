package cmd

import (
	"context"
	"github.com/emrgen/authbase"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var userCommand = &cobra.Command{
	Use:   "user",
	Short: "user commands",
}

func init() {
	userCommand.AddCommand(createUserCommand())
	userCommand.AddCommand(listUserCommand())
	userCommand.AddCommand(updateUserCommand())
	userCommand.AddCommand(deleteUserCommand())
}

func createUserCommand() *cobra.Command {
	var username string
	var email string

	command := &cobra.Command{
		Use:   "create",
		Short: "create user",
		Run: func(cmd *cobra.Command, args []string) {
			verifyContext()
			if username == "" {
				logrus.Errorf("missing required flag: --username")
				return
			}

			if email == "" {
				logrus.Errorf("missing required flag: --email")
				return
			}

			_, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}

			//client.CreateUser()

			logrus.Infof("user created successfully %v", "")
		},
	}

	bindContextFlags(command)
	command.Flags().StringVarP(&username, "username", "u", "", "username")
	command.Flags().StringVarP(&email, "email", "e", "", "email")

	return command
}

func listUserCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "list",
		Short: "list user",
		Run: func(cmd *cobra.Command, args []string) {
			verifyContext()

			_, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
			}
		},
	}

	bindContextFlags(command)

	return command
}

func updateUserCommand() *cobra.Command {
	var username string
	var email string

	command := &cobra.Command{
		Use:   "update",
		Short: "update user",
		Run: func(cmd *cobra.Command, args []string) {
			verifyContext()
			if username == "" {
				logrus.Errorf("missing required flag: --username")
				return
			}

			if email == "" {
				logrus.Errorf("missing required flag: --email")
				return
			}

			_, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}

			logrus.Infof("update user")
		},
	}

	bindContextFlags(command)
	command.Flags().StringVarP(&username, "username", "u", "", "username")
	command.Flags().StringVarP(&email, "email", "e", "", "email")

	return command
}

func deleteUserCommand() *cobra.Command {
	var username string

	command := &cobra.Command{
		Use:   "delete",
		Short: "delete user",
		Run: func(cmd *cobra.Command, args []string) {
			verifyContext()

			if username == "" {
				logrus.Errorf("missing required flag: --username")
				return
			}

			_, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}

			logrus.Infof("delete user")
		},
	}

	bindContextFlags(command)

	command.Flags().StringVarP(&username, "username", "u", "", "username")

	return command
}

func registerUserCommand() *cobra.Command {
	var organizationId string
	var username string
	var email string
	var password string

	command := &cobra.Command{
		Use:   "register",
		Short: "register user",
		Run: func(cmd *cobra.Command, args []string) {
			if organizationId == "" {
				logrus.Errorf("missing required flag: --organization-id")
				return
			}

			if username == "" {
				logrus.Errorf("missing required flag: --username")
				return
			}

			if email == "" {
				logrus.Errorf("missing required flag: --email")
				return
			}

			if password == "" {
				logrus.Errorf("missing required flag: --password")
				return
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}

			res, err := client.Register(context.Background(), &v1.RegisterRequest{
				OrganizationId: organizationId,
				Username:       username,
				Email:          email,
				Password:       password,
			})
			if err != nil {
				logrus.Errorf("failed to register user: %v", err)
				return
			}

			logrus.Infof("register user: %v", res)
		},
	}

	bindContextFlags(command)

	command.Flags().StringVarP(&organizationId, "organization-id", "o", "", "organization id")
	command.Flags().StringVarP(&username, "username", "u", "", "username")
	command.Flags().StringVarP(&email, "email", "e", "", "email")
	command.Flags().StringVarP(&password, "password", "p", "", "password")

	return command

}
