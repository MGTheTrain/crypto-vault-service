#!/bin/bash

set -euo pipefail

SCRIPT_DIR=$(dirname "$BASH_SOURCE")
ROOT_PROJECT_DIR=$SCRIPT_DIR/..

cd $ROOT_PROJECT_DIR

BLUE='\033[0;34m'
NC='\033[0m'

# Default flag values
RUN_UNIT_TESTS=true
RUN_INTEGRATION_TESTS=true

# Parse arguments
while getopts "uia" opt; do
  case ${opt} in
    u)
      RUN_UNIT_TESTS=true
      RUN_INTEGRATION_TESTS=false
      ;;
    i)
      RUN_UNIT_TESTS=false
      RUN_INTEGRATION_TESTS=true
      ;;
    a)
      RUN_UNIT_TESTS=true
      RUN_INTEGRATION_TESTS=true
      ;;
    *)
      echo "Usage: $0 [-u] (for unit tests) [-i] (for integration tests) [-a] (running unit and integration tests)"
      exit 1
      ;;
  esac
done

echo "#####################################################################################################"

if [ "$RUN_UNIT_TESTS" = true ] && [ "$RUN_INTEGRATION_TESTS" = true ]; then
  echo -e "$BLUE INFO: $NC Running unit and integration tests..."
  go test ./internal/... --tags="unit integration" -cover
  exit $?
fi

if [ "$RUN_UNIT_TESTS" = true ]; then
  echo -e "$BLUE INFO: $NC Running unit tests..."
  go test ./internal/... --tags=unit -cover
  exit $?
fi

if [ "$RUN_INTEGRATION_TESTS" = true ]; then
  echo -e "$BLUE INFO: $NC Running integration tests..."
  go test ./internal/... --tags=integration -cover
  exit $?
fi

cd $SCRIPT_DIR
