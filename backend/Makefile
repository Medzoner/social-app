.PHONY: help build test bench fmt vet lint docker release fuzz clean
help: Makefile
	@sed -n 's|^##||p' $< | column -t -s ':' | sort

.PHONY: init
init:
	go install "github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
	go install "golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@master"
	go install "github.com/google/wire/cmd/wire@latest"
	go install go.uber.org/mock/mockgen@latest
	go mod tidy

APP_NAME=social-app
VERSION := $(shell git describe --tags --always)
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -ldflags="-s -w -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)"

init: wire

.PHONY: wire
wire:
	wire ./...

docs: swagger

.PHONY: swag
swag:
	swag init

check: fmt vet lint bench fuzz

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint:
	golangci-lint run -c ./golangci.yml -v --fix

.PHONY: fuzz
fuzz:
	go test -fuzz=FuzzTransform -fuzztime=10s ./internal/domains/transform/handler

.PHONY: bench
bench:
	go test -bench=. -benchmem ./internal/domains/transform

#(pprof) top
#(pprof) list transform
#(pprof) web  # → si tu as graphviz
pprof-cpu:
	go tool pprof http://localhost:6060/debug/pprof/profile?seconds=10

pprof-heap:
	go tool pprof http://localhost:6060/debug/pprof/heap
