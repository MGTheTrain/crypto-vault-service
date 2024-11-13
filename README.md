# crypto-vault-service

## Table of Contents

+ [Summary](#summary)
+ [References](#references)
+ [Features](#features)
+ [Getting started](#getting-started)

## Summary

RESTful Web API for managing cryptographic keys and securing data at rest (metadata, BLOB)

## References

TBD

## Features

### Functional

- [ ] **Provide RESTful API for cryptographic operations**: Expose endpoints for managing cryptographic material and securing data (files, metadata) at rest.
- [x] **Asymmetric encryption and decryption**: Support RSA and other asymmetric encryption algorithms for data protection.
- [x] **Symmetric encryption**: Support for symmetric key encryption (e.g. AES) for data protection.
- [x] **Hashing and signature verification**: Support hashing algorithms (e.g. SHA-256, SHA-512) and verify signatures using asymmetric keys (RSA, ECDSA, etc.).
- [ ] **PKCS#11 integration**:  Enable key management in FIPS-compliant hardware or software.
- [ ] **Manage cryptographic material**: Enable management of private/public key pairs and symmetric keys (generation, import/export, rotation, etc.).
- [ ] **Key management lifecycle**: Implement key lifecycle management (generation, rotation, revocation, expiration).
- [ ] **Secure file storage integration**: Provide mechanisms to securely store encrypted files in BLOB storage (e.g. AWS S3, Azure Blob Storage, Google Cloud Storage).
- [ ] **Access control**: Implement role-based access control (RBAC) for APIs ensuring that only authorized users can perform operations on cryptographic material.

### Non-functional

- [x] **Scalable and maintainable project structure**: Refer to the [project-layout GitHub repo](https://github.com/golang-standards/project-layout) and adopt Domain-Driven Design to create a **modular, flexible and maintainable** project structure with a focus on the **domain at its core**
- [x] **CI workflows for quality checks**: Set up continuous integration workflows with GitHub Actions for automated linting, functional and non-functional testing, building and pushing artifacts.
- [ ] **Security checks in CI workflows**: Consider non-functional testing (vulnerability scanning, SBOM generation, Static Code Analysis) in GitHub Actions.
- [ ] **Performance optimization**: Ensure cryptographic operations are optimized for performance, especially for large files and high throughput environments.
- [ ] **Logging and monitoring**: Integrate logging (e.g. using structured logging with `logrus`) and monitoring (e.g. Prometheus, Grafana) to track API usage, performance and errors.
- [ ] **Error handling and resiliency**: Implement comprehensive error handling and retries for operations that may fail, with clear error messages and status codes for the API.
- [ ] **Security**: Ensure all cryptographic material is securely stored and encrypted, protect APIs with authentication (e.g. OAuth2, JWT) and follow best practices for handling sensitive data.
- [ ] **Documentation**: Provide clear API documentation (e.g. Swagger/OpenAPI) for ease of integration by other developers.
- [ ] **Versioning**: Implement proper API versioning to maintain backward compatibility as the API evolves.
- [ ] **Audit logging**: Maintain logs of all cryptographic operations and key management activities for compliance and auditing purposes.


## Getting Started

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

To run `unit` tests on Unix systems either execute

```sh
cd scripts
./run-test.sh -u
```

or

```sh
make run-unit-tests
```

**TBD** To run `integration` tests on Unix systems either execute

```sh
cd scripts
./run-test.sh -i
```

**TBD** or

```sh
make run-integration-tests
```

### Applications

You can find applications utilizing [internal packages](./internal/) in the [cmd folder](./cmd/).