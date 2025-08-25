package cmd

import (
	"context"
	"crypto/ed25519"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/xboshy/pure/internal/did"
	"github.com/xboshy/pure/internal/vault"
	xLog "go.bryk.io/pkg/log"
)

var didCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"gen", "new"},
	Short:   "Create a new DID",
	Example: "pure did new",
	RunE: func(_ *cobra.Command, args []string) error {
		uuid, err := uuid.NewV6()
		if err != nil {
			return err
		}

		name := uuid.String()

		client, err := vault.NewClient(cfgVault, algorithm)
		if err != nil {
			return err
		}

		ctx := context.Background()
		cfg := vault.KeyConfig{
			Algo: algorithm,
			Name: name,
		}

		err = client.NewKey(ctx, cfg)
		if err != nil {
			return err
		}

		key := vault.Key{
			Name:    name,
			Version: 1,
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

		log.WithFields(xLog.Fields{
			"key": key,
			"did": d.String(),
		}).Info("new wallet created")
		return nil
	},
}

func init() {
	didCmd.AddCommand(didCreateCmd)
}
