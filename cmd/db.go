package cmd

import (
	"github.com/emrgen/authbase/pkg/config"
	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "db commands",
}

func init() {
	dbCmd.AddCommand(Migrate())
}

func Migrate() *cobra.Command {
	command := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate the database",
		Run: func(cmd *cobra.Command, args []string) {
			db := config.GetDB()
			err := db.Migrate()
			if err != nil {
				panic(err)
			}
		},
	}

	return command
}
