package cmd

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jws"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

func generateSecret() (string, error) {
	_, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString([]byte(priv.Seed())), nil
}

func deserializeSecret(secret string) (ed25519.PrivateKey, error) {
	b, err := base64.RawURLEncoding.DecodeString(secret)
	if err != nil {
		return nil, err
	}

	priv := ed25519.NewKeyFromSeed(b)
	return priv, nil
}

func sign(secret ed25519.PrivateKey, headers, claims map[string]any) ([]byte, error) {
	c := jwt.New()

	for k, v := range claims {
		c.Set(k, v)
	}

	var ho []jws.Option
	if headers != nil {
		h := jws.NewHeaders()
		for k, v := range headers {
			h.Set(k, v)
		}

		ho = append(ho, jws.WithProtectedHeaders(h))
	}

	signed, err := jwt.Sign(c, jwt.WithKey(jwa.EdDSA(), secret, ho...))
	if err != nil {
		return nil, err
	}

	return signed, nil
}

func readValue(prompt string) (string, error) {
	var dest string
	fmt.Printf("%s: ", prompt)
	_, err := fmt.Scanln(&dest)
	return dest, err
}
