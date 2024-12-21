# Build stage
FROM golang:1.21-alpine AS build

WORKDIR /app

RUN apk update && apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o crypto_vault_service ./cmd/crypto-vault-service/crypto_vault_service.go

# Serve stage
FROM alpine:latest

WORKDIR /root/

RUN apk update && apk add --no-cache libc6-compat
COPY --from=build /app/crypto_vault_service .
COPY ./configs ./configs

EXPOSE 8080

CMD ["./crypto_vault_service"]
