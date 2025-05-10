#!/bin/bash

apt-get update
apt-get install -y openssl opensc softhsm libssl-dev libengine-pkcs11-openssl protobuf-compiler bc
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
go install golang.org/x/tools/cmd/goimports@latest

# shfmt for bash script formatting
apt-get install -y shfmt

# prettier for Markdown formatting
npm install -g prettier
