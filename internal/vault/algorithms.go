package vault

import (
	"fmt"

	"github.com/xboshy/pure/internal/vault/algorithms"
)

type Algorithm interface {
	// VaultKeyType returns the required Vault Transit key type
	// Examples: "ecdsa-p256", "rsa-2048"
	Name() string

	// SigningParams returns algorithm-specific Vault signing parameters
	// Including base params (prehashed, hash_algorithm) and algorithm specific ones
	// All algorithms use marshaling_algorithm=jws
	Params() map[string]any

	// Verify verifies the signature against the message using given public key
	// Key must be *ecdsa.PublicKey or *rsa.PublicKey matching the algorithm
	Verify(message, signature []byte, key any) error

	// KeyCheck validates the key type for verification
	// Returns ErrInvalidKeyType if key type doesn't match algorithm
	KeyCheck(key any) error
}

var registry map[string]Algorithm

func init() {
	registry = make(map[string]Algorithm)

	var alg Algorithm

	alg = algorithms.Ed25519()
	registry[alg.Name()] = alg
}

func GetAlgorithm(name string) (Algorithm, error) {
	alg, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("algorithm found: %s", name)
	}
	return alg, nil
}
