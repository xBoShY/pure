package cmd

import (
	"context"
	"crypto"
	"crypto/ed25519"
	"encoding/base64"
	"time"

	"github.com/kennygrant/sanitize"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/spf13/cobra"
	"github.com/xboshy/pure/internal/did"
	"github.com/xboshy/pure/internal/vault"
	"go.bryk.io/pkg/errors"
	xlog "go.bryk.io/pkg/log"
)

var authCmd = &cobra.Command{
	Use:     "authorization",
	Aliases: []string{"auth"},
	Short:   "Generates the Authorization header for an HTTP request",
	Example: "pure req [alias] [secret]",
	RunE: func(_ *cobra.Command, args []string) error {
		var err error

		// Get parameters
		if len(args) < 1 {
			return errors.New("you must provide an alias for your did")
		}
		name := sanitize.Name(args[0])

		secretStr := ""
		if len(args) >= 2 {
			secretStr = args[1]
		}

		if secretStr == "" {
			secretStr, err = generateSecret()
			if err != nil {
				return err
			}
		}

		secret, err := deserializeSecret(secretStr)
		if err != nil {
			return err
		}

		jwkey, err := jwk.Import(secret)
		if err != nil {
			return err
		}

		jkt, err := jwkey.Thumbprint(crypto.SHA256)
		if err != nil {
			return err
		}
		jktStr := base64.RawURLEncoding.EncodeToString([]byte(jkt))

		client, err := vault.NewClient(cfgVault, algorithm)
		if err != nil {
			return err
		}

		ctx := context.Background()

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

		iss := d.String()
		iat := time.Now().Unix()
		exp := iat + 60
		cnf := map[string]any{
			"jkt": jktStr,
		}

		accessClaims := map[string]any{
			"iss": iss,
			"sub": iss,
			"aud": "http://api.pure.com",
			"iat": iat,
			"exp": exp,
			"cnf": cnf,
		}

		access, err := sign(secret, nil, accessClaims)
		if err != nil {
			return err
		}

		log.Infof("Authorization: DPoP %s", string(access))

		log.WithFields(xlog.Fields{
			"name":   name,
			"did":    iss,
			"secret": secretStr,
		}).Info("Authorization header generated")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}
