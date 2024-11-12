# crypto-vault-cli

`crypto-vault-cli` is a command-line tool for file encryption and decryption using AES, RSA and EC algorithms. It provides an easy interface to securely encrypt and decrypt files using symmetric (AES) and asymmetric (RSA, EC) cryptography.

## Prerequisites

Before you begin, ensure you have the following tools installed:

- Install Go from the official Go website, or use this [devcontainer.json](../../.devcontainer/devcontainer.json) with the [DevContainer extensions in VS Code or other IDE supporting DevContainers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)

## Getting Started

## AES example

```sh
# Encryption
go run crypto-cli.go encrypt-aes --input data/input.txt --output data/output.enc --keySize 16 --keyDir data/
# Decryption
go run crypto-cli.go decrypt-aes --input data/output.enc --output data/decrypted.txt --keyDir data/
```

## RSA Example

### External key generation

```sh
cd assets
openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -pubout -in private_key.pem -out public_key.pem
cd -

# Encryption
go run crypto-cli.go encrypt-rsa --input data/input.txt --output data/encryptedII.txt --publicKey data/public_key.pem

# Decryption
go run crypto-cli.go decrypt-rsa --input data/encryptedII.txt --output data/decryptedII.txt --privateKey data/private_key.pem
```

### Internal key generation

```sh
# Encryption
go run crypto-cli.go encrypt-rsa --input data/input.txt --output data/encryptedII.txt

# Decryption
go run crypto-cli.go decrypt-rsa --input data/encryptedII.txt --output data/decryptedII.txt --privateKey data/private_key.pem
```

## EC Example

### External key generation

```sh
TBD
```

### Internal key generation

```sh
TBD
```

## RSA with PKCS#11 

### External key generation

```sh
TBD
```

### Internal key generation

```sh
TBD
```