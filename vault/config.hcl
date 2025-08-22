listener "tcp" {
  address     = "0.0.0.0:8200"    # Listen on port 8200
  tls_disable = 1                  # Disable SSL for local testing
}

storage "file" {
  path = "/vault/file"             # Where to store Vault data
}

disable_mlock = true               # Don't lock memory (for Docker)

ui = true                          # Enable web interface
