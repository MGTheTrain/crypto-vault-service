# Build stage
FROM golang:1.23-alpine AS build

WORKDIR /app

RUN apk update
COPY . .
RUN go mod tidy
RUN go build -o crypto_vault_service ./cmd/crypto-vault-rest-service/crypto_vault_service.go

# Serve stage
FROM alpine:latest

WORKDIR /root/

RUN apk update && apk add --no-cache libc6-compat
COPY --from=build /app/crypto_vault_service .

EXPOSE 8080

ENTRYPOINT ["./crypto_vault_service"]
