package cmd

import (
	"github.com/spf13/cobra"
	"go.bryk.io/pkg/errors"
)

var didRekeyCmd = &cobra.Command{
	Use:     "rekey",
	Aliases: []string{},
	Short:   "Rekey",
	Example: "pure did rekey [did] [new-did-controller]",
	RunE: func(_ *cobra.Command, args []string) error {
		return errors.New("not implemented")
	},
}

func init() {
	didCmd.AddCommand(didRekeyCmd)
}
