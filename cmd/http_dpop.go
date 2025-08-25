package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/spf13/cobra"
	"go.bryk.io/pkg/errors"
	xLog "go.bryk.io/pkg/log"
)

var dpopCmd = &cobra.Command{
	Use:     "dpop",
	Aliases: []string{"req", "call"},
	Short:   "Generates the DPoP header for an HTTP request",
	Example: "pure http dpop [secret]",
	RunE: func(_ *cobra.Command, args []string) error {
		var err error

		uuid, err := uuid.NewV6()
		if err != nil {
			return err
		}

		jti := uuid.String()

		// Get parameters
		if len(args) != 1 {
			return errors.New("you must provide a secret")
		}

		secretStr := args[0]
		secret, err := deserializeSecret(secretStr)
		if err != nil {
			return err
		}

		jwkey, err := jwk.Import(secret)
		if err != nil {
			return err
		}
		jwkey.Remove("d") // "d" is the the private key's field

		htm, err := readValue("htm")
		if err != nil {
			return err
		}

		htu, err := readValue("htu")
		if err != nil {
			return err
		}

		bsh, err := readValue("bsh")
		if err != nil && err.Error() != "unexpected newline" {
			return err
		}

		iat := time.Now().Unix()
		exp := iat + 60

		header := map[string]any{
			"typ": "DPoP+JWT",
			"alg": "EdDSA",
			"jwk": jwkey,
		}

		claims := map[string]any{
			"jti": jti,
			"htm": htm,
			"htu": htu,
			"iat": iat,
			"exp": exp,
		}
		if bsh != "" {
			claims["bsh"] = bsh
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

		signature, err := sign([]byte(tokenB64), secret)
		if err != nil {
			return err
		}
		signatureB64 := base64.RawURLEncoding.EncodeToString(signature)

		log.Infof("DPoP: %s.%s", tokenB64, signatureB64)

		log.Info("DPoP header generated")
		return nil
	},
}

func init() {
	httpCmd.AddCommand(dpopCmd)
}
