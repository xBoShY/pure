package cmd

import (
	"context"
	"crypto/ed25519"

	"github.com/kennygrant/sanitize"
	"github.com/spf13/cobra"
	"github.com/xboshy/pure/internal/did"
	"github.com/xboshy/pure/internal/vault"
	"go.bryk.io/pkg/errors"
	xlog "go.bryk.io/pkg/log"
)

var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new"},
	Short:   "Create a new DID",
	Example: "pure new [alias]",
	RunE: func(_ *cobra.Command, args []string) error {
		// Get parameters
		if len(args) != 1 {
			return errors.New("you must provide an alias for your did")
		}
		name := sanitize.Name(args[0])

		client, err := vault.NewClient(cfgVault, algorithm)
		if err != nil {
			return err
		}

		ctx := context.Background()
		cfg := vault.KeyConfig{
			KeyAlgo: algorithm,
			KeyName: name,
		}

		err = client.NewKey(ctx, cfg)
		if err != nil {
			return err
		}

		key := vault.Key{
			KeyName:    name,
			KeyVersion: 1,
		}

		pub, err := client.GetPublicKey(ctx, key)
		if err != nil {
			return err
		}

		obj, err := did.NewObject(ed25519.PublicKey(pub))
		if err != nil {
			return err
		}

		d, err := obj.DID("mainnet")

		log.WithFields(xlog.Fields{
			"name": name,
			"did":  d.String(),
		}).Info("new wallet created")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
