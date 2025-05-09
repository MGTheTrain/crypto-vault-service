# crypto-vault-rest-service

## Summary

REST service capable of managing cryptographic keys and securing data at rest (metadata, BLOB)

## Getting Started

Set up your IDE with the necessary Go tooling (such as the `delve` debugger) or use the provided [devcontainer.json file](../../.devcontainer/devcontainer.json). You can start the service by either running `go run main.go --config ../../configs/rest-app.yaml` from this directory or by using the `spin-up-docker-containers Make target` from the [Makefile](../../Makefile). To explore the Swagger Web UI you need to either visit `http://localhost:8080/api/v1/cvs/swagger/index.html`.
