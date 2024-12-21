package cmd

import (
	"github.com/emrgen/authbase/pkg/store"
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
			db := store.GetDB()
			err := db.Migrate()
			if err != nil {
				panic(err)
			}
		},
	}

	return command
}
