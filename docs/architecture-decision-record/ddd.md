# Domain-driven Design

Incorporating DDD as an architecture decision impacts the structure and development process of the entire application. By aligning the software architecture with the business domain, we promote a more **modular, flexible and maintainable** design. The goal is to deeply understand the domain, represent it through domain models and ensure that the application closely matches the business processes.

```sh
.
├── cmd
│   ├── crypto-vault-cli
│   │   ├── README.md
│   │   ├── crypto-vault-cli.go
│   │   └── data
│   │       ├── decrypted.txt
│   │       ├── decryptedII.txt
│   │       ├── encryptedII.txt
│   │       ├── encryption_key.bin
│   │       ├── input.txt
│   │       ├── output.enc
│   │       ├── private_key.pem
│   │       ├── public_key.pem
│   │       └── signature.sig
│   └── crypto-vault-service
│       ├── Dockerfile
│       └── crypto-vault-service.go
├── configs
│   ├── dev.yml
│   ├── prd.yml
│   └── qas.yml
├── go.mod
├── go.sum
├── internal
│   ├── api
│   │   └── v1
│   │   └── v2
│   │   └── ...
│   ├── app
│   ├── domain
│   ├── infrastructure
│   └── persistence
└── test
    ├── data
    ├── integration
    │   ├── domain
    │   ├── infrastructure
    │   └── persistence
    └── unit
        ├── domain
        ├── infrastructure
        └── persistence
```