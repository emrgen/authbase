package cmd

import (
	"context"
	"fmt"
	goset "github.com/deckarep/golang-set/v2"
	"github.com/emrgen/authbase"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
)

var userCommand = &cobra.Command{
	Use:   "account",
	Short: "account commands",
}

func init() {
	userCommand.AddCommand(createUserCommand())
	userCommand.AddCommand(getUserCommand())
	userCommand.AddCommand(checkEmailUsedCommand())
	userCommand.AddCommand(registerUserCommand())
	userCommand.AddCommand(loginUserCommand())
	userCommand.AddCommand(logoutUserCommand())
	userCommand.AddCommand(forgotPasswordCommand())
	userCommand.AddCommand(resetPasswordCommand())
	userCommand.AddCommand(changePasswordCommand())
	userCommand.AddCommand(revokeUserSessionsCommand())
	userCommand.AddCommand(listUserCommand())
	userCommand.AddCommand(updateUserCommand())
	userCommand.AddCommand(deleteUserCommand())
	userCommand.AddCommand(enableUserCommand())
	userCommand.AddCommand(disableUserCommand())
	userCommand.AddCommand(listUserSessionsCommand())
	userCommand.AddCommand(listActiveSessionsCommand())

	userCommand.AddCommand(getUserPermissionsCommand())
}

func createUserCommand() *cobra.Command {
	var poolID string
	var username string
	var password string
	var email string

	command := &cobra.Command{
		Use:   "create",
		Short: "create account",
		Run: func(cmd *cobra.Command, args []string) {
			if poolID == "" {
				logrus.Errorf("missing required flag: --client-id")
				return
			}

			if username == "" {
				username = strings.Split(email, "@")[0]
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
			user, err := client.CreateAccount(ctx, &v1.CreateAccountRequest{
				PoolId:   poolID,
				Email:    email,
				Username: username,
				Password: password,
			})
			if err != nil {
				logrus.Errorf("failed to create user: %v", err)
				return
			}

			logrus.Info("user created successfully")

			// print response in table
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"#", "ID", "Email", "Username", "CreatedAt"})
			account := user.Account
			table.Append([]string{
				"1", account.Id, account.Email, account.Username, account.CreatedAt.AsTime().Format("2006-01-02 15:04:05"),
			})
			table.Render()
		},
	}

	command.Flags().StringVarP(&poolID, "pool-id", "c", "", "pool id")
	command.Flags().StringVarP(&username, "username", "u", "", "username")
	command.Flags().StringVarP(&email, "email", "e", "", "email")
	command.Flags().StringVarP(&password, "password", "p", "", "password")

	return command
}

func getUserCommand() *cobra.Command {
	var accountID string

	command := &cobra.Command{
		Use:   "get",
		Short: "get account",
		Run: func(cmd *cobra.Command, args []string) {
			loadToken()

			if Token == "" {
				logrus.Errorf("missing required flags: --token")
				return
			}

			if accountID == "" {
				logrus.Errorf("missing required flag: --user-id")
				return
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}
			defer client.Close()

			res, err := client.GetAccount(tokenContext(), &v1.GetAccountRequest{
				Id: accountID,
			})
			if err != nil {
				logrus.Errorf("failed to get user: %v", err)
				return
			}

			// print response in table
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"#", "ID", "Email", "Username", "Verified", "Active", "Member"})
			table.Append([]string{
				"1",
				res.Account.Id,
				res.Account.Email,
				res.Account.Username,
				strconv.FormatBool(res.Account.VerifiedAt.AsTime().Format("2006-01-02 15:04:05") != "1970-01-01 00:00:00"),
				strconv.FormatBool(!res.Account.Disabled),
				strconv.FormatBool(res.Account.Member),
			})
			table.Render()
		},
	}

	bindContextFlags(command)
	command.Flags().StringVarP(&accountID, "account-id", "a", "", "account id")

	return command
}

func checkEmailUsedCommand() *cobra.Command {
	var projectID string
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

			res, err := client.AccountEmailExists(tokenContext(), &v1.AccountEmailExistsRequest{
				ProjectId: projectID,
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

	command.Flags().StringVarP(&projectID, "project-id", "r", "", "project id")
	command.Flags().StringVarP(&email, "email", "e", "", "email")
	command.Flags().StringVarP(&username, "username", "u", "", "username")

	return command

}

func listUserCommand() *cobra.Command {
	var projectID string
	var poolID string
	var roleName string

	command := &cobra.Command{
		Use:   "list",
		Short: "list account",
		Run: func(cmd *cobra.Command, args []string) {
			loadToken()
			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}
			defer client.Close()

			if projectID == "" && poolID == "" {
				if Token != "" {
					account := getAccount(client)
					poolID = account.PoolId
				} else {
					logrus.Errorf("missing required flags: --project-id or --pool-id")
					return
				}
			}

			req := &v1.ListAccountsRequest{}

			if projectID != "" {
				req.ProjectId = &projectID
			}

			if poolID != "" {
				req.PoolId = &poolID
			}

			if roleName != "" {
				req.RoleName = &roleName
			}

			res, err := client.ListAccounts(tokenContext(), req)
			if err != nil {
				logrus.Errorf("failed to list users: %v", err)
				return
			}

			// print response in table
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"#", "ID", "Email", "Name", "CreatedAt", "Verified", "Active", "Member"})
			for i, user := range res.Accounts {
				verified := user.VerifiedAt.AsTime().Format("2006-01-02 15:04:05") != "1970-01-01 00:00:00"
				table.Append([]string{
					strconv.Itoa(i + 1),
					user.Id,
					user.Email,
					user.VisibleName,
					user.CreatedAt.AsTime().Format("2006-01-02 15:04:05"),
					strconv.FormatBool(verified),
					strconv.FormatBool(!user.Disabled),
					strconv.FormatBool(user.Member),
				})
			}

			table.Render()

			fmt.Printf("Users: page: %v, showing: %v, total: %v\n", res.Meta.Page, len(res.Accounts), res.Meta.Total)
		},
	}

	bindContextFlags(command)
	command.Flags().StringVarP(&poolID, "pool-id", "p", "", "pool id")
	command.Flags().StringVarP(&projectID, "project-id", "r", "", "project id")
	command.Flags().StringVarP(&roleName, "role", "s", "", "role name")

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
	var projectID string
	var username string
	var email string
	var password string

	command := &cobra.Command{
		Use:   "register",
		Short: "register user",
		Run: func(cmd *cobra.Command, args []string) {
			if projectID == "" {
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

			res, err := client.RegisterUsingPassword(context.Background(), &v1.RegisterUsingPasswordRequest{
				ProjectId: projectID,
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
	var clientID string
	var email string
	var password string

	command := &cobra.Command{
		Use:   "login",
		Short: "login user",
		Run: func(cmd *cobra.Command, args []string) {
			if clientID == "" {
				logrus.Errorf("missing required flag: --client-id")
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

			res, err := client.LoginUsingPassword(context.Background(), &v1.LoginUsingPasswordRequest{
				Email:    email,
				Password: password,
				ClientId: clientID,
			})
			if err != nil {
				logrus.Errorf("failed to login user: %v", err)
				return
			}

			cmd.Printf("Name: %v\n", res.Account.VisibleName)
			cmd.Printf("Email: %v\n", res.Account.Email)
			cmd.Printf("AccessToken: %v\n", res.Token.AccessToken)
			cmd.Printf("RefreshToken: %v\n", res.Token.RefreshToken)
		},
	}

	bindContextFlags(command)
	command.Flags().StringVarP(&clientID, "client-id", "c", "", "client id")
	command.Flags().StringVarP(&email, "email", "e", "", "email")
	command.Flags().StringVarP(&password, "password", "p", "", "password")

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

func forgotPasswordCommand() *cobra.Command {
	var clientID string
	var email string

	command := &cobra.Command{
		Use:   "forgot",
		Short: "forgot password",
		Run: func(cmd *cobra.Command, args []string) {
			if clientID == "" {
				logrus.Errorf("missing required flag: --client-id")
				return
			}

			if email == "" {
				logrus.Errorf("missing required flag: --email")
				return
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}
			defer client.Close()

			_, err = client.ForgotPassword(context.Background(), &v1.ForgotPasswordRequest{
				ClientId: clientID,
				Email:    email,
			})
			if err != nil {
				logrus.Errorf("failed to forgot password: %v", err)
				return
			}

			logrus.Infof("password reset code sent successfully")
		},
	}

	command.Flags().StringVarP(&email, "email", "e", "", "email")

	return command
}

func resetPasswordCommand() *cobra.Command {
	var code string
	var password string

	command := &cobra.Command{
		Use:   "reset",
		Short: "reset password",
		Run: func(cmd *cobra.Command, args []string) {
			if code == "" {
				logrus.Errorf("missing required flag: --code")
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

			_, err = client.ResetPassword(context.Background(), &v1.ResetPasswordRequest{
				Code:        code,
				NewPassword: password,
			})
			if err != nil {
				logrus.Errorf("failed to reset password: %v", err)
				return
			}

			logrus.Infof("password reset successfully")
		},
	}

	command.Flags().StringVarP(&code, "code", "c", "", "reset code")
	command.Flags().StringVarP(&password, "password", "p", "", "new password")

	return command
}

func changePasswordCommand() *cobra.Command {
	var accountID string
	var newPassword string

	command := &cobra.Command{
		Use:   "change",
		Short: "change password",
		Run: func(cmd *cobra.Command, args []string) {
			if accountID == "" {
				logrus.Errorf("missing required flag: --account-id")
				return
			}

			if newPassword == "" {
				logrus.Errorf("missing required flag: --new-password")
				return
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}
			defer client.Close()

			_, err = client.ChangePassword(tokenContext(), &v1.ChangePasswordRequest{
				AccountId:   accountID,
				NewPassword: newPassword,
			})
			if err != nil {
				logrus.Errorf("failed to change password: %v", err)
				return
			}

			logrus.Infof("password changed successfully")
		},
	}

	bindContextFlags(command)

	command.Flags().StringVarP(&accountID, "account-id", "a", "", "account id")
	command.Flags().StringVarP(&newPassword, "new-password", "n", "", "new password")

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
				AccountId: userID,
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

func listUserSessionsCommand() *cobra.Command {
	var userID string

	command := &cobra.Command{
		Use:   "sessions",
		Short: "list active user sessions",
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
		},
	}

	bindContextFlags(command)

	command.Flags().StringVarP(&userID, "user-id", "u", "", "user id")

	return command
}

func listActiveSessionsCommand() *cobra.Command {
	var poolID string

	command := &cobra.Command{
		Use:   "list-active",
		Short: "list active sessions",
		Run: func(cmd *cobra.Command, args []string) {
			loadToken()

			if Token == "" {
				logrus.Errorf("missing required flags: --token")
				return
			}

			if poolID == "" {
				logrus.Errorf("missing required flag: --pool-id")
				return
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}
			defer client.Close()

			res, err := client.ListActiveAccounts(tokenContext(), &v1.ListActiveAccountsRequest{
				PoolId: poolID,
			})
			if err != nil {
				logrus.Errorf("failed to list active sessions: %v", err)
				return
			}

			// print response in table
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"#", "Pool ID", "User ID", "Name", "Last Used At"})
			for i, account := range res.Accounts {
				table.Append([]string{
					strconv.Itoa(i + 1),
					account.PoolId,
					account.Id,
					account.VisibleName,
					account.LastUsedAt.AsTime().Format("2006-01-02 15:04:05"),
				})
			}

			table.Render()

			fmt.Printf("Sessions: page: %v, showing: %v, total: %v\n", res.Meta.Page, len(res.Accounts), res.Meta.Total)
		},
	}

	bindContextFlags(command)

	command.Flags().StringVarP(&poolID, "pool-id", "p", "", "pool id")

	return command
}

func enableUserCommand() *cobra.Command {
	var projectID string
	var userID string

	command := &cobra.Command{
		Use:   "enable",
		Short: "enable account",
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
			_, err = client.EnableAccount(ctx, &v1.EnableAccountRequest{
				AccountId: userID,
				ProjectId: projectID,
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

			_, err = client.DisableAccount(tokenContext(), &v1.DisableAccountRequest{
				AccountId: userID,
			})

			logrus.Infof("user disabled successfully")
		},
	}

	bindContextFlags(command)
	command.Flags().StringVarP(&userID, "user-id", "u", "", "user id")

	return command
}

func getUserPermissionsCommand() *cobra.Command {
	var accountID string

	command := &cobra.Command{
		Use:   "perm",
		Short: "get user permissions",
		Run: func(cmd *cobra.Command, args []string) {
			if accountID == "" {
				logrus.Errorf("missing required flag: --user-id")
				return
			}

			client, err := authbase.NewClient(":4000")
			if err != nil {
				logrus.Errorf("failed to create client: %v", err)
				return
			}

			res, err := client.ListGroups(tokenContext(), &v1.ListGroupsRequest{
				AccountId: &accountID,
			})
			if err != nil {
				logrus.Errorf("failed to get user permissions: %v", err)
				return
			}

			roles := goset.NewSet[string]()
			for _, group := range res.Groups {
				for _, role := range group.Roles {
					roles.Add(role.Name)
				}
			}

			// print response in table
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"#", "Roles"})
			for i, role := range roles.ToSlice() {
				table.Append([]string{
					strconv.Itoa(i + 1),
					role,
				})
			}
			table.Render()
		},
	}

	bindContextFlags(command)
	command.Flags().StringVarP(&accountID, "account-id", "a", "", "account id")

	return command
}
