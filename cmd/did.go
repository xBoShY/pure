package cmd

import (
	"context"
	"crypto/ed25519"

	"github.com/spf13/cobra"
	"github.com/xboshy/pure/internal/did"
	"github.com/xboshy/pure/internal/vault"
)

var didCmd = &cobra.Command{
	Use:     "did",
	Aliases: []string{},
	Short:   "DID tools",
}

func init() {
	rootCmd.AddCommand(didCmd)
}

func didList(ctx context.Context, client *vault.Client) (map[did.DID]vault.Key, error) {
	keys, err := client.ListKeys(ctx)
	if err != nil {
		return nil, err
	}

	result := map[did.DID]vault.Key{}

	for _, key := range keys {
		pub, err := client.GetPublicKey(ctx, key)
		if err != nil {
			return nil, err
		}

		obj, err := did.NewObject(ed25519.PublicKey(pub))
		if err != nil {
			return nil, err
		}

		d, err := obj.DID("mainnet")

		result[d] = key
	}

	return result, nil
}
