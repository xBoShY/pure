
## Contents
- [Contents](#contents)
- [Setup Hashicorp Vault](#setup-hashicorp-vault)
  - [Start Vault](#start-vault)
  - [Initialize Vault](#initialize-vault)
  - [Unseal the Vault](#unseal-the-vault)
  - [Enable HashiCorp Vault Transit Engine](#enable-hashicorp-vault-transit-engine)
- [Using Pure CLI](#using-pure-cli)
  - [Configuration](#configuration)
  - [DIDs](#dids)
    - [Generating DIDs](#generating-dids)
    - [Listing DIDs](#listing-dids)
  - [Secrets](#secrets)
    - [Generating a New Key Pair](#generating-a-new-key-pair)
  - [HTTP Headers / JWT](#http-headers--jwt)
    - [Generating an Authorization Header (DPoP Bound Access Token)](#generating-an-authorization-header-dpop-bound-access-token)
    - [Generating a DPoP Proof Header](#generating-a-dpop-proof-header)
    - [Verifying a JWT Signature](#verifying-a-jwt-signature)


## Setup Hashicorp Vault


### Start Vault
This command starts a HashiCorp Vault server in development mode inside a Docker container.
**Note**: Development mode is for testing only and is not secure for production use.



1. Run the provided script to start the Vault server with Docker:
```shell
./vault/docker.sh
```

The vault is now running in the Docker container called `pureVault`.


### Initialize Vault
When you first start a Vault server, it must be initialized. This process generates the encryption keys and the initial root token.

```shell
docker exec pureVault vault operator init
```

You will see output similar to the following. This is the only time this information will be displayed in its entirety:
```
Unseal Key 1: abc123...
Unseal Key 2: def456...
Unseal Key 3: ghi789...
Unseal Key 4: jkl012...
Unseal Key 5: mno345...

Initial Root Token: hvs.xyz789...
```

**IMPORTANT**: 
- Securely store all five Unseal Keys and the Root Token.
- You will need multiple unseal keys to unseal the Vault after restarts.
- The Root Token has full administrative privileges; treat it like a root password.

The Vault starts in a "sealed" state. You must unseal it before it can be used.


### Unseal the Vault
To unseal the Vault, provide a quorum of unseal keys (3 out of 5).
```shell
docker exec pureVault vault operator unseal <unseal-key-1>
docker exec pureVault vault operator unseal <unseal-key-2>
docker exec pureVault vault operator unseal <unseal-key-3>
```

After providing the required number of keys, the Vault will become unsealed and operational.


### Enable HashiCorp Vault Transit Engine
The Pure CLI uses Vault's Transit engine to manage cryptographic keys for DIDs and secrets. Enable it using the root token.

```shell
docker exec pureVault sh -c 'vault login <your-root-token-here> && vault secrets enable transit'
```
Replace *`<your-root-token-here>`* with the *`Initial Root Token`* from the initialization step.


## Using Pure CLI


### Configuration
Before using the CLI, you must configure it with your Vault server's address and authentication token.

Update the `config.yaml` file. For testing, you can use the root token, but for production, a more restricted token is recommended.

```yaml
address: http://localhost:8200
token: hvs.xyz789... # Replace with your actual root token
transit-path: transit
```


### DIDs


#### Generating DIDs
This command creates a new Decentralized Identifier (DID) and stores the associated private key in your Vault.

```shell
pure did new
```

**Output:**
```log
INF new wallet created did=did:pure:mainnet:U4C254D7KHG362UHT3ABO5QZ3CBW6BRYAJ6JZ4LZ4WO3MO7QK6TWBTPEQI key={"Name":"01f081d6-4f7b-6b2d-be58-7e069b1728d5","Version":1}
```
*Take note of the `did`, as you will need them for signing operations.*


#### Listing DIDs
This command lists all DIDs and their corresponding keys stored in your Vault.
```shell
pure did ls
```

**Output:**
```log
INF did=did:pure:mainnet:U4C254D7KHG362UHT3ABO5QZ3CBW6BRYAJ6JZ4LZ4WO3MO7QK6TWBTPEQI key={"Name":"01f081d6-4f7b-6b2d-be58-7e069b1728d5","Version":1}
INF all keys retrieved
```


### Secrets


#### Generating a New Key Pair
This command generates a new public/private key pair. The private key is secret.
```shell
pure secret new
```

**Output:**
```log
INF secret generated public=HLobDamYvr7ItQTOPVq064-O7EK7Udfb7-iAODk6Qv8 secret=uxdCWX5GweRHA9AAcK0hSGve-MYVsMNKPw2_ruJEwt0
```
*The `public` key is used for verification. The `secret` (private key) is used for signing and must be kept confidential.*


### HTTP Headers / JWT


#### Generating an Authorization Header (DPoP Bound Access Token)
This command generates a long-lived DPoP-bound access token (a JWT). It is signed by your DID and binds the token to a specific public key (which will be used for the DPoP proofs).

**Syntax:**
```shell
pure http auth <your-did> <your-public-key>
```

**Example:**
```shell
pure http auth did:pure:mainnet:U4C254D7KHG362UHT3ABO5QZ3CBW6BRYAJ6JZ4LZ4WO3MO7QK6TWBTPEQI HLobDamYvr7ItQTOPVq064-O7EK7Udfb7-iAODk6Qv8
```

**Output:**
```log
INF header header={"alg":"EdDSA","typ":"JWT"}
INF claims claims={"aud":"http://api.pure.com","cnf":{"jkt":"PDAEkcu6PUpxDa34gRLI0hv32X0v-ozWG7LkA_f68tQ"},"exp":1756142002,"iat":1756141942,"iss":"did:pure:mainnet:U4C254D7KHG362UHT3ABO5QZ3CBW6BRYAJ6JZ4LZ4WO3MO7QK6TWBTPEQI","sub":"did:pure:mainnet:U4C254D7KHG362UHT3ABO5QZ3CBW6BRYAJ6JZ4LZ4WO3MO7QK6TWBTPEQI"}
INF Authorization: DPoP eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJodHRwOi8vYXBpLnB1cmUuY29tIiwiY25mIjp7ImprdCI6IlBEQUVrY3U2UFVweERhMzRnUkxJMGh2MzJYMHYtb3pXRzdMa0FfZjY4dFEifSwiZXhwIjoxNzU2MTQyMDAyLCJpYXQiOjE3NTYxNDE5NDIsImlzcyI6ImRpZDpwdXJlOm1haW5uZXQ6VTRDMjU0RDdLSEczNjJVSFQzQUJPNVFaM0NCVzZCUllBSjZKWjRMWjRXTzNNTzdRSzZUV0JUUEVRSSIsInN1YiI6ImRpZDpwdXJlOm1haW5uZXQ6VTRDMjU0RDdLSEczNjJVSFQzQUJPNVFaM0NCVzZCUllBSjZKWjRMWjRXTzNNTzdRSzZUV0JUUEVRSSJ9.LOFvIAAe319toxR5mjeeykaXcwRnfrXZzP7vPcRUXcjAeRyJkYl79PPtTwi32139ZZYLbVXjk2B604hks0TKAA
INF Authorization header generated did=did:pure:mainnet:U4C254D7KHG362UHT3ABO5QZ3CBW6BRYAJ6JZ4LZ4WO3MO7QK6TWBTPEQI key={"Name":"01f081d6-4f7b-6b2d-be58-7e069b1728d5","Version":1}
```

Use the entire *`Authorization: DPoP ...`* string as the *`Authorization`* header in your HTTP requests.


#### Generating a DPoP Proof Header
This command generates a short-lived DPoP proof JWT for a specific HTTP request. It is signed with the private key bound to the access token.

**Syntax:**
```shell
pure http dpop <your-private-key>
```

**Example:**
```shell
pure http dpop uxdCWX5GweRHA9AAcK0hSGve-MYVsMNKPw2_ruJEwt0
```

The CLI will then prompt you for the details of the HTTP request you are making:
```shell
htm: GET # The HTTP Method (e.g., GET, POST)
htu: https://pure.xboshy.io/api/v1/object # The full URL of the request
bsh: (Optional request body hash) Press Enter to skip.
```

**Output:**
```log
INF header header={"alg":"EdDSA","jwk":{"crv":"Ed25519","kty":"OKP","x":"HLobDamYvr7ItQTOPVq064-O7EK7Udfb7-iAODk6Qv8"},"typ":"DPoP+JWT"}
INF claims claims={"exp":1756142700,"htm":"GET","htu":"https://pure.xboshy.io/api/v1/object","iat":1756142640,"jti":"01f081d8-4334-6fb3-916a-7e069b1728d5"}
INF DPoP: eyJhbGciOiJFZERTQSIsImp3ayI6eyJjcnYiOiJFZDI1NTE5Iiwia3R5IjoiT0tQIiwieCI6IkhMb2JEYW1ZdnI3SXRRVE9QVnEwNjQtTzdFSzdVZGZiNy1pQU9EazZRdjgifSwidHlwIjoiRFBvUCtKV1QifQ.eyJleHAiOjE3NTYxNDI3MDAsImh0bSI6IkdFVCIsImh0dSI6Imh0dHBzOi8vcHVyZS54Ym9zaHkuaW8vYXBpL3YxL29iamVjdCIsImlhdCI6MTc1NjE0MjY0MCwianRpIjoiMDFmMDgxZDgtNDMzNC02ZmIzLTkxNmEtN2UwNjliMTcyOGQ1In0.V6BE_2RBYCwgNjNk5nHd4Nnqntj61KtbKD7G-VluOf3oT-UwPDJYT8KymxDX8sYNdym9zcOftfCijSuq1OrqDQ
INF DPoP header generated
```

*Use the `DPoP: ...` string as the `DPoP` header in your HTTP request.*


#### Verifying a JWT Signature
This command allows you to verify the signature of any JWT (e.g., an Auth token or a DPoP proof) using the appropriate public key.

**Syntax:**
```shell
pure http verify <your-jwt-token>
```

**Example (Verifying an Auth Token):**
```shell
pure http verify eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJhdWQ...
```

**Example (Verifying a DPoP Proof):**
```shell
pure http verify eyJhbGciOiJFZERTQSIsImp3ayI6eyJjcnYiOiJFZDI1NTE5...
```

**Expected Output for a Valid Token:**
```log
INF JWT verified with success!
```
