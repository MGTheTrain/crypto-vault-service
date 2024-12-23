#!/bin/bash

set -euo pipefail

SCRIPT_DIR=$(dirname "$BASH_SOURCE")
ROOT_PROJECT_DIR=$SCRIPT_DIR/..
INTERNAL_GRPC_PROTO_V1_DIR=$ROOT_PROJECT_DIR/internal/api/grpc/v1/proto

cd $INTERNAL_GRPC_PROTO_V1_DIR

BLUE='\033[0;34m'
NC='\033[0m'

echo "#####################################################################################################"
echo -e "$BLUE INFO: $NC About to generate Go gRPC files from proto files"

protoc --go_out=. --go-grpc_out=. ./service.proto
find . -type f -name '*.go' -exec sed -i 's/^package __/package proto/' {} \;

cd $SCRIPT_DIR
