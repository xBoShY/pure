package cmd

import (
	"github.com/spf13/cobra"
	"go.bryk.io/pkg/errors"
)

var didRotateCmd = &cobra.Command{
	Use:     "rotate",
	Aliases: []string{},
	Short:   "Rotate",
	Example: "pure did rotate [did]",
	RunE: func(_ *cobra.Command, args []string) error {
		return errors.New("not implemented")
	},
}

func init() {
	didCmd.AddCommand(didRotateCmd)
}
