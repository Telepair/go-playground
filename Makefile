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
	@echo "  build-all                   Build all projects"
	@echo "  build-cellular-automaton    Build the cellular automaton"
	@echo "  build-conway-game-of-life   Build the conway game of life"
	@echo "  build-mandelbrot-set        Build the mandelbrot set"
	@echo "  test                        Test all projects"
	@echo "  bench                       Run benchmarks"
	@echo "  clean                       Clean binary and cache"
	@echo ""
	@echo "$(GREEN)Cellular Automaton:$(RESET)" 
	@echo "  cellular-automaton       Run the cellular automaton"
	@echo ""
	@echo "$(GREEN)Conway Game of Life:$(RESET)"
	@echo "  conway-game-of-life             Run the conway game of life"
	@echo ""
	@echo "$(GREEN)Mandelbrot Set:$(RESET)"
	@echo "  mandelbrot-set              Run the mandelbrot set fractal visualization"

# Build targets
.PHONY: build-all
build-all: build-cellular-automaton build-conway-game-of-life build-mandelbrot-set

.PHONY: build-cellular-automaton
build-cellular-automaton: tidy fmt vet lint osv 
	@echo "  >  Building cellular automaton..."
	@mkdir -p bin
	go build -ldflags="-s -w" -o ./bin/cellular-automaton ./cellular-automaton
	@echo "  >  Cellular automaton built successfully."

.PHONY: build-conway-game-of-life
build-conway-game-of-life: tidy fmt vet lint osv 
	@echo "  >  Building conway game of life..."
	@mkdir -p bin
	go build -ldflags="-s -w" -o ./bin/conway-game-of-life ./conway-game-of-life
	@echo "  >  Conway game of life built successfully."

.PHONY: build-mandelbrot-set
build-mandelbrot-set: tidy fmt vet lint osv 
	@echo "  >  Building mandelbrot set..."
	@mkdir -p bin
	go build -ldflags="-s -w" -o ./bin/mandelbrot-set ./mandelbrot-set
	@echo "  >  Mandelbrot set built successfully."


.PHONY: test
test: tidy fmt vet lint osv
	@echo "  >  Testing ..."
	go test -v -race ./...
	go test -v -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

.PHONY: bench
bench:
	go test -v -bench=. -benchmem -run=^$$ ./...

.PHONY: tidy fmt vet lint osv
tidy:
	go mod tidy

fmt: 
	go fmt ./...
	@if [ ! -x "$(GOIMPORTS)" ]; then \
		echo "$(YELLOW)  >  Installing goimports...$(RESET)"; \
		go install golang.org/x/tools/cmd/goimports@latest; \
	fi
	$(GOIMPORTS) -l -w .

vet: 
	go vet ./...

lint:
	@if [ ! -x "$(GOLANGCI_LINT)" ]; then \
		echo "$(YELLOW)  >  Installing golangci-lint...$(RESET)"; \
		go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6; \
	fi
	$(GOLANGCI_LINT) run ./...

osv: 
	@if [ ! -x "$(OSV_SCANNER)" ]; then \
		echo "$(YELLOW)  >  Installing osv-scanner...$(RESET)"; \
		go install github.com/google/osv-scanner/v2/cmd/osv-scanner@v2; \
	fi
	$(OSV_SCANNER) -r .

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

# Conway Game of Life demos
.PHONY: conway-game-of-life
conway-game-of-life: build-conway-game-of-life
	@echo "Demo Conway Game of Life: Default Settings..."
	./bin/conway-game-of-life

# Mandelbrot Set demos
.PHONY: mandelbrot-set
mandelbrot-set: build-mandelbrot-set
	@echo "Demo Mandelbrot Set: Fractal Visualization..."
	./bin/mandelbrot-set
