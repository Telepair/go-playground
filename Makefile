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
	@echo "  cellular-automaton-basic       Run the basic cellular automaton"
	@echo "  cellular-automaton-sierpinski  Run the sierpinski cellular automaton"
	@echo "  cellular-automaton-turing      Run the turing cellular automaton"
	@echo "  cellular-automaton-traffic     Run the traffic cellular automaton"
	@echo "  cellular-automaton-infinite    Run the infinite cellular automaton"
	@echo "  cellular-automaton-colorful    Run the colorful cellular automaton"
	@echo "  cellular-automaton-fixed       Run the fixed cellular automaton"
	@echo "  cellular-automaton-periodic    Run the periodic cellular automaton"
	@echo "  cellular-automaton-reflect     Run the reflect cellular automaton"
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
	cd cellular-automaton && go test -v -bench=. ./...

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
.PHONY: cellular-automaton-basic
cellular-automaton-basic: build-cellular-automaton
	@echo "Demo Cellular Automaton: Basic Rule 30..."
	./bin/cellular-automaton -rule 30 

.PHONY: cellular-automaton-sierpinski
cellular-automaton-sierpinski: build-cellular-automaton
	@echo "Demo Cellular Automaton: Sierpinski Triangle (Rule 90)..."
	./bin/cellular-automaton -rule 90 

.PHONY: cellular-automaton-turing
cellular-automaton-turing: build-cellular-automaton
	@echo "Demo Cellular Automaton: Turing Complete (Rule 110)..."
	./bin/cellular-automaton -rule 110 

.PHONY: cellular-automaton-traffic
cellular-automaton-traffic: build-cellular-automaton
	@echo "Demo Cellular Automaton: Traffic Flow (Rule 184)..."
	./bin/cellular-automaton -rule 184 -steps 30 -alive-char "🚗" -dead-char "░"

.PHONY: cellular-automaton-infinite
cellular-automaton-infinite: build-cellular-automaton
	@echo "Demo Cellular Automaton: Infinite Mode (Rule 30) - Press q to quit..."
	./bin/cellular-automaton -rule 30 -steps 0 -refresh 100ms

.PHONY: cellular-automaton-colorful
cellular-automaton-colorful: build-cellular-automaton
	@echo "Demo Cellular Automaton: Colorful Pattern (Rule 90)..."
	./bin/cellular-automaton -rule 90 -rows 35 -cols 60 -steps 30 -cellsize 1 -alive-color "#00FF00" -dead-color "#FF0000" -lang en

.PHONY: cellular-automaton-fixed
cellular-automaton-fixed: build-cellular-automaton
	@echo "Demo Cellular Automaton: Fixed Mode (Rule 30)..."
	./bin/cellular-automaton -rule 30 -steps 50 -boundary fixed

.PHONY: cellular-automaton-periodic
cellular-automaton-periodic: build-cellular-automaton
	@echo "Demo Cellular Automaton: Periodic Mode (Rule 30)..."
	./bin/cellular-automaton -rule 30 -steps 50 -boundary periodic

.PHONY: cellular-automaton-reflect
cellular-automaton-reflect: build-cellular-automaton
	@echo "Demo Cellular Automaton: Reflect Mode (Rule 30)..."
	./bin/cellular-automaton -rule 30 -steps 50 -boundary reflect