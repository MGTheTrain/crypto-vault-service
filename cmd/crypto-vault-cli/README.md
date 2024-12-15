# crypto_vault_cli

## Table of Contents

- [Summary](#summary)
- [Getting Started](#getting-started)
  - [Encryption and Decryption](#encryption-and-decryption)
    - [AES Example](#aes-example)
    - [RSA Example](#rsa-example)
    - [PKCS#11 Encryption and Decryption](#pkcs11-encryption-and-decryption)
  - [Signing and Verifying Signatures](#signing-and-verifying-signatures)
    - [ECDSA Example](#ecdsa-example)
    - [PKCS#11 Signing and Verifying](#pkcs11-signing-and-verifying)
  - [PKCS#11 key management operations](#pkcs11-key-management-operations)

## Summary

`crypto_vault_cli` is a command-line tool for file encryption and decryption using AES, RSA and EC algorithms. It provides an easy interface to securely encrypt and decrypt files using symmetric (AES) and asymmetric (RSA, EC) cryptography.

## Getting Started

### Encryption and Decryption

#### AES example

*NOTE:* Keys will be generated internally during the encryption operations.

```sh
uuid=$(cat /proc/sys/kernel/random/uuid)
# Encryption
go run crypto_vault_cli.go encrypt-aes --input-file data/input.txt --output-file data/${uuid}-output.enc --key-size 16 --key-dir data/
# Decryption
go run crypto_vault_cli.go decrypt-aes --input-file data/${uuid}-output.enc --output-file data/${uuid}-decrypted.txt --symmetric-key <your generated symmetric key from previous encryption operation>
```

#### RSA Example

*NOTE:* Keys will be generated internally during the encryption operations.

```sh
uuid=$(cat /proc/sys/kernel/random/uuid)

# Encryption
go run crypto_vault_cli.go encrypt-rsa --input-file data/input.txt --output-file data/${uuid}-encrypted.txt --key-dir data/

# Decryption
go run crypto_vault_cli.go decrypt-rsa --input-file data/${uuid}-encrypted.txt --output-file data/${uuid}-decrypted.txt --private-key <your generated private key from previous encryption operation>
```

#### PKCS#11 encryption and decryption

*NOTE:* Requires RSA keys managed in FIPS-compliant software or hardware trough `pkcs11-tool` or utilize commands in [PKCS#11 key management operations](#pkcs11-key-management-operations):

```sh
# RSA-PKCS
# Encryption
go run crypto_vault_cli.go encrypt --token-label my-token --object-label my-rsa-key --key-type RSA --input-file data/input.txt --output-file data/encrypted-output.enc

# Decryption
go run crypto_vault_cli.go decrypt --token-label my-token --object-label my-rsa-key --key-type RSA --input-file data/encrypted-output.enc --output-file data/decrypted-output.txt
```

---

### Signing and Verifying signatures

#### ECDSA Example

*NOTE:* Keys will be generated internally during signature generation operations.

```sh
# Sign a file with a newly generated ECC key pair (internally generated)
go run crypto_vault_cli.go sign-ecc --input-file data/input.txt --key-dir data

# Verify the signature using the generated public key
go run crypto_vault_cli.go verify-ecc --input-file data/input.txt --public-key <your generated public key from previous signing operation> --signature-file <your generated signature file from previous signing operation>
```

#### PKCS#11 signing and verifying

*NOTE:* Requires RSA or EC keys managed in FIPS-compliant software or hardware trough `pkcs11-tool` or utilize commands in [PKCS#11 key management operations](#pkcs11-key-management-operations):

```sh
# RSA-PSS
# Sign data with a PKCS#11 token
go run crypto_vault_cli.go sign --token-label my-token --object-label my-rsa-key --key-type RSA --data-file data/input.txt --signature-file data/signature.sig

# Verify the signature using the generated public key from the PKCS#11 token
go run crypto_vault_cli.go verify --token-label my-token --object-label my-rsa-key --key-type RSA --data-file data/input.txt --signature-file data/signature.sig

# ECDSA
# Sign data with a PKCS#11 token
go run crypto_vault_cli.go sign --token-label my-token --object-label my-ecdsa-key --key-type ECDSA --data-file data/input.txt --signature-file data/signature.sig

# Verify the signature using the generated public key from the PKCS#11 token
go run crypto_vault_cli.go verify --token-label my-token --object-label my-ecdsa-key --key-type ECDSA --data-file data/input.txt --signature-file data/signature.sig
```

---

### PKCS#11 key management operations

```sh
# Configure settings
go run crypto_vault_cli.go store-pkcs11-settings --module /usr/lib/softhsm/libsofthsm2.so --so-pin 1234 --user-pin 5678 --slot-id "0x0"

# List token slots
go run crypto_vault_cli.go list-slots

# Initialize a PKCS#11 token
go run crypto_vault_cli.go initialize-token --token-label my-token


# Adding keys to tokens
# Add an RSA or EC key pair (private and public key) to a PKCS#11 token
go run crypto_vault_cli.go add-key --token-label my-token --object-label my-rsa-key --key-type RSA --key-size 2048
go run crypto_vault_cli.go add-key --token-label my-token --object-label my-ecdsa-key --key-type ECDSA --key-size 256

# List token objects
go run crypto_vault_cli.go list-objects --token-label "my-token"

# Deleting keys from tokens
# Delete an object (e.g., RSA or EC key) from the PKCS#11 token
go run crypto_vault_cli.go delete-object --token-label my-token --object-label my-rsa-key --object-type pubkey
go run crypto_vault_cli.go delete-object --token-label my-token --object-label my-rsa-key --object-type privkey
```

### Running the e2e-test

In order to e2e-test the entire flow from encryption to decryption, key management, signing, and verifying signatures as outlined in previous [Getting Started](#getting-started) sections run `go test ./crypto_vault_cli_test.go`.