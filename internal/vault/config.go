package vault

// Config holds configuration for the Vault client
type Config struct {
	// Address is the Vault server address (e.g., "http://localhost:8200")
	Address string

	// Token is the authentication token
	// Must have permissions for transit engine operations
	Token string

	// TransitPath is the path to the transit engine
	// Example: "transit" at transit/keys
	TransitPath string
}
