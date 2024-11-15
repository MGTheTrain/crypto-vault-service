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
    - [List Commands](#list-commands)
    - [Setup Commands](#setup-commands)
    - [Adding Keys to Tokens](#adding-keys-to-tokens)
    - [Deleting Keys from Tokens](#deleting-keys-from-tokens)


## Summary

`crypto-vault-cli` is a command-line tool for file encryption and decryption using AES, RSA and EC algorithms. It provides an easy interface to securely encrypt and decrypt files using symmetric (AES) and asymmetric (RSA, EC) cryptography.

## Getting Started

**NOTE**: Keys will be generated internally during the encryption or signature generation operations.

### Encryption/Decryption

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

---

### Signing / Verifying signatures

#### ECDSA Example

```sh
# Sign a file with a newly generated ECC key pair (internally generated)
go run crypto-vault-cli.go sign-ecc --input data/input.txt --keyDir data

# Verify the signature using the generated public key
go run crypto-vault-cli.go verify-ecc --input data/input.txt --publicKey <your generated public key from previous signing operation> --signature <your generated signature file from previous signing operation>
```

---

### PKCS#11  

#### List commands

```sh
# Get the next available uninitialized slot in the PKCS#11 module
go run crypto-vault-cli.go get-free-slot --module /usr/lib/softhsm/libsofthsm2.so

# Check if PKCS#11 token is set
go run crypto-vault-cli.go is-token-set --module /usr/lib/softhsm/libsofthsm2.so --token-label my-token

# Check if an object (e.g., key) exists in the PKCS#11 token
go run crypto-vault-cli.go is-object-set --module /usr/lib/softhsm/libsofthsm2.so --token-label my-token --object-label my-rsa-key --user-pin 5678
```

#### Setup commands

```sh
# Initialize a PKCS#11 token
go run crypto-vault-cli.go initialize-token --module /usr/lib/softhsm/libsofthsm2.so --token-label my-token --so-pin 1234 --user-pin 5678

# Set up a PKCS#11 token slot
go run crypto-vault-cli.go setup-slot --module /usr/lib/softhsm/libsofthsm2.so --token-label my-token --so-pin 1234 --user-pin 5678
```

#### Adding keys to tokens

```sh
# Add an RSA or ECDSA key to a PKCS#11 token
go run crypto-vault-cli.go add-key --module /usr/lib/softhsm/libsofthsm2.so --token-label my-token --object-label my-rsa-key --key-type RSA --key-size 2048 --user-pin 5678
```

#### Deleting keys from tokens

```sh
# Delete an object (e.g., RSA or ECDSA key) from the PKCS#11 token
go run crypto-vault-cli.go delete-object --module /usr/lib/softhsm/libsofthsm2.so --token-label my-token --object-label my-rsa-key --object-type privkey --user-pin 5678
```

#### Encryption/Decryption

TBD

#### Signing / Verifying signatures

TBD