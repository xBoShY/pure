package cmd

import (
	"github.com/spf13/cobra"
	xlog "go.bryk.io/pkg/log"
)

var secretCmd = &cobra.Command{
	Use:     "secret",
	Aliases: []string{},
	Short:   "Generates a new secret key",
	Example: "pure secret",
	RunE: func(_ *cobra.Command, args []string) error {
		var err error

		secretStr, err := generateSecret()
		if err != nil {
			return err
		}

		log.WithFields(xlog.Fields{
			"secret": secretStr,
		}).Info("Secret generated")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(secretCmd)
}
