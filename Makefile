# =========================================================
# Vayload Build System
# =========================================================

.DEFAULT_GOAL := help

SHELL := /bin/bash
.ONESHELL:

# =========================================================
# Application
# =========================================================

APP_SERVER := vayload-server
APP_CLI := vayload

CMD_SERVER := ./cmd/server
CMD_CLI := ./cmd/cli

# =========================================================
# Directories
# =========================================================

BIN_DIR := ./bin
DIST_DIR := ./dist
RELEASE_DIR := ./release

# =========================================================
# Go / Platform
# =========================================================

GO := go

GOOS ?= $(shell $(GO) env GOOS)
GOARCH ?= $(shell $(GO) env GOARCH)

EXT :=
ifeq ($(GOOS),windows)
	EXT := .exe
endif

# =========================================================
# Build Metadata
# =========================================================

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# =========================================================
# Binary Names
# =========================================================

SERVER_BIN := $(BIN_DIR)/$(APP_SERVER)$(EXT)
CLI_BIN := $(BIN_DIR)/$(APP_CLI)$(EXT)

PACKAGE_NAME := vayload_$(VERSION)_$(GOOS)_$(GOARCH)

# =========================================================
# Build Configuration
# =========================================================

CGO_ENABLED ?= 1
BUILD_MODE ?= release

# Embedded metadata
LDFLAGS_VERSION := \
	-X 'main.Version=$(VERSION)' \
	-X 'main.BuildTime=$(BUILD_TIME)' \
	-X 'main.GitCommit=$(GIT_COMMIT)'

# Base optimizations
LDFLAGS := -s -w $(LDFLAGS_VERSION)
GCFLAGS :=
BUILD_TAGS :=

# =========================================================
# Platform Optimizations
# =========================================================

# Linux
ifeq ($(GOOS),linux)
	LDFLAGS += -extldflags '-static'
	GCFLAGS += -trimpath
	BUILD_TAGS += netgo osusergo
endif

# macOS
ifeq ($(GOOS),darwin)
	GCFLAGS += -trimpath

	ifeq ($(GOARCH),arm64)
		GCFLAGS += -N -l
	endif
endif

# Windows
ifeq ($(GOOS),windows)
	BUILD_TAGS += netgo
endif

# Debug mode
ifeq ($(BUILD_MODE),debug)
	LDFLAGS := -w $(LDFLAGS_VERSION)
	GCFLAGS += -N -l
endif

# =========================================================
# Go Build Command
# =========================================================

GO_BUILD = \
	CGO_ENABLED=$(CGO_ENABLED) \
	GOOS=$(GOOS) \
	GOARCH=$(GOARCH) \
	$(GO) build \
	-v \
	-trimpath \
	-buildvcs=false \
	-ldflags="$(LDFLAGS)" \
	-gcflags="$(GCFLAGS)" \
	-tags="$(BUILD_TAGS)"

# =========================================================
# Systemd Install
# =========================================================

SERVICE_NAME := $(APP_SERVER)
SERVICE_FILE := /etc/systemd/system/$(SERVICE_NAME).service
SERVICE_TEMPLATE := resources/$(SERVICE_NAME).service.tpl

SYSTEM_SERVER_BIN := /usr/local/bin/$(APP_SERVER)
SYSTEM_CLI_BIN := /usr/local/bin/$(APP_CLI)

# =========================================================
# Frontend Static Assets
# =========================================================

FRONTEND_BUILD_DIR := web/build
STATIC_DIR := static

# =========================================================
# Helpers
# =========================================================

define HEADER
	@echo ""
	@echo "================================================="
	@echo ">> $(1)"
	@echo "================================================="
endef

define BUILD_BINARY
	$(call HEADER,Building $(1))

	@mkdir -p $(BIN_DIR)

	$(GO_BUILD) \
		-o $(BIN_DIR)/$(1)$(EXT) \
		$(2)

	@echo ""
	@echo ">> Binary generated:"
	@ls -lh $(BIN_DIR)/$(1)$(EXT)
endef

# =========================================================
# PHONY
# =========================================================

.PHONY: \
	help \
	build \
	build-server \
	build-cli \
	build-client \
	package \
	install \
	dev \
	test \
	lint \
	clean \
	info \
	version \
	setup \
	keys \
	token

# =========================================================
# Help
# =========================================================

help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-22s\033[0m %s\n", $$1, $$2}'

# =========================================================
# Build
# =========================================================

build: build-server build-cli ## Build all binaries

build-server: ## Build server binary
	$(call BUILD_BINARY,$(APP_SERVER),$(CMD_SERVER))

build-cli: ## Build CLI binary
	$(call BUILD_BINARY,$(APP_CLI),$(CMD_CLI))

compress: ## Compress binaries with UPX (if available)
	@if command -v upx >/dev/null 2>&1; then \
		echo ">> Compressing binaries with UPX"; \
		upx $(SERVER_BIN); \
		upx $(CLI_BIN); \
	else \
		echo ">> UPX not installed"; \
	fi

# =========================================================
# Frontend
# =========================================================

build-client: ## Build frontend client
	$(call HEADER,Building frontend client)

	@if [ ! -d "web" ]; then \
		echo ">> No web directory found"; \
		exit 0; \
	fi

	@if ! command -v node >/dev/null 2>&1; then \
		echo ">> Node.js not installed"; \
		exit 0; \
	fi

	@if ! command -v pnpm >/dev/null 2>&1; then \
		echo ">> pnpm not installed"; \
		exit 0; \
	fi

	cd web

	if [ ! -d "node_modules" ]; then \
		echo ">> Installing dependencies"; \
		pnpm install --frozen-lockfile; \
	fi

	echo ">> Building frontend"
	pnpm run build

# =========================================================
# Package (GitHub Actions / Releases)
# =========================================================

package: build compress build-client ## Create release tar.gz package
	$(call HEADER,Packaging release)

	@mkdir -p $(DIST_DIR)/$(PACKAGE_NAME)
	@mkdir -p $(RELEASE_DIR)

	# -----------------------------------------------------
	# Binaries
	# -----------------------------------------------------

	@cp $(SERVER_BIN) $(DIST_DIR)/$(PACKAGE_NAME)/
	@cp $(CLI_BIN) $(DIST_DIR)/$(PACKAGE_NAME)/

	# -----------------------------------------------------
	# Frontend Static Assets
	# -----------------------------------------------------

	@if [ -d "$(FRONTEND_BUILD_DIR)" ]; then \
		echo ">> Including frontend static assets"; \
		mkdir -p $(DIST_DIR)/$(PACKAGE_NAME)/$(STATIC_DIR); \
		cp -R $(FRONTEND_BUILD_DIR)/* \
			$(DIST_DIR)/$(PACKAGE_NAME)/$(STATIC_DIR)/; \
	else \
		echo ">> No frontend build found"; \
	fi

	# -----------------------------------------------------
	# Metadata
	# -----------------------------------------------------

	@cp README.md $(DIST_DIR)/$(PACKAGE_NAME)/ 2>/dev/null || true
	@cp LICENSE $(DIST_DIR)/$(PACKAGE_NAME)/ 2>/dev/null || true

	@echo "$(VERSION)" > $(DIST_DIR)/$(PACKAGE_NAME)/VERSION
	@echo "$(GIT_COMMIT)" > $(DIST_DIR)/$(PACKAGE_NAME)/COMMIT

	# -----------------------------------------------------
	# Archive
	# -----------------------------------------------------

	@tar -czf \
		$(RELEASE_DIR)/$(PACKAGE_NAME).tar.gz \
		-C $(DIST_DIR) \
		$(PACKAGE_NAME)

	# -----------------------------------------------------
	# Checksums
	# -----------------------------------------------------

	@cd $(RELEASE_DIR) && \
	sha256sum $(PACKAGE_NAME).tar.gz > $(PACKAGE_NAME).sha256

	@echo ""
	@echo ">> Release generated:"
	@echo "   $(RELEASE_DIR)/$(PACKAGE_NAME).tar.gz"


# =========================================================
# Install From Source
# =========================================================

install: build compress build-client ## Install server locally with systemd
ifeq ($(GOOS),linux)
	$(call HEADER,Installing $(SERVICE_NAME))

	# -----------------------------------------------------
	# Install Server Binary
	# -----------------------------------------------------

	@sudo cp $(SERVER_BIN) $(SYSTEM_SERVER_BIN)
	@sudo chmod +x $(SYSTEM_SERVER_BIN)

	# -----------------------------------------------------
	# Install CLI Binary
	# -----------------------------------------------------

	@sudo cp $(CLI_BIN) $(SYSTEM_CLI_BIN)
	@sudo chmod +x $(SYSTEM_CLI_BIN)

	# -----------------------------------------------------
	# Install Frontend Assets
	# -----------------------------------------------------

	@if [ -d "$(FRONTEND_BUILD_DIR)" ]; then \
		echo ">> Installing frontend assets"; \
		sudo mkdir -p /var/lib/$(SERVICE_NAME)/$(STATIC_DIR); \
		sudo cp -R $(FRONTEND_BUILD_DIR)/* \
			/var/lib/$(SERVICE_NAME)/$(STATIC_DIR)/; \
	else \
		echo ">> No frontend assets found"; \
	fi

	# -----------------------------------------------------
	# Systemd Service
	# -----------------------------------------------------

	@if [ -f "$(SERVICE_FILE)" ]; then \
		echo ">> Service already exists"; \
	else \
		echo ">> Rendering systemd service"; \
		sed \
			-e "s|{{DESCRIPTION}}|vayload-server Daemon Server|g" \
			-e "s|{{DOCUMENTATION}}|https://vayload.dev/docs|g" \
			-e "s|{{USER}}|$(USER)|g" \
			-e "s|{{GROUP}}|$(USER)|g" \
			-e "s|{{WORKDIR}}|$(HOME)|g" \
			-e "s|{{EXEC_START}}|$(SYSTEM_SERVER_BIN)|g" \
			-e "s|{{ENV}}|production|g" \
			-e "s|{{ENV_FILE}}|/etc/$(SERVICE_NAME).env|g" \
			-e "s|{{NAME}}|$(SERVICE_NAME)|g" \
			$(SERVICE_TEMPLATE) | sudo tee $(SERVICE_FILE) > /dev/null; \
	fi

	# -----------------------------------------------------
	# Reload Systemd
	# -----------------------------------------------------

	@sudo systemctl daemon-reload
	@sudo systemctl enable $(SERVICE_NAME)
	@sudo systemctl restart $(SERVICE_NAME)

	@echo ""
	@echo ">> Service installed successfully"

else
	@echo ">> systemd install only supported on Linux"
	@exit 1
endif

# =========================================================
# Development
# =========================================================

dev: ## Run development server
	@command -v air >/dev/null 2>&1 || { \
		echo "Air not installed. Run: make setup"; \
		exit 1; \
	}

	air

setup: ## Install development dependencies
	$(call HEADER,Setup environment)

	$(GO) mod tidy

	@if ! command -v air >/dev/null 2>&1; then \
		echo ">> Installing Air"; \
		$(GO) install github.com/air-verse/air@latest; \
	fi

# =========================================================
# Quality
# =========================================================

test: ## Run tests
	$(call HEADER,Running tests)

	@$(GO) test -v -race -coverprofile=coverage.out ./...

lint: ## Run linter
	$(call HEADER,Running linter)

	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo ">> golangci-lint not installed"; \
	fi

# =========================================================
# Utilities
# =========================================================

clean: ## Clean build artifacts
	$(call HEADER,Cleaning)

	@rm -rf $(BIN_DIR)
	@rm -rf $(DIST_DIR)
	@rm -rf $(RELEASE_DIR)

	@rm -f coverage.out
	@rm -f coverage.html

version: ## Show version
	@echo $(VERSION)

info: ## Show build information
	@echo "Version:    $(VERSION)"
	@echo "Commit:     $(GIT_COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "GOOS:       $(GOOS)"
	@echo "GOARCH:     $(GOARCH)"
	@echo "CGO:        $(CGO_ENABLED)"

keys: ## Generate key pairs
	@./scripts/generate-key-pairs.sh

token: ## Generate token (SECRET=...)
	@if [ -z "$(SECRET)" ]; then \
		echo "Usage: make token SECRET=your_secret"; \
		exit 1; \
	fi

	@./scripts/generate-fmc-key.sh $(SECRET)
