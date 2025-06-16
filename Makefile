# Binary output directory and name
BINARY_DIR := $(shell pwd)/bin

# Tools
GOIMPORTS := $(shell go env GOPATH)/bin/goimports
GOLANGCI_LINT := $(shell go env GOPATH)/bin/golangci-lint
OSV_SCANNER := $(shell go env GOPATH)/bin/osv-scanner

# Colors for output
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
RESET := \033[0m

.DEFAULT_GOAL := help

.PHONY: help
help:
	@echo "$(GREEN)Usage: make <target>$(RESET)"
	@echo "  build-cellular-automaton  Build the cellular automaton"
	@echo "  test-cellular-automaton   Test the cellular automaton"
	@echo ""
	@echo "$(GREEN)Cellular Automaton:$(RESET)" 
	@echo "  cellular-automaton       Run the cellular automaton"
	@echo ""
	@echo "$(GREEN)Others:$(RESET)"
	@echo "  ensure-tools              Ensure tools are installed"
	@echo "  clean                     Clean binary and cache"
	@echo "  help                      Show this help message"

# Build target
.PHONY: build-cellular-automaton
build-cellular-automaton: ensure-tools
	@echo "  >  Formatting cellular automaton..."
	cd cellular-automaton && go mod tidy
	cd cellular-automaton && go fmt ./...
	cd cellular-automaton && $(GOIMPORTS) -l -w .
	@echo "  >  Vetting cellular automaton..."
	cd cellular-automaton && go vet ./...
	cd cellular-automaton && $(GOLANGCI_LINT) run ./...
	@echo "  >  Scanning for vulnerabilities..."
	cd cellular-automaton && $(OSV_SCANNER) -r .
	@echo "  >  Building cellular automaton..."
	cd cellular-automaton && go build -o ../bin/cellular-automaton .
	@echo "  >  Cellular automaton built successfully."

.PHONY: test-cellular-automaton
test-cellular-automaton:
	@echo "  >  Testing cellular automaton..."
	cd cellular-automaton && go test -v -race ./...
	cd cellular-automaton && go test -v -bench=. -benchmem -run=^$
	cd cellular-automaton && go test -v -coverprofile=coverage.out
	cd cellular-automaton && go tool cover -func=coverage.out

.PHONY: ensure-tools
ensure-tools:
	@echo "  >  Ensuring tools..."
	@if [ ! -x "$(GOIMPORTS)" ]; then \
		echo "$(YELLOW)  >  Installing goimports...$(RESET)"; \
		go install golang.org/x/tools/cmd/goimports@latest; \
	fi
	@if [ ! -x "$(GOLANGCI_LINT)" ]; then \
		echo "$(YELLOW)  >  Installing golangci-lint...$(RESET)"; \
		go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6; \
	fi
	@if [ ! -x "$(OSV_SCANNER)" ]; then \
		echo "$(YELLOW)  >  Installing osv-scanner...$(RESET)"; \
		go install github.com/google/osv-scanner/v2/cmd/osv-scanner@v2; \
	fi
	@echo "  >  Tools ensured successfully."

.PHONY: clean
clean:
	@echo "  >  Cleaning cellular automaton..."
	git clean -xdf
	@echo "  >  Cellular automaton cleaned successfully."

# Cellular Automaton demos
.PHONY: cellular-automaton
cellular-automaton: build-cellular-automaton
	@echo "Demo Cellular Automaton: Basic Rule 30..."
	./bin/cellular-automaton