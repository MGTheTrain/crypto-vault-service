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
echo -e "$BLUE INFO: $NC Formatting Go files with go fmt..."
go fmt ./...
golangci-lint run

echo -e "$BLUE INFO: $NC Formatting shell scripts (*.sh) with shfmt..."
find . -name '*.sh' -exec shfmt -w -i 2 -ci {} +

echo -e "$BLUE INFO: $NC Formatting Markdown files (*.md) with prettier..."
find . -name '*.md' -exec prettier --write {} +
prettier --write README.md

cd $SCRIPT_DIR
