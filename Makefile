# Dir where build binaries are generated. The dir should be gitignored
BUILD_OUT_DIR := "bin/"

# go binary. Change this to experiment with different versions of go.
GO       = go

MODULE   = $(shell $(GO) list -m)
SERVICE  = $(shell basename $(MODULE))

# Proto gen info
PROTO_ROOT := proto/
RPC_ROOT := rpc/

BUF_VERSION := v1.17.0
GRPC_GATEWAY_VERSION := v2.15
OPENAPI_VERSION := v2.15
GEN_GO_VERSION := v1.28
GEN_GO_GRPC_VERSION := v1.2
INJECT_TAG_VERSION := v1.4.0

# Fetch OS info
GOVERSION=$(shell go version)
UNAME_OS=$(shell go env GOOS)
UNAME_ARCH=$(shell go env GOARCH)


VERBOSE = 0
Q 		= $(if $(filter 1,$VERBOSE),,@)
M 		= $(shell printf "\033[34;1m▶\033[0m")


BIN 	 = $(CURDIR)/bin
PKGS     = $(or $(PKG),$(shell $(GO) list ./...))

$(BIN)/%: | $(BIN) ; $(info $(M) building package: $(PACKAGE))
	tmp=$$(mktemp -d); \
	   env GOBIN=$(BIN) go get $(PACKAGE) \
		|| ret=$$?; \
	   rm -rf $$tmp ; exit $$ret

$(BIN)/golint: PACKAGE=golang.org/x/lint/golint

GOLINT = $(BIN)/golint

.PHONY: lint
lint: | $(GOLINT) ; $(info $(M) running golint) @ ## Run golint
	$Q $(GOLINT) -set_exit_status $(PKGS)

.PHONY: fmt
fmt: ; $(info $(M) running gofmt) @ ## Run gofmt on all source files
	$Q $(GO) fmt $(PKGS)

.PHONY: all
all: build

.PHONY: deps
deps: proto-deps

.PHONY: proto-deps
proto-deps::
	@echo "\n + Fetching dependencies \n"
	@go install github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION)
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@$(GRPC_GATEWAY_VERSION)
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@$(OPENAPI_VERSION)
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@$(GEN_GO_VERSION)
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(GEN_GO_GRPC_VERSION)
	@buf mod update


.PHONY: proto-generate ## Compile protobuf to pb and twirp files
proto-generate:
	@echo "\n + Generating pb language bindings\n"
	@buf generate

.PHONY: proto-refresh ## Download and re-compile protobuf
proto-refresh: clean proto-generate ## Fetch proto files frrm remote repo

.PHONY: build
build: docker-build

.PHONY: build-info
build-info:
	@echo "\nBuild Info:\n"
	@echo "\t\033[33mOS\033[0m: $(UNAME_OS)"
	@echo "\t\033[33mArch\033[0m: $(UNAME_ARCH)"
	@echo "\t\033[33mGo Version\033[0m: $(GOVERSION)\n"

.PHONY: go-build-api ## Build the binary file for API server
go-build-api:
	@CGO_ENABLED=0 GOOS=$(UNAME_OS) GOARCH=$(UNAME_ARCH) go build -v -o $(API_OUT) $(API_MAIN_FILE)

.PHONY: clean ## Remove previous builds, protobuf files, and proto compiled code
clean:
	@echo " + Removing cloned and generated files\n"
	@rm -rf $(API_OUT) $(RPC_ROOT)

.PHONY: pre-build
pre-build: clean deps proto-refresh

.PHONY: docker-build
docker-build: clean docker-build-api

.PHONY: dev-docker-up ## Bring up docker-compose for local dev-setup
dev-docker-up:
	@COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker compose -f deployment/dev/monitoring/docker-compose.yml up -d
	@COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker compose -f deployment/dev/docker-compose.yml up -d --build

.PHONY: dev-docker-rebuild ## Rebuild
dev-docker-rebuild: dev-docker-up

.PHONY: dev-docker-down ## Shutdown docker-compose for local dev-setup
dev-docker-down:
	@COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker compose -f deployment/dev/docker-compose.yml down --remove-orphans
	@COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker compose -f deployment/dev/monitoring/docker-compose.yml down --remove-orphans

.PHONY: docker-build-api
docker-build-api:
	@DOCKER_BUILDKIT=1 docker build . --secret id=GIT_TOKEN --build-arg=GIT_USERNAME -f build/docker/prod/Dockerfile.api -t razorpay/catalyst:latest

.PHONY: docs-uml ## Generates UML file
docs-uml:
	@go-plantuml generate --recursive --directories cmd --directories internal --directories pkg --out "docs/uml_graph.puml"

.PHONY: docs ## Generates project documentation
docs: docs-uml

check: go-build-api
.PHONY: list-packages

list-packages:
	@touch  app.packages
	@go list ./...| grep e2e -v | grep slit -v > app.packages

