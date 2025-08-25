package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/xboshy/pure/internal/did"
	"github.com/xboshy/pure/internal/vault"
	"go.bryk.io/pkg/errors"
	xLog "go.bryk.io/pkg/log"
)

var authCmd = &cobra.Command{
	Use:     "authorization",
	Aliases: []string{"auth"},
	Short:   "Generates the Authorization header for an HTTP request",
	Example: "pure http auth [did] [public]",
	RunE: func(_ *cobra.Command, args []string) error {
		// Get parameters
		if len(args) < 1 {
			return errors.New("you must provide a did and public key")
		}

		publicStr := ""
		if len(args) < 2 {
			return errors.New("you must provide a public key")
		}
		publicStr = args[1]
		public, err := deserializePublic(publicStr)
		if err != nil {
			return err
		}

		jkt, err := jktFromPublicKey(public)
		if err != nil {
			return err
		}
		jktStr := base64.RawURLEncoding.EncodeToString([]byte(jkt))

		d, err := did.DidDecode(args[0])
		if err != nil {
			return err
		}

		ctx := context.Background()
		client, err := vault.NewClient(cfgVault, algorithm)
		if err != nil {
			return err
		}

		keys, err := didList(ctx, client)
		if err != nil {
			return err
		}

		key := keys[d]

		iss := d.String()
		iat := time.Now().Unix()
		exp := iat + 60
		cnf := map[string]any{
			"jkt": jktStr,
		}

		header := map[string]any{
			"typ": "JWT",
			"alg": "EdDSA",
		}

		claims := map[string]any{
			"iss": iss,
			"sub": iss,
			"aud": "http://api.pure.com",
			"iat": iat,
			"exp": exp,
			"cnf": cnf,
		}

		log.WithFields(xLog.Fields{
			"header": header,
		}).Info("header")

		log.WithFields(xLog.Fields{
			"claims": claims,
		}).Info("claims")

		headerBytes, err := json.Marshal(header)
		if err != nil {
			return err
		}

		claimsBytes, err := json.Marshal(claims)
		if err != nil {
			return err
		}

		headerB64 := base64.RawURLEncoding.EncodeToString(headerBytes)
		claimsB64 := base64.RawURLEncoding.EncodeToString(claimsBytes)
		tokenB64 := fmt.Sprintf("%s.%s", headerB64, claimsB64)

		signature, err := client.SignData(ctx, []byte(tokenB64), key)
		if err != nil {
			return err
		}

		signatureB64 := base64.RawURLEncoding.EncodeToString(signature)

		log.Infof("Authorization: DPoP %s.%s", tokenB64, signatureB64)

		log.WithFields(xLog.Fields{
			"key": key,
			"did": iss,
		}).Info("Authorization header generated")
		return nil
	},
}

func init() {
	httpCmd.AddCommand(authCmd)
}
