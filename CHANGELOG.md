# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Added [.golangci.yml](./.golangci.yml) with selected linters and formatters
- Introduced a [reusable workflow call](./.github/workflows/_test.yml) utilized by various CI workflows

### Updated

- Placed unit and integration tests next to their respective implementations, with unit test files named `<component>_test.go` and integration test files named `<component>_integration_test.go`
- Followed Go naming conventions by removing the `I` prefix from interfaces and making the structs that implement them private by default, e.g. `type SomeService interface {}` and `type someService struct {}`. This ensures that in this example, `someService` objects can only be created through the `func NewSomeService(...) -> *someService` function
- Refactored `crypto-vault-service` cli tool to rely on environment variables for PKCS#11 operations and refactored e2e-test related to it
- Ran `go fmt ./...` prior to `golangci-lint run` in the [format-and-lint.sh script](./scripts/format-and-lint.sh), and incorporated `shfmt` for shell script formatting and `prettier` for markdown formatting
- Renamed entrypoint files in cmd folder to `main.go`
- Added a Make target to verify that code coverage meets the `70% threshold` across unit and integration tests
- Consolidated standalone scripts into dedicated Make targets for running unit, integration and end-to-end tests

### Fixed

- Enabled use of cancellation contexts in repository components
- Resolved findings from various linters, including `errcheck`, `govet`, `staticcheck`, `wrapcheck`, `importas`, `unused`, `ineffassign`, `errorlint`, `gocritic`, `gosec`, `misspell`, `nakedret` and `revive`
- Fixed `README.md` sections related to commands executed against internal REST and gRPC service APIs
- Specified the database name as an argument in the `QueryRow` function during service startup to ensure proper database creation
- Changed license to MIT: LGPL CLI tools are invoked in internal components only

### Removed

- Removed obsolete `.vscode/launch.json`

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
