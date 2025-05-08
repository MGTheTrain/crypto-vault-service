#!/bin/bash

set -euo pipefail

SCRIPT_DIR=$(dirname "$BASH_SOURCE")
ROOT_PROJECT_DIR=$SCRIPT_DIR/..

cd $ROOT_PROJECT_DIR

BLUE='\033[0;34m'
NC='\033[0m'

echo "#####################################################################################################"
echo -e "$BLUE INFO: $NC About to apply auto-formatting"

# "golangci-lint - Fast linters runners for Go. Bundle of gofmt, govet, errcheck, staticcheck, revive and many other linters. Recommended by the original author to replace gometalinter (Drop-in replacement).""
# Refer to: https://go.dev/wiki/CodeTools
golangci-lint run

cd $SCRIPT_DIR
