package cmd

import (
	"crypto"
	"crypto/ed25519"
	"encoding/base64"

	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/spf13/cobra"
)

var secretCmd = &cobra.Command{
	Use:     "secret",
	Aliases: []string{},
	Short:   "Secret key tools",
}

func init() {
	rootCmd.AddCommand(secretCmd)
}

func generateSecret() (string, string, error) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		return "", "", err
	}

	return base64.RawURLEncoding.EncodeToString([]byte(pub)),
		base64.RawURLEncoding.EncodeToString([]byte(priv.Seed())),
		nil
}

func deserializePublic(public string) (ed25519.PublicKey, error) {
	b, err := base64.RawURLEncoding.DecodeString(public)
	if err != nil {
		return nil, err
	}

	pub := ed25519.PublicKey(b)
	return pub, nil
}

func deserializeSecret(secret string) (ed25519.PrivateKey, error) {
	b, err := base64.RawURLEncoding.DecodeString(secret)
	if err != nil {
		return nil, err
	}

	priv := ed25519.NewKeyFromSeed(b)
	return priv, nil
}

func jktFromPublicKey(public ed25519.PublicKey) ([]byte, error) {
	jwKey, err := jwk.Import(public)
	if err != nil {
		return nil, err
	}

	jkt, err := jwKey.Thumbprint(crypto.SHA256)
	if err != nil {
		return nil, err
	}

	return jkt, nil
}

func sign(data []byte, secret ed25519.PrivateKey) ([]byte, error) {
	signature := ed25519.Sign(secret, data)

	return signature, nil
}
