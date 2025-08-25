#!/usr/bin/env bash

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)

docker run -d --name pureVault \
  -p 8200:8200 \
  -v ${SCRIPT_DIR}/config.hcl:/vault/config/config.hcl \
  --cap-add=IPC_LOCK \
  -e VAULT_ADDR=http://localhost:8200 \
  hashicorp/vault:latest server
