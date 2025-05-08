#!/bin/bash

set -euo pipefail

SCRIPT_DIR=$(dirname "$BASH_SOURCE")
ROOT_PROJECT_DIR=$SCRIPT_DIR/..

cd $ROOT_PROJECT_DIR

BLUE='\033[0;34m'
NC='\033[0m'

echo "#####################################################################################################"
echo -e "$BLUE INFO: $NC About to convert Go annotations to Swagger Documentation 2.0"

swag init -g cmd/crypto-vault-rest-service/crypto_vault_service.go -o cmd/crypto-vault-rest-service/docs

cd $SCRIPT_DIR
