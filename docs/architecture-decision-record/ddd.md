# Domain-driven Design

Incorporating DDD as an architecture decision impacts the structure and development process of the entire application. By aligning the software architecture with the business domain, we promote a more **modular, flexible and maintainable** design. The goal is to deeply understand the domain, represent it through domain models and ensure that the application closely matches the business processes.

```sh
.
├── cmd
│   ├── crypto-vault-cli
│       ├── Dockerfile
│       ├── crypto-vault-cli.go
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
│       └── blob
│       └── key
│       └── permission
│   ├── infrastructure
│   └── persistence
└── test
    ├── data
    ├── integration
    │   ├── domain
    │       └── blob
    │       └── key
    │       └── permission
    │   ├── infrastructure
    │   └── persistence
    └── unit
        ├── domain
        ├── infrastructure
        └── persistence
```