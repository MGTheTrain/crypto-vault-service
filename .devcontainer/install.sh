#!/bin/bash

apt-get update
apt-get install -y openssl opensc softhsm libssl-dev libengine-pkcs11-openssl protobuf-compiler
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest