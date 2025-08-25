package vault

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/vault/api"
)

type Client struct {
	vault *api.Client
	path  string
}

type KeyConfig struct {
	Algo Algorithm
	Name string
}

type Key struct {
	Name    string
	Version int64
}

// NewClient creates a new Vault client
func NewClient(config Config, algorithm Algorithm) (*Client, error) {
	vaultConfig := api.DefaultConfig()
	vaultConfig.Address = config.Address

	apiClient, err := api.NewClient(vaultConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	apiClient.SetToken(config.Token)

	// Create client instance
	vc := &Client{
		vault: apiClient,
		path:  config.TransitPath,
	}

	return vc, nil
}

func (c *Client) ListKeys(ctx context.Context) ([]Key, error) {
	path := fmt.Sprintf("%s/keys", c.path)

	secret, err := c.vault.Logical().List(path)
	if err != nil {
		return nil, fmt.Errorf("failed to request new key: %w", err)
	}

	if secret == nil {
		return nil, fmt.Errorf("transit path does not exist: %s", c.path)
	}

	obj, ok := secret.Data["keys"]
	if !ok {
		return nil, fmt.Errorf("failed to read keys")
	}

	keys := []Key{}
	for _, v := range obj.([]any) {
		key := Key{
			Name:    v.(string),
			Version: 1,
		}

		keys = append(keys, key)
	}

	return keys, nil
}

// NewKey requests the transit engine to generate a new key
// with algorithm-specific params
func (c *Client) NewKey(ctx context.Context, cfg KeyConfig) error {
	path := fmt.Sprintf("%s/keys/%s", c.path, cfg.Name)

	params := cfg.Algo.Params()
	_, err := c.vault.Logical().Write(path, params)
	if err != nil {
		return fmt.Errorf("failed to request new key: %w", err)
	}

	return nil
}

// GetPublicKey returns the Public Key for the provided Key
func (c *Client) GetPublicKey(ctx context.Context, key Key) ([]byte, error) {
	versionStr := strconv.FormatInt(key.Version, 10)
	path := fmt.Sprintf("%s/keys/%s", c.path, key.Name)

	secret, err := c.vault.Logical().Read(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key: %w", err)
	}

	if secret == nil {
		return nil, fmt.Errorf("key %s does not exist", key.Name)
	}

	obj, ok := secret.Data["keys"]
	if !ok {
		return nil, fmt.Errorf("failed to read versions")
	}
	versions := obj.(map[string]any)

	obj, ok = versions[versionStr]
	if !ok {
		return nil, fmt.Errorf("failed to read version %s", versionStr)
	}
	retrievedKey := obj.(map[string]any)

	obj, ok = retrievedKey["public_key"]
	if !ok {
		return nil, fmt.Errorf("failed to retrieve public_key")
	}
	publicKeyStr := obj.(string)

	publicKey, err := base64.StdEncoding.DecodeString(publicKeyStr)

	return publicKey, err
}

// GetParams fetch the key configuration from the vault and maps it with the
// appropriate parameters to be used on SignData
func (c *Client) GetParams(ctx context.Context, keyName string) (map[string]any, error) {
	path := fmt.Sprintf("%s/keys/%s", c.path, keyName)
	secret, err := c.vault.Logical().Read(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read key info: %w", err)
	}

	if secret == nil {
		return nil, fmt.Errorf("key not found at path: %s", keyName)
	}

	// TODO: check the type
	keyType := secret.Data["type"].(string)

	alg, err := GetAlgorithm(keyType)
	if err != nil {
		return nil, fmt.Errorf("failed to get algorithm %s: %w", keyType, err)
	}

	params := alg.Params()

	return params, nil
}

// SignData signs data using the transit engine with algorithm-specific params
// Returns the signature in JWS format (base64url encoded)
// Uses marshaling_algorithm=jws for consistent JWT format
func (c *Client) SignData(
	ctx context.Context, data []byte, key Key,
) ([]byte, error) {
	path := fmt.Sprintf("%s/sign/%s", c.path, key.Name)
	params, err := c.GetParams(ctx, key.Name)
	if err != nil {
		return nil, fmt.Errorf("key %s not found: %w", key.Name, err)
	}

	input := base64.StdEncoding.EncodeToString(data)
	params["input"] = input
	params["key_version"] = key.Version

	secret, err := c.vault.Logical().Write(path, params)
	if err != nil {
		return nil, fmt.Errorf("failed to sign data: %w", err)
	}

	if secret == nil {
		return nil, fmt.Errorf("no signature returned")
	}

	signature, ok := secret.Data["signature"].(string)
	if !ok {
		return nil, fmt.Errorf("signature not found in Vault response")
	}

	parts := strings.Split(signature, ":")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid signature format")
	}

	sigStr := parts[2]

	signatureB64, err := base64.StdEncoding.DecodeString(sigStr)
	if err != nil {
		fmt.Printf("\n%s\n", sigStr)
		return nil, err
	}

	return signatureB64, nil
}
