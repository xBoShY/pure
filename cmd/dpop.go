package cmd

import (
	"time"

	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/spf13/cobra"
	"go.bryk.io/pkg/errors"
	xlog "go.bryk.io/pkg/log"
)

var dpopCmd = &cobra.Command{
	Use:     "dpop",
	Aliases: []string{"req", "call"},
	Short:   "Generates the DPoP header for an HTTP request",
	Example: "pure req [secret]",
	RunE: func(_ *cobra.Command, args []string) error {
		var err error

		// Get parameters
		if len(args) != 1 {
			return errors.New("you must provide an alias for your did")
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

		headers := map[string]any{
			"typ": "DPoP+JWT",
			"jwk": jwkey,
		}

		claims := map[string]any{
			"htm": htm,
			"htu": htu,
			"iat": iat,
			"exp": exp,
		}

		if bsh != "" {
			claims["bsh"] = bsh
		}

		access, err := sign(secret, headers, claims)
		if err != nil {
			return err
		}

		log.Infof("DPoP: %s", string(access))

		log.WithFields(xlog.Fields{
			"secret": secretStr,
		}).Info("DPoP header generated")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(dpopCmd)
}
