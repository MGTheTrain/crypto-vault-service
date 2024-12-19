# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
- **Asymmetric encryption and decryption**: Supported RSA encryption algorithm for data protection.
- **Symmetric encryption**: Supported symmetric key encryption (e.g. AES) for data protection.
**Signature creation and verification:** Support for hashing algorithms (e.g. SHA-256, SHA-512) to create digital signatures, and the ability to verify these signatures using asymmetric keys (RSA, ECDSA).
- **Scalable and maintainable project structure**: Referred to the [project-layout GitHub repo](https://github.com/golang-standards/project-layout) and adopted Domain-Driven Design to create a **modular, flexible and maintainable** project structure with a focus on the **domain at its core**
- **CI workflows for quality checks**: Set up continuous integration workflows with GitHub Actions for automated linting, functional testing, building and pushing artifacts.
- **PKCS#11 integration**: Enabled key management and cryptographic operations (such as RSA-PKCS encryption/decryption and RSA-PSS or ECDSA signing/verification) through PKCS#11 interfaces supporting both FIPS-compliant hardware and software environments.
- **Logging**: Integrated console and file logging (e.g. using structured logging with `logrus`) 
- **Manage cryptographic material**: Enable management of private/public key pairs and symmetric keys and implement key lifecycle management including key generation and key export
- **Secure file storage integration**: Provide mechanisms to securely store encrypted files in Azure Blob Storage 

## [0.1.0] - TBD-TBD-TBD

### Added

TBD