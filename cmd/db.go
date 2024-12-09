package cmd

import (
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
			//db := config.GetDb(config.FromEnv())
			//err := model.Migrate(db)
			//if err != nil {
			//	panic(err)
			//}
		},
	}

	return command
}
