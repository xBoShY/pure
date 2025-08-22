package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/xboshy/pure/internal/vault"
	xlog "go.bryk.io/pkg/log"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List the current vault keys",
	Example: "pure ls",
	RunE: func(_ *cobra.Command, args []string) error {
		client, err := vault.NewClient(cfgVault, algorithm)
		if err != nil {
			return err
		}

		ctx := context.Background()

		keys, err := client.ListKeys(ctx)
		if err != nil {
			return err
		}

		log.WithFields(xlog.Fields{
			"keys": keys,
		}).Info("keys retrieved")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
