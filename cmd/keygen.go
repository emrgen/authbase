package cmd

import (
	"fmt"
	"github.com/emrgen/authbase/x"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var keygenCommand = &cobra.Command{
	Use:   "keygen",
	Short: "generate a random key",
	Run: func(cmd *cobra.Command, args []string) {
		key := x.Keygen()

		logrus.Info(key)
	},
}

func init() {
	keygenCommand.AddCommand(shortCmd())
}

func shortCmd() *cobra.Command {
	var size int
	cmd := &cobra.Command{
		Use:   "short",
		Short: "generate a short key",
		Run: func(cmd *cobra.Command, args []string) {
			key := x.KeygenSize(size)

			fmt.Printf("%s\n", key)
		},
	}

	cmd.Flags().IntVarP(&size, "size", "s", 32, "size of the key")

	return cmd
}
