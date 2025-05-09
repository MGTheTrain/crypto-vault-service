# crypto-vault-cli

## Summary

`crypto-vault-cli` is a command-line tool for file encryption and decryption using AES, RSA and EC algorithms. It provides an easy interface to securely encrypt and decrypt files using symmetric (AES) and asymmetric (RSA, EC) cryptography.

## Getting Started

### AES example

```sh
uuid=$(cat /proc/sys/kernel/random/uuid)
# Generate AES keys
go run main.go generate-aes-keys --key-size 16 --key-dir data/
# Encryption
go run main.go encrypt-aes --input-file data/input.txt --output-file data/${uuid}-output.enc --symmetric-key <your generated symmetric key>
# Decryption
go run main.go decrypt-aes --input-file data/${uuid}-output.enc --output-file data/${uuid}-decrypted.txt --symmetric-key <your generated symmetric key>
```

### RSA Example

```sh
uuid=$(cat /proc/sys/kernel/random/uuid)

# Generate RSA keys
go run main.go generate-rsa-keys --key-size 2048 --key-dir data/

# Encryption
go run main.go encrypt-rsa --input-file data/input.txt --output-file data/${uuid}-encrypted.txt --public-key <your generated public key>

# Decryption
go run main.go decrypt-rsa --input-file data/${uuid}-encrypted.txt --output-file data/${uuid}-decrypted.txt --private-key <your generated private key>

# Sign
go run main.go sign-rsa --input-file data/input.txt --output-file data/${uuid}-signature.bin --private-key <your generated private key>

# Verify
go run main.go verify-rsa --input-file data/input.txt --signature-file data/${uuid}-signature.bin --public-key <your generated public key>
```

### ECDSA Example

```sh
uuid=$(cat /proc/sys/kernel/random/uuid)

# Generate ECC keys
go run main.go generate-ecc-keys --key-size 256 --key-dir data/

# Sign
go run main.go sign-ecc --input-file data/input.txt  --output-file data/${uuid}-signature.bin --private-key <your generated private key>

# Verify
go run main.go verify-ecc --input-file data/input.txt --signature-file data/${uuid}-signature.bin --public-key <your generated public key>
```

### PKCS#11 example

Make sure the following environment variables are exported as a prerequisite:

```sh
export PKCS11_MODULE_PATH="/usr/lib/softhsm/libsofthsm2.so"
export PKCS11_SO_PIN="1234"
export PKCS11_USER_PIN="5678"
export PKCS11_SLOT_ID="0x0"
```

Next, execute:

```sh
# List token slots
go run main.go list-slots

# Initialize a PKCS#11 token
go run main.go initialize-token --token-label my-token

# Adding keys to tokens
# Add an RSA or EC key pair (consisting of private and public key) to a PKCS#11 token
go run main.go add-key --token-label my-token --object-label my-rsa-key --key-type RSA --key-size 2048
go run main.go add-key --token-label my-token --object-label my-ecdsa-key --key-type ECDSA --key-size 256

# List token objects
go run main.go list-objects --token-label "my-token"

# Deleting keys from tokens
# Delete an object (e.g. RSA or EC key) from the PKCS#11 token
go run main.go delete-object --token-label my-token --object-label my-rsa-key --object-type pubkey
go run main.go delete-object --token-label my-token --object-label my-rsa-key --object-type privkey

# RSA-PKCS
# Encryption
go run main.go encrypt --token-label my-token --object-label my-rsa-key --key-type RSA --input-file data/input.txt --output-file data/encrypted-output.enc

# Decryption
go run main.go decrypt --token-label my-token --object-label my-rsa-key --key-type RSA --input-file data/encrypted-output.enc --output-file data/decrypted-output.txt

# RSA-PSS
# Sign data with a PKCS#11 token
go run main.go sign --token-label my-token --object-label my-rsa-key --key-type RSA --data-file data/input.txt --signature-file data/signature.sig

# Verify the signature using the generated public key from the PKCS#11 token
go run main.go verify --token-label my-token --object-label my-rsa-key --key-type RSA --data-file data/input.txt --signature-file data/signature.sig

# ECDSA
# Sign data with a PKCS#11 token
go run main.go sign --token-label my-token --object-label my-ecdsa-key --key-type ECDSA --data-file data/input.txt --signature-file data/signature.sig

# Verify the signature using the generated public key from the PKCS#11 token
go run main.go verify --token-label my-token --object-label my-ecdsa-key --key-type ECDSA --data-file data/input.txt --signature-file data/signature.sig
```

## e2e-test

An [e2e testing](../../test/e2e/e2e_test.go) the entire flow from encryption to decryption, key management, signing and verifying signatures exists.
