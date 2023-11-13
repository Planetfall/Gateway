[![Go Reference](https://pkg.go.dev/badge/github.com/planetfall/gateway.svg)](https://pkg.go.dev/github.com/planetfall/gateway)
[![Go Report Card](https://goreportcard.com/badge/github.com/planetfall/gateway)](https://goreportcard.com/report/github.com/planetfall/gateway)
[![codecov](https://codecov.io/gh/planetfall/gateway/graph/badge.svg?token=QWPH8FP2BO)](https://codecov.io/gh/planetfall/gateway)
[![Tests](https://github.com/planetfall/gateway/actions/workflows/gateway.yml/badge.svg)](https://github.com/Planetfall/Gateway/actions/workflows/gateway.yml)
[![Release](https://img.shields.io/github/release/gin-gonic/gin.svg?style=flat-square)](https://github.com/gin-gonic/gin/releases)

# Gateway

Main entrypoint to GRPC microservices.
- Go 1.21
- [genproto](https://github.com/Planetfall/genproto)
- [framework](https://github.com/Planetfall/Framework)

## Run

Get the gateway service account as a key JSON file.
```
gcloud iam service-accounts keys create ./gateway-key.json \
    --iam-account=echo-slam-planetfall-gateway@echo-slam-planetfall.iam.gserviceaccount.com
```

Set the `GOOGLE_APPLICATION_CREDENTIALS` as the path to this JSON file.
```
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/key/json/gateway-key.json
```

Setup your configuration file (ex: config.dev.yaml)

Run
```
go run ./cmd/server/main.go --env development --config ./config/config.dev.yaml
```

## Tests

Run the tests
```
go test ./...
```

Run the tests with coverage
```
go test -v -race -covermode=atomic -coverprofile=coverage.out ./...
```

Print the coverage in HTML
```
go tool cover -html=coverage.out
```

## Lint

Report card
```
goreportcard-cli
```

Golang Lint
```
golangci-lint run
```