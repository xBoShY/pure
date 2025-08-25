package cmd

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/xboshy/pure/internal/did"

	"github.com/spf13/cobra"
	"go.bryk.io/pkg/errors"
)

var verifyCmd = &cobra.Command{
	Use:     "verify",
	Aliases: []string{"check"},
	Short:   "Verifies a JWT token",
	Example: "pure http verify [jwt]",
	RunE: func(_ *cobra.Command, args []string) error {
		// Get parameters
		if len(args) != 1 {
			return errors.New("you must provide a secret")
		}

		partsB64 := strings.Split(args[0], ".")
		if len(partsB64) != 3 {
			return errors.New("invalid jwt")
		}

		signingB64 := partsB64[0] + "." + partsB64[1]
		signature, err := base64.RawURLEncoding.DecodeString(partsB64[2])
		if err != nil {
			return err
		}

		headerBytes, err := base64.RawURLEncoding.DecodeString(partsB64[0])
		if err != nil {
			return err
		}

		header := make(map[string]any)
		err = json.Unmarshal(headerBytes, &header)
		if err != nil {
			return err
		}

		alg, ok := header["alg"]
		if !ok {
			return errors.New("invalid header")
		}
		if strings.ToUpper(alg.(string)) != "EDDSA" {
			return errors.New("unsupported algorithm")
		}

		typ, ok := header["typ"]
		if !ok {
			return errors.New("invalid header")
		}

		var publicKey ed25519.PublicKey
		switch strings.ToUpper(typ.(string)) {
		case "JWT":
			claimsBytes, err := base64.RawURLEncoding.DecodeString(partsB64[1])
			if err != nil {
				return err
			}

			claims := make(map[string]any)
			err = json.Unmarshal(claimsBytes, &claims)
			if err != nil {
				return err
			}

			issBytes, ok := claims["iss"]
			if !ok {
				return errors.New("missing iss in jwt claims")
			}

			iss := issBytes.(string)
			d, err := did.DidDecode(iss)
			if err != nil {
				return err
			}

			publicKey, err = did.Object(d.Id).CanonicalPublicKey()
			if err != nil {
				return err
			}

		case "DPOP+JWT":
			jwk, ok := header["jwk"]
			if !ok {
				return errors.New("missing jwk in jwt header")
			}

			jwkObj := jwk.(map[string]any)
			pub, ok := jwkObj["x"]
			if !ok {
				return errors.New("missing d in jwk")
			}

			publicKey, err = deserializePublic(pub.(string))
			if err != nil {
				return err
			}
		default:
			return errors.New("unsupported jwt type")
		}

		meth := jwt.GetSigningMethod("EdDSA")
		err = meth.Verify(signingB64, signature, publicKey)
		if err != nil {
			return err
		}

		log.Info("JWT verified with success!")

		return nil
	},
}

func init() {
	httpCmd.AddCommand(verifyCmd)
}
