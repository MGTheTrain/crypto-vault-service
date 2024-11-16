# crypto-vault-cli

## Table of Contents

- [Summary](#summary)
- [Getting Started](#getting-started)
  - [Encryption and Decryption](#encryption-and-decryption)
    - [AES Example](#aes-example)
    - [RSA Example](#rsa-example)
  - [Signing and Verifying Signatures](#signing-and-verifying-signatures)
    - [ECDSA Example](#ecdsa-example)
- [PKCS#11 Integration](#pkcs11-integration)

## Summary

`crypto-vault-cli` is a command-line tool for file encryption and decryption using AES, RSA and EC algorithms. It provides an easy interface to securely encrypt and decrypt files using symmetric (AES) and asymmetric (RSA, EC) cryptography.

## Getting Started

### Encryption and Decryption

**NOTE**: Keys will be generated internally during the encryption operations.

#### AES example

```sh
uuid=$(cat /proc/sys/kernel/random/uuid)
# Encryption
go run crypto-vault-cli.go encrypt-aes --input data/input.txt --output data/${uuid}-output.enc --keySize 16 --keyDir data/
# Decryption
go run crypto-vault-cli.go decrypt-aes --input data/${uuid}-output.enc --output data/${uuid}-decrypted.txt --symmetricKey <your generated symmetric key from previous encryption operation>
```

#### RSA Example

```sh
uuid=$(cat /proc/sys/kernel/random/uuid)

# Encryption
go run crypto-vault-cli.go encrypt-rsa --input data/input.txt --output data/${uuid}-encrypted.txt --keyDir data/

# Decryption
go run crypto-vault-cli.go decrypt-rsa --input data/${uuid}-encrypted.txt --output data/${uuid}-decrypted.txt --privateKey <your generated private key from previous encryption operation>
```

#### PKCS#11 encryption and decryption

```sh
# Encryption
go run crypto-vault-cli.go encrypt --module /usr/lib/softhsm/libsofthsm2.so --token-label my-token --object-label my-rsa-key --user-pin 5678 --input-file data/input.txt --output-file data/encrypted-output.enc

# Decryption
go run crypto-vault-cli.go decrypt --module /usr/lib/softhsm/libsofthsm2.so --token-label my-token --object-label my-rsa-key --user-pin 5678 --input-file data/encrypted-output.enc --output-file data/decrypted-output.txt
```

---

### Signing and Verifying signatures

**NOTE**: Keys will be generated internally during signature generation operations.

#### ECDSA Example

```sh
# Sign a file with a newly generated ECC key pair (internally generated)
go run crypto-vault-cli.go sign-ecc --input data/input.txt --keyDir data

# Verify the signature using the generated public key
go run crypto-vault-cli.go verify-ecc --input data/input.txt --publicKey <your generated public key from previous signing operation> --signature <your generated signature file from previous signing operation>
```

#### PKCS#11 signing and verifying

```sh
# Sign data with a PKCS#11 token
go run crypto-vault-cli.go sign --module /usr/lib/softhsm/libsofthsm2.so --token-label my-token --object-label my-rsa-key --user-pin 5678 --input-file data/input.txt --output-file data/signature.sig

# Verify the signature using the generated public key from the PKCS#11 token
go run crypto-vault-cli.go verify --module /usr/lib/softhsm/libsofthsm2.so --token-label my-token --object-label my-rsa-key --user-pin 5678 --data-file data/input.txt --signature-file data/signature.sig
```

---

### PKCS#11 Integration

```sh
# Check available slots
pkcs11-tool --module /usr/lib/softhsm/libsofthsm2.so -L
# Initialize a PKCS#11 token
go run crypto-vault-cli.go initialize-token --module /usr/lib/softhsm/libsofthsm2.so --token-label my-token --so-pin 1234 --user-pin 5678 --slot "0x0"

# Check if PKCS#11 token is set
go run crypto-vault-cli.go is-token-set --module /usr/lib/softhsm/libsofthsm2.so --token-label my-token

# Check if an object (e.g., key) exists in the PKCS#11 token
go run crypto-vault-cli.go is-object-set --module /usr/lib/softhsm/libsofthsm2.so --token-label my-token --object-label my-rsa-key --user-pin 5678
# Check all keys
pkcs11-tool --module /usr/lib/softhsm/libsofthsm2.so -O --token-label "my-token" --pin 5678

# Adding keys to tokens
# Add an RSA or ECDSA key pair (private and public key) to a PKCS#11 token
go run crypto-vault-cli.go add-key --module /usr/lib/softhsm/libsofthsm2.so --token-label my-token --object-label my-rsa-key --key-type RSA --key-size 2048 --user-pin 5678

# Deleting keys from tokens
# Delete an object (e.g., RSA or ECDSA key) from the PKCS#11 token
go run crypto-vault-cli.go delete-object --module /usr/lib/softhsm/libsofthsm2.so --token-label my-token --object-label my-rsa-key --object-type pubkey --user-pin 5678
go run crypto-vault-cli.go delete-object --module /usr/lib/softhsm/libsofthsm2.so --token-label my-token --object-label my-rsa-key --object-type privkey --user-pin 5678
```