package cmd

import (
	"context"
	"github.com/emrgen/authbase"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var userCommand = &cobra.Command{
	Use:   "user",
	Short: "user commands",
}

func init() {
	userCommand.AddCommand(createUserCommand())
	userCommand.AddCommand(checkEmailUsedCommand())
	userCommand.AddCommand(registerUserCommand())
	userCommand.AddCommand(loginUserCommand())
	userCommand.AddCommand(logoutUserCommand())
	userCommand.AddCommand(revokeUserSessionsCommand())
	userCommand.AddCommand(listUserCommand())
	userCommand.AddCommand(updateUserCommand())
	userCommand.AddCommand(deleteUserCommand())
	userCommand.AddCommand(enableUserCommand())
	userCommand.AddCommand(disableUserCommand())
}

func createUserCommand() *cobra.Command {
	var username string
	var password string
	var email string

	command := &cobra.Command{
		Use:   "create",
		Short: "create user",
		Run: func(cmd *cobra.Command, args []string) {
			loadToken()

			if Token == "" {
				logrus.Errorf("missing required flags: --token")
				return
			}

			if ProjectId == "" {
				logrus.Errorf("missing required flag: --project-id")
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
			defer client.Close()

			ctx := tokenContext()
			user, err := client.CreateUser(ctx, &v1.CreateUserRequest{
				ProjectId: ProjectId,
				Email:     email,
				Username:  username,
				Password:  password,
			})
			if err != nil {
				logrus.Errorf("failed to create user: %v", err)
				return
			}

			logrus.Info("user created successfully")

			// print response in table
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"#", "ID", "Email", "Username", "CreatedAt"})
			table.Append([]string{
				"1", user.Id, user.Email, user.Username, user.CreatedAt.AsTime().Format("2006-01-02 15:04:05"),
			})
			table.Render()
		},
	}

	bindContextFlags(command)
	command.Flags().StringVarP(&username, "username", "u", "", "username")
	command.Flags().StringVarP(&email, "email", "e", "", "email")
	command.Flags().StringVarP(&password, "password", "p", "", "password")

	return command
}

func checkEmailUsedCommand() *cobra.Command {
	var email string
	var username string
	command := &cobra.Command{
		Use:   "check",
		Short: "check email and username availability",
		Run: func(cmd *cobra.Command, args []string) {

			if email == "" {
				logrus.Errorf("missing required flag: --email")
				return
			}

			if username == "" {
				logrus.Errorf("missing required flag: --username")
				return
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}
			defer client.Close()

			res, err := client.UserEmailExists(tokenContext(), &v1.UserEmailExistsRequest{
				ProjectId: ProjectId,
				Email:     email,
				Username:  username,
			})
			if err != nil {
				logrus.Errorf("failed to check email: %v", err)
				return
			}

			if res.UsernameExists {
				logrus.Warn("username already exists within the project")
			} else {
				logrus.Infof("username is available")
			}

			if res.EmailExists {
				logrus.Warn("email already exists within the project")
			} else {
				logrus.Infof("email is available")
			}
		},
	}

	bindContextFlags(command)
	command.Flags().StringVarP(&email, "email", "e", "", "email")
	command.Flags().StringVarP(&username, "username", "u", "", "username")

	return command

}

func listUserCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "list",
		Short: "list user",
		Run: func(cmd *cobra.Command, args []string) {
			loadToken()
			verifyToken()

			if ProjectId == "" {
				logrus.Errorf("missing required flag: --project-id")
				return
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}
			defer client.Close()

			ctx := tokenContext()
			res, err := client.ListUsers(ctx, &v1.ListUsersRequest{
				ProjectId: ProjectId,
			})
			if err != nil {
				logrus.Errorf("failed to list users: %v", err)
				return
			}

			// print response in table
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"#", "ID", "Email", "Username", "CreatedAt", "Verified", "Active", "Member"})
			for i, user := range res.Users {
				verified := user.VerifiedAt.AsTime().Format("2006-01-02 15:04:05") != "1970-01-01 00:00:00"
				table.Append([]string{
					strconv.Itoa(i + 1),
					user.Id,
					user.Email,
					user.Username,
					user.CreatedAt.AsTime().Format("2006-01-02 15:04:05"),
					strconv.FormatBool(verified),
					strconv.FormatBool(!user.Disabled),
					strconv.FormatBool(user.Member),
				})
			}

			table.Render()
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
	var username string
	var email string
	var password string

	command := &cobra.Command{
		Use:   "register",
		Short: "register user",
		Run: func(cmd *cobra.Command, args []string) {
			if ProjectId == "" {
				logrus.Errorf("missing required flag: --project-id")
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
				ProjectId: ProjectId,
				Username:  username,
				Email:     email,
				Password:  password,
			})
			if err != nil {
				logrus.Errorf("failed to register user: %v", err)
				return
			}

			logrus.Infof("register user: %v", res)
		},
	}

	bindContextFlags(command)

	command.Flags().StringVarP(&username, "username", "u", "", "username")
	command.Flags().StringVarP(&email, "email", "e", "", "email")
	command.Flags().StringVarP(&password, "password", "p", "", "password")

	return command

}

func loginUserCommand() *cobra.Command {
	var email string
	var password string

	command := &cobra.Command{
		Use:   "login",
		Short: "login user",
		Run: func(cmd *cobra.Command, args []string) {
			if ProjectId == "" {
				logrus.Errorf("missing required flag: --project-id")
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
			defer client.Close()

			res, err := client.Login(context.Background(), &v1.LoginRequest{
				Email:     email,
				Password:  password,
				ProjectId: ProjectId,
			})
			if err != nil {
				logrus.Errorf("failed to login user: %v", err)
				return
			}

			// print response in table
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"#", "ID", "Email", "Username", "CreatedAt", "UpdatedAt", "Token"})
			table.Append([]string{
				"1", res.User.Id, res.User.Email, res.User.Username,
				res.User.CreatedAt.AsTime().Format("2006-01-02 15:04:05"),
				res.User.UpdatedAt.AsTime().Format("2006-01-02 15:04:05"),
				res.Token.AccessToken,
			})

			table.Render()
		},
	}

	bindContextFlags(command)

	command.Flags().StringVarP(&email, "email", "e", "", "email")
	command.Flags().StringVarP(&password, "password", "w", "", "password")

	return command
}

func logoutUserCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "logout",
		Short: "logout user",
		Run: func(cmd *cobra.Command, args []string) {
			if Token == "" {
				logrus.Errorf("missing required flag: --token")
				return
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}
			defer client.Close()

			_, err = client.Logout(tokenContext(), &v1.LogoutRequest{})
			if err != nil {
				logrus.Errorf("failed to logout user: %v", err)
				return
			}

			logrus.Infof("user logged out successfully")
		},
	}

	bindContextFlags(command)

	return command

}

func revokeUserSessionsCommand() *cobra.Command {
	var userID string

	command := &cobra.Command{
		Use:   "revoke",
		Short: "revoke user sessions",
		Run: func(cmd *cobra.Command, args []string) {
			loadToken()

			if Token == "" {
				logrus.Errorf("missing required flags: --token")
				return
			}

			if userID == "" {
				logrus.Errorf("missing required flag: --user-id")
				return
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}
			defer client.Close()

			_, err = client.DeleteAllSessions(tokenContext(), &v1.DeleteAllSessionsRequest{
				UserId: userID,
			})
			if err != nil {
				logrus.Errorf("failed to revoke user sessions: %v", err)
				return
			}

			logrus.Infof("user sessions revoked successfully")
		},
	}

	bindContextFlags(command)

	command.Flags().StringVarP(&userID, "user-id", "u", "", "user id")

	return command
}

func enableUserCommand() *cobra.Command {
	var userID string

	command := &cobra.Command{
		Use:   "enable",
		Short: "enable user",
		Run: func(cmd *cobra.Command, args []string) {
			loadToken()

			if Token == "" {
				logrus.Errorf("missing required flags: --token")
				return
			}

			if userID == "" {
				logrus.Errorf("missing required flag: --user-id")
				return
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}

			ctx := tokenContext()
			_, err = client.EnableUser(ctx, &v1.EnableUserRequest{
				UserId:    userID,
				ProjectId: ProjectId,
			})

			logrus.Infof("user enabled successfully")
		},
	}

	bindContextFlags(command)

	command.Flags().StringVarP(&userID, "user-id", "u", "", "user id")

	return command
}

func disableUserCommand() *cobra.Command {
	var userID string

	command := &cobra.Command{
		Use:   "disable",
		Short: "disable user",
		Run: func(cmd *cobra.Command, args []string) {
			loadToken()

			if Token == "" {
				logrus.Errorf("missing required flags: --token")
				return
			}

			if userID == "" {
				logrus.Errorf("missing required flag: --user-id")
				return
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}

			ctx := tokenContext()
			_, err = client.DisableUser(ctx, &v1.DisableUserRequest{
				UserId: userID,
			})

			logrus.Infof("user disabled successfully")
		},
	}

	bindContextFlags(command)

	command.Flags().StringVarP(&userID, "user-id", "u", "", "user id")

	return command
}
