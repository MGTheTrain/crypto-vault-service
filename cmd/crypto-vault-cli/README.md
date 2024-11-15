# crypto-vault-cli

## Table of Contents

+ [Summary](#summary)
+ [Getting started](#getting-started)

## Summary

`crypto-vault-cli` is a command-line tool for file encryption and decryption using AES, RSA and EC algorithms. It provides an easy interface to securely encrypt and decrypt files using symmetric (AES) and asymmetric (RSA, EC) cryptography.

## Getting Started

**NOTE**: Keys will be generated internally during the encryption or signature generation operations.

### Encryption/Decryption

**AES example**

```sh
uuid=$(cat /proc/sys/kernel/random/uuid)
# Encryption
go run crypto-vault-cli.go encrypt-aes --input data/input.txt --output data/${uuid}-output.enc --keySize 16 --keyDir data/
# Decryption
go run crypto-vault-cli.go decrypt-aes --input data/${uuid}-output.enc --output data/${uuid}-decrypted.txt --symmetricKey <your generated symmetric key from previous encryption operation>
```

**RSA Example**

```sh
uuid=$(cat /proc/sys/kernel/random/uuid)

# Encryption
go run crypto-vault-cli.go encrypt-rsa --input data/input.txt --output data/${uuid}-encrypted.txt --keyDir data/

# Decryption
go run crypto-vault-cli.go decrypt-rsa --input data/${uuid}-encrypted.txt --output data/${uuid}-decrypted.txt --privateKey <your generated private key from previous encryption operation>
```

**RSA with PKCS#11 Example** 

```sh
TBD
```

### Hashing / Verifying signatures

**ECDSA Example**

```sh
# Sign a file with a newly generated ECC key pair (internally generated)
go run crypto-vault-cli.go sign-ecc --input data/input.txt --keyDir data

# Verify the signature using the generated public key
go run crypto-vault-cli.go verify-ecc --input data/input.txt --publicKey <your generated public key from previous signing operation> --signature <your generated signature file from previous signing operation>
```