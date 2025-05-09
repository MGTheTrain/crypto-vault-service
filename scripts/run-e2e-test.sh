#!/bin/bash

set -euo pipefail

SCRIPT_DIR=$(dirname "$BASH_SOURCE")
ROOT_PROJECT_DIR=$SCRIPT_DIR/..

cd $ROOT_PROJECT_DIR

BLUE='\033[0;34m'
NC='\033[0m'

echo "#####################################################################################################"

echo -e "$BLUE INFO: $NC Running e2e tests..."
go test ./test/... --tags="e2e" -cover
exit $?

cd $SCRIPT_DIR
