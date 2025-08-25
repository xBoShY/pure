package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/xboshy/pure/internal/vault"
	xLog "go.bryk.io/pkg/log"
)

var didListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List the current vault keys",
	Example: "pure did ls",
	RunE: func(_ *cobra.Command, args []string) error {
		ctx := context.Background()
		client, err := vault.NewClient(cfgVault, algorithm)
		if err != nil {
			return err
		}

		didKeys, err := didList(ctx, client)
		if err != nil {
			return err
		}

		for d, key := range didKeys {
			log.WithFields(xLog.Fields{
				"key": key,
				"did": d.String(),
			}).Info("")
		}

		log.Info("all keys retrieved")
		return nil
	},
}

func init() {
	didCmd.AddCommand(didListCmd)
}
