# GRPC mTLS example

## Install protoc(for MacOS)

```
brew install protobuf clang-format
```

clang-format is used for format proto files

## Install golang code generator

```
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

## Build proto files

```
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative greet/*.proto
```
## Generate x509 certificates

1. Download [cfssl](https://github.com/cloudflare/cfssl)
2. Run the certGen.sh script to re-generate certificates.

```
sh certGen.sh
```

## Build client and server

```
rm -rf dist && mkdir -p dist
go build -o ./dist ./cmd/...
```

## Run and test

```console
$ ./dist/server
```

```console
$ ./dist/client
```

