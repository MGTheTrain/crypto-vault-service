# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Added [.golangci.yml](./.golangci.yml) with selected linters and formatters

### Updated

- Followed Go idioms by relocating unit and integration tests alongside their corresponding implementations
- Invoked `go fmt ./...` before `golangci-lint run` in [format-and-lint.sh script](./scripts/format-and-lint.sh)
- Renamed entrypoint files in cmd folder to `main.go` 

### Fixed 

- Enabled use of cancellation contexts in repository components
- Addressed linter findings related to error wrapping, using switch statements over if/else, `import` dependency ordering, handling all errors

## [0.3.0] - 01-02-2025

### Updated

- Considered passing context as an input argument to manage concurrent operations (control cancellation, set timeouts/deadlines, propagate values).
- Logged code coverage when running unit or integration tests

## [0.2.0] - 14-01-2025

### Updated

- Relocated repository interfaces to the domain layer adhering to domain-driven design concepts
- Utilized `ubuntu-22.04` tag for runners 
- Eliminated the outdated GitHub run number suffix for release Docker images in the release CI pipeline

### Fixed

- Fixed `postCreateCommand` in [devcontainer.json file](./.devcontainer/devcontainer.json)

## [0.1.0] - 23-12-2024

### Added

- **Asymmetric encryption and decryption**: Supported RSA encryption algorithm for data protection.
- **Symmetric encryption**: Supported symmetric key encryption (e.g. AES) for data protection.
**Signature creation and verification:** Support for hashing algorithms (e.g. SHA-256, SHA-512) to create digital signatures, and the ability to verify these signatures using asymmetric keys (RSA, ECDSA).
- **Scalable and maintainable project structure**: Referred to the [project-layout GitHub repo](https://github.com/golang-standards/project-layout) and adopted Domain-Driven Design to create a **modular, flexible and maintainable** project structure with a focus on the **domain at its core**
- **CI workflows for quality checks**: Set up continuous integration workflows with GitHub Actions for automated linting, functional testing, building and pushing artifacts.
- **PKCS#11 integration**: Enabled key management and cryptographic operations (such as RSA-PKCS encryption/decryption and RSA-PSS or ECDSA signing/verification) through PKCS#11 interfaces supporting both FIPS-compliant hardware and software environments.
- **Logging**: Integrated console and file logging (e.g. using structured logging with `logrus`) 
- **Manage cryptographic material**: Enabled management of private/public key pairs and symmetric keys and implemented key lifecycle management including primarily key generation and key export
- **Secure file storage integration**: Provided mechanisms to securely store encrypted files in Azure Blob Storage 
- **RESTful API**: Provided HTTP endpoints to manage cryptographic material and secure data at rest.
- **Documentation**: Provided clear API documentation (e.g. Swagger/OpenAPI) for ease of integration by other developers.
- **Versioning**: Implemented proper API versioning to maintain backward compatibility as the API evolves.
- **gRPC API**: Provided gRPC endpoints to manage cryptographic material and secure data at rest