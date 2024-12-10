package cmd

import (
	"github.com/emrgen/authbase/x"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var keygenCommand = &cobra.Command{
	Use:   "keygen",
	Short: "generate a new key",
	Run: func(cmd *cobra.Command, args []string) {
		key := x.Keygen()

		logrus.Info(key)
	},
}
