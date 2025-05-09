# crypto-vault-grpc-service

## Table of Contents

+ [Summary](#summary)
+ [Getting started](#getting-started)

## Summary

gRPC service capable of managing cryptographic keys and securing data at rest (metadata, BLOB)

## Getting Started

Set up your IDE with the necessary Go tooling (such as the `delve` debugger or `grpcurl`) or use the provided [devcontainer.json file](../../.devcontainer/devcontainer.json). You can start the service by either running `go run main.go --config ../../configs/grpc-app.yaml` from this directory or by using the `spin-up-docker-containers Make target` from the [Makefile](../../Makefile). 

### List available services

Run `grpcurl -plaintext localhost:50051 list`

The output should resemble:

```sh
grpc.reflection.v1.ServerReflection
grpc.reflection.v1alpha.ServerReflection
internal.BlobDownload
internal.BlobMetadata
internal.BlobUpload
internal.CryptoKeyDownload
internal.CryptoKeyMetadata
internal.CryptoKeyUpload
```

### Upload blob

**NOTE:** Multipart file uploads are not supported with grpc-gateway and `curl`. For more details checkout: `https://grpc-ecosystem.github.io/grpc-gateway/docs/mapping/binary_file_uploads/`. 

Run:

```sh
cd ../../ # Navigate to project root
echo "This is some test content" > task.tmp
grpcurl -import-path ./internal/api/grpc/v1/proto -proto internal/api/grpc/v1/proto/internal/service.proto -d '{
  "file_name": "task.tmp",
  "file_content": "'$(base64 -w 0 task.tmp)'"
}' -plaintext localhost:50051 internal.BlobUpload/Upload
rm task.tmp
```

### List blob metadata

Run `curl -X 'GET' 'http://localhost:8090/api/v1/cvs/blobs' -H 'accept: application/json'`

Optionally:

```sh
cd ../../ # Navigate to project root
grpcurl -import-path ./internal/api/grpc/v1/proto -proto internal/api/grpc/v1/proto/internal/service.proto -d '{
     "name": null,
     "size": null,
     "type": null,
     "date_time_created": null,
     "limit": null,
     "offset": null,
     "sort_by": null,
     "sort_order": null
    }' -plaintext localhost:50051 internal.BlobMetadata/ListMetadata
```

#### Get blob metadata

Run `curl -X 'GET' 'http://localhost:8090/api/v1/cvs/blobs/<blob_id>' -H 'accept: application/json'`

Optionally:

```sh
cd ../../ # Navigate to project root
grpcurl -import-path ./internal/api/grpc/v1/proto -proto internal/api/grpc/v1/proto/internal/service.proto -d '{
    "id": "<blob_id>"
}' -plaintext localhost:50051 internal.BlobMetadata/GetMetadataById
```

### Download blob

Run `curl -X 'GET' 'http://localhost:8090/api/v1/cvs/blobs/<blob_id>/file' -H 'accept: application/json'`

Optionally:

```sh
cd ../../ # Navigate to project root
grpcurl -import-path ./internal/api/grpc/v1/proto -proto internal/api/grpc/v1/proto/internal/service.proto -d '{
    "id": "<blob_id>",
    "decryption_key_id": ""
}' -plaintext localhost:50051 internal.BlobDownload/DownloadById
```

### Delete blob

Run `curl -X 'DELETE' 'http://localhost:8090/api/v1/cvs/blobs/<blob_id>' -H 'accept: application/json'`

Optionally:

```sh
cd ../../ # Navigate to project root
grpcurl -import-path ./internal/api/grpc/v1/proto -proto internal/api/grpc/v1/proto/internal/service.proto -d '{
    "id": "<blob_id>"
}' -plaintext localhost:50051 internal.BlobMetadata/DeleteById
```

### Generate and upload keys

Run:

```sh
cd ../../ # Navigate to project root
grpcurl -import-path ./internal/api/grpc/v1/proto -proto internal/api/grpc/v1/proto/internal/service.proto -d '{
  "algorithm": "RSA",
  "key_size": "2048"
}' -plaintext localhost:50051 internal.CryptoKeyUpload/Upload
```

### List key metadata

Run:

```sh
cd ../../ # Navigate to project root
grpcurl -import-path ./internal/api/grpc/v1/proto -proto internal/api/grpc/v1/proto/internal/service.proto -d '{
    "algorithm": null,
    "type": null,
    "date_time_created": null,
    "limit": null,
    "offset": null,
    "sort_by": null,
    "sort_order": null
}' -plaintext localhost:50051 internal.CryptoKeyMetadata/ListMetadata
```

### Get key metadata

Run: 

```sh
cd ../../ # Navigate to project root
grpcurl -import-path ./internal/api/grpc/v1/proto -proto internal/api/grpc/v1/proto/internal/service.proto -d '{
    "id": "<key_id>"
}' -plaintext localhost:50051 internal.CryptoKeyMetadata/GetMetadataById
```

### Download key

Run:

```sh
cd ../../ # Navigate to project root
grpcurl -import-path ./internal/api/grpc/v1/proto -proto internal/api/grpc/v1/proto/internal/service.proto -d '{
    "id": "<key_id>"
}' -plaintext localhost:50051 internal.CryptoKeyDownload/DownloadById
```

### Delete key

Run:

```sh
cd ../../ # Navigate to project root
curl -X 'DELETE' 'http://localhost:8090/api/v1/cvs/keys/<key_id>' -H 'accept: application/json'
```