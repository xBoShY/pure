package cmd

import (
	"github.com/spf13/cobra"
	xLog "go.bryk.io/pkg/log"
)

var secretGenCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen", "new"},
	Short:   "Generates a new secret key",
	Example: "pure secret new",
	RunE: func(_ *cobra.Command, args []string) error {
		var err error

		pubStr, secretStr, err := generateSecret()
		if err != nil {
			return err
		}

		log.WithFields(xLog.Fields{
			"public": pubStr,
			"secret": secretStr,
		}).Info("secret generated")
		return nil
	},
}

func init() {
	secretCmd.AddCommand(secretGenCmd)
}
