GOOS ?= $(shell go env GOOS || echo linux)
GOARCH ?= $(shell go env GOARCH || echo amd64)
CGO_ENABLED ?= 0

vendor: ## vendor
	go mod download

build: ## build binary files
	GOOS=linux GOARCH=amd64 go build -o bin/app-linux-amd64 main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/app-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/app-darwin-arm64 main.go

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
