# crypto-vault-service

## Table of Contents

+ [Summary](#summary)
+ [References](#references)
+ [Features](#features)
+ [Getting started](#getting-started)
+ [Documentation](#documentation)

## Summary

RESTful Web API for managing cryptographic keys and securing data at rest (metadata, BLOB)

## References

- [OpenSSL with libp11 for Signing, Verifying and Encrypting, Decrypting](https://docs.yubico.com/hardware/yubihsm-2/hsm-2-user-guide/hsm2-openssl-libp11.html#rsa-pkcs)
- [pkcs11-tool usage](https://docs.nitrokey.com/nethsm/pkcs11-tool#id1)
- [OpenFGA online editor](https://play.fga.dev/sandbox/?store=github)

## Features

### Functional

- [ ] **Provide RESTful API for cryptographic operations**: Expose endpoints for managing cryptographic material and securing data (files, metadata) at rest.
- [x] **Asymmetric encryption and decryption**: Support RSA encryption algorithm for data protection.
- [x] **Symmetric encryption**: Support for symmetric key encryption (e.g. AES) for data protection.
- [x] **Signature creation and verification:** Support for hashing algorithms (e.g. SHA-256, SHA-512) to create digital signatures, and the ability to verify these signatures using asymmetric keys (RSA, ECDSA).
- [x] **PKCS#11 integration**: Enable key management and cryptographic operations (such as RSA-PKCS encryption/decryption and RSA-PSS or ECDSA signing/verification) through PKCS#11 interfaces supporting both FIPS-compliant hardware and software environments.
- [ ] **Manage cryptographic material**: Enable management of private/public key pairs and symmetric keys (generation, import/export, rotation, etc.).
- [ ] **Key management lifecycle**: Implement key lifecycle management (generation, rotation, revocation, expiration).
- [ ] **Secure file storage integration**: Provide mechanisms to securely store encrypted files in BLOB storage (e.g. AWS S3, Azure Blob Storage, Google Cloud Storage).
- [ ] **Access control**:  Implement relationship-based access control (ReBAC) for APIs, ensuring that users can only perform operations on cryptographic material based on their defined relationships and permissions within the system.

### Non-functional

- [x] **Scalable and maintainable project structure**: Refer to the [project-layout GitHub repo](https://github.com/golang-standards/project-layout) and adopt Domain-Driven Design to create a **modular, flexible and maintainable** project structure with a focus on the **domain at its core**
- [x] **CI workflows for quality checks**: Set up continuous integration workflows with GitHub Actions for automated linting, functional and non-functional testing, building and pushing artifacts.
- [ ] **Security checks in CI workflows**: Consider non-functional testing (vulnerability scanning, SBOM generation, Static Code Analysis) in GitHub Actions.
- [ ] **Performance optimization**: Ensure cryptographic operations are optimized for performance, especially for large files and high throughput environments.
- [ ] **Logging and monitoring**: Integrate logging (e.g. using structured logging with `logrus`) and monitoring (e.g. Prometheus, Grafana) to track API usage, performance and errors.
- [ ] **Error handling and resiliency**: Implement comprehensive error handling and retries for operations that may fail, with clear error messages and status codes for the API.
- [ ] **Security**: Ensure that all cryptographic material is securely encrypted before storing it in a key vault using a master key. Additionally, protect APIs with authentication mechanisms such as OAuth2 or JWT, and follow best practices for handling sensitive data.
- [ ] **Documentation**: Provide clear API documentation (e.g. Swagger/OpenAPI) for ease of integration by other developers.
- [ ] **Versioning**: Implement proper API versioning to maintain backward compatibility as the API evolves.
- [ ] **Audit logging**: Maintain logs of all cryptographic operations and key management activities for compliance and auditing purposes.


## Getting Started

### Preconditions

- Install Go from the official Go website, or use this [devcontainer.json](../../.devcontainer/devcontainer.json) with the [DevContainer extensions in VS Code or other IDE supporting DevContainers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)
- If the `devcontainer.json` is not used, install the necessary dependencies for PKCS#11 integration on a later Linux distribution such as `Debian 12` or `Ubuntu 22.04`: 

```sh
apt-get update 
apt-get install -y openssl opensc softhsm libssl-dev libengine-pkcs11-openssl
```

### Formatting and linting

For formatting and linting run either on Unix systems

```sh
cd scripts
./format-and-lint.sh
```

or

```sh
make format-and-lint
```

### Run Tests

To run `unit tests` on Unix systems execute

```sh
make run-unit-tests
```

To run `integration tests` on Unix systems execute

```sh
make spin-up-integration-test-docker-containers
make run-integration-tests
make shut-down-integration-test-docker-containers # Optionally clear docker resources
```

### Applications

You can find applications utilizing [internal packages](./internal/) in the [cmd folder](./cmd/).

### Documentation

You can find documentation on architectural decisions, diagrams and concepts [here](./docs).
