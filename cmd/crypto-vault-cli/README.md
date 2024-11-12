# crypto-vault-cli

`crypto-vault-cli` is a command-line tool for file encryption and decryption using AES, RSA and EC algorithms. It provides an easy interface to securely encrypt and decrypt files using symmetric (AES) and asymmetric (RSA, EC) cryptography.

## Prerequisites

- Install Go from the official Go website, or use this [devcontainer.json](../../.devcontainer/devcontainer.json) with the [DevContainer extensions in VS Code or other IDE supporting DevContainers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)

## Getting Started

### Encryption/Decryption

**AES example**

```sh
# Encryption
go run crypto-cli.go encrypt-aes --input data/input.txt --output data/output.enc --keySize 16 --keyDir data/
# Decryption
go run crypto-cli.go decrypt-aes --input data/output.enc --output data/decrypted.txt --keyDir data/
```

**RSA Example considering external key generation**

```sh
cd data
openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -pubout -in private_key.pem -out public_key.pem
cd -

# Encryption
go run crypto-cli.go encrypt-rsa --input data/input.txt --output data/encryptedII.txt --publicKey data/public_key.pem

# Decryption
go run crypto-cli.go decrypt-rsa --input data/encryptedII.txt --output data/decryptedII.txt --privateKey data/private_key.pem
```

**RSA Example considering internal key generation**

```sh
# Encryption
go run crypto-cli.go encrypt-rsa --input data/input.txt --output data/encryptedII.txt

# Decryption
go run crypto-cli.go decrypt-rsa --input data/encryptedII.txt --output data/decryptedII.txt --privateKey data/private_key.pem
```

**RSA with PKCS#11 Example considering external key generation** 

```sh
TBD
```

**RSA with PKCS#11 Example considering internal key generation** 

```sh
TBD
```

### Hashing / Verifying signatures

**ECDSA Example considering internal key generation**

```sh
# Sign a file with a newly generated ECC key pair (internally generated)
go run crypto-cli.go sign-ecc --input data/input.txt --keyDir data

# Verify the signature using the generated public key
go run crypto-cli.go verify-ecc --input data/input.txt --publicKey data/public_key.pem --signature data/signature.sig
```