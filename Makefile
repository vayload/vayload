APP_SERVER=vayload-server
APP_CLI=vayload
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

## Directories
BIN_DIR=./bin
DIST_DIR=./dist
RELEASE_DIR=./release
INSTALL_DIR=$(HOME)/bin
CMD_SERVER=./cmd/server
CMD_CLI=./cmd/cli

## Platform detection
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

## Extension for executables
EXT=
ifeq ($(GOOS),windows)
	EXT=.exe
endif

## Build information
LDFLAGS_VERSION=-X 'main.Version=$(VERSION)' \
	-X 'main.BuildTime=$(BUILD_TIME)' \
	-X 'main.GitCommit=$(GIT_COMMIT)'

## Platform-specific optimizations
LDFLAGS=-s -w $(LDFLAGS_VERSION)
GCFLAGS=
BUILD_TAGS=

# Linux optimizations
ifeq ($(GOOS),linux)
	LDFLAGS += -extldflags '-static'
	GCFLAGS += -trimpath
	BUILD_TAGS += netgo osusergo
endif

# macOS optimizations
ifeq ($(GOOS),darwin)
	GCFLAGS += -trimpath
	ifeq ($(GOARCH),arm64)
		GCFLAGS += -N -l
	endif
endif

# Windows optimizations
ifeq ($(GOOS),windows)
	BUILD_TAGS += netgo
endif

## CGO settings for static builds
CGO_ENABLED ?= 1
ifeq ($(GOOS),darwin)
	CGO_ENABLED = 1
endif

## Build mode and optimization level
BUILD_MODE ?= release
ifeq ($(BUILD_MODE),debug)
	LDFLAGS = -w $(LDFLAGS_VERSION)
	GCFLAGS += -N -l
endif

GO_BUILD=CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) \
	go build -v \
		-ldflags='$(LDFLAGS)' \
		-gcflags='$(GCFLAGS)' \
		-tags='$(BUILD_TAGS)'

.PHONY: all build build-server build-cli build-client install-cli clean \
	gen-fmc-keys gen-pair-keys test lint \
	docker-build docker-push

## Default target
all: build

## Help
help:
	@echo 'Vayload Build System'
	@echo ''
	@echo 'Usage:'
	@echo '  make [target]'
	@echo ''
	@echo 'Targets:'
	@echo '  build           Build all components (server, cli, client)'
	@echo '  build-server    Build server binary'
	@echo '  build-cli       Build CLI binary'
	@echo '  build-client    Build frontend client'
	@echo '  install-cli     Build and install CLI to ~/bin'
	@echo '  clean           Clean build artifacts'
	@echo '  test            Run tests'
	@echo '  lint            Run linter'
	@echo '  docker-build    Build Docker image'
	@echo ''
	@echo 'Variables:'
	@echo '  VERSION         Version string (default: git describe)'
	@echo '  GOOS            Target OS (default: current)'
	@echo '  GOARCH          Target architecture (default: current)'
	@echo '  BUILD_MODE      Build mode: release|debug (default: release)'
	@echo ''
	@echo 'Examples:'
	@echo '  make build'
	@echo '  make build GOOS=linux GOARCH=amd64'

## Build all components
build: build-server build-cli build-client
	@echo ''
	@echo '  ================================================='
	@echo '  >> Complete build for $(GOOS)/$(GOARCH)'
	@echo '  >> Version: $(VERSION)'
	@echo '  ================================================='

## Build server
build-server:
	@echo '  ================================================='
	@echo '  >> Building VAYLOAD-SERVER'
	@echo '  >> OS:        $(GOOS)'
	@echo '  >> ARCH:      $(GOARCH)'
	@echo '  >> CGO:       $(CGO_ENABLED)'
	@echo '  >> VERSION:   $(VERSION)'
	@echo '  >> COMMIT:    $(GIT_COMMIT)'
	@echo '  >> LDFLAGS:   $(LDFLAGS)'
	@echo '  >> GCFLAGS:   $(GCFLAGS)'
	@echo '  >> TAGS:      $(BUILD_TAGS)'
	@echo '  ================================================='
	@mkdir -p $(BIN_DIR)
	$(GO_BUILD) -o $(BIN_DIR)/$(APP_SERVER)$(EXT) $(CMD_SERVER)
	@echo ''
	@echo '  >> Server binary generated:'
	@ls -lh $(BIN_DIR)/$(APP_SERVER)$(EXT)
	@echo '  ================================================='

## Build CLI
build-cli:
	@echo '  ================================================='
	@echo '  >> Building VAYLOAD-CLI'
	@echo '  >> OS:        $(GOOS)'
	@echo '  >> ARCH:      $(GOARCH)'
	@echo '  >> CGO:       $(CGO_ENABLED)'
	@echo '  >> VERSION:   $(VERSION)'
	@echo '  >> COMMIT:    $(GIT_COMMIT)'
	@echo '  >> LDFLAGS:   $(LDFLAGS)'
	@echo '  >> GCFLAGS:   $(GCFLAGS)'
	@echo '  >> TAGS:      $(BUILD_TAGS)'
	@echo '  ================================================='
	@mkdir -p $(BIN_DIR)
	$(GO_BUILD) -o $(BIN_DIR)/$(APP_CLI)$(EXT) $(CMD_CLI)
	@echo ''
	@echo '  >> CLI binary generated:'
	@ls -lh $(BIN_DIR)/$(APP_CLI)$(EXT)
	@echo '  ================================================='

## Build frontend client
build-client:
	@echo '  ================================================='
	@echo '  >> Building frontend client'
	@echo '  ================================================='
	@if [ ! -d "web" ]; then \
		echo '  >> No web directory found, skipping frontend build'; \
		exit 0; \
	fi
	@if ! command -v node >/dev/null 2>&1; then \
		echo '  >> Node.js not installed, skipping frontend build'; \
		exit 0; \
	fi
	@if ! command -v pnpm >/dev/null 2>&1; then \
		echo '  >> pnpm not installed, skipping frontend build'; \
		exit 0; \
	fi
	@if [ ! -d "web/node_modules" ]; then \
		echo '  >> Installing frontend dependencies (pnpm install)'; \
		cd web && pnpm install; \
	else \
		echo '  >> node_modules found, skipping install'; \
	fi
	@echo '  >> Building frontend (pnpm run build)'
	@cd web && pnpm run build
	@echo '  >> Frontend build complete'

## Install CLI
install-cli: build-cli
	@echo '  >> Installing CLI'
	@mkdir -p $(INSTALL_DIR)
	@cp $(BIN_DIR)/$(APP_CLI)$(EXT) $(INSTALL_DIR)/$(APP_CLI)
	@echo '  >> Installed $(APP_CLI) into $(INSTALL_DIR)'
	@echo '  >> Binary size:'
	@ls -lh $(INSTALL_DIR)/$(APP_CLI)
	@echo '  >> Ensure $$HOME/bin is in your PATH'

## Run tests
test:
	@echo '  >> Running tests'
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo '  >> Coverage report generated: coverage.html'

## Run linter
lint:
	@echo '  >> Running linter'
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo '  >> golangci-lint not installed, skipping'; \
		echo '  >> Install: https://golangci-lint.run/usage/install/'; \
	fi

## Clean build artifacts
clean:
	@echo '  >> Cleaning build artifacts'
	@rm -rf $(BIN_DIR)/*
	@rm -rf $(DIST_DIR)/*
	@rm -rf $(RELEASE_DIR)/*
	@rm -f $(INSTALL_DIR)/$(APP_CLI)
	@rm -rf web/build web/node_modules
	@rm -f coverage.out coverage.html
	@echo '  >> Clean completed'

## Generate FMC keys (use: make gen-fmc-keys SECRET=your_secret)
gen-fmc-keys:
	@if [ -z "$(SECRET)" ]; then \
		echo "Usage: make gen-fmc-keys SECRET=your_secret"; \
		exit 1; \
	fi
	@./scripts/create-fmc-key.sh $(SECRET)

## Generate key pair
gen-pair-keys:
	@./scripts/gen-pair-keys.sh


## Build Docker image
docker-build:
	@echo '  >> Building Docker image'
	@docker build -t vayload:$(VERSION) .
	@docker tag vayload:$(VERSION) vayload:latest
	@echo '  >> Docker image built: vayload:$(VERSION)'

## Push Docker image
docker-push: docker-build
	@echo '  >> Pushing Docker image'
	@docker push vayload:$(VERSION)
	@docker push vayload:latest

## Show version
version:
	@echo $(VERSION)

## Show build info
info:
	@echo 'Build Information:'
	@echo '  Version:    $(VERSION)'
	@echo '  Commit:     $(GIT_COMMIT)'
	@echo '  Build Time: $(BUILD_TIME)'
	@echo '  OS:         $(GOOS)'
	@echo '  Arch:       $(GOARCH)'
	@echo '  CGO:        $(CGO_ENABLED)'
