# Binary output directory and name
BINARY_DIR := $(shell pwd)/bin
BINARY_NAME := go-playground

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
	@echo ""
	@echo "$(GREEN)Build:$(RESET)"
	@echo "  build                       Build the main go-playground binary"
	@echo "  build-all                   Build all standalone projects (legacy)"
	@echo "  build-cellular-automaton    Build the cellular automaton (standalone)"
	@echo "  build-conway-game-of-life   Build the conway game of life (standalone)"
	@echo "  build-mandelbrot-set        Build the mandelbrot set (standalone)"
	@echo "  build-random-walk           Build the random walk visualization (standalone)"
	@echo "  build-digital-rain          Build the digital rain (standalone)"
	@echo ""
	@echo "$(GREEN)Demos:$(RESET)" 
	@echo "  cellular-automaton       Run the cellular automaton"
	@echo "  conway-game-of-life      Run the conway game of life"
	@echo "  mandelbrot-set           Run the mandelbrot set fractal visualization"
	@echo "  random-walk              Run the random walk visualization"
	@echo "  digital-rain             Run the digital rain (Matrix effect)"
	@echo ""
	@echo "$(GREEN)Test:$(RESET)"
	@echo "  test                        Test all projects"
	@echo "  bench                       Run benchmarks"
	@echo "  clean                       Clean binary and cache"
	@echo ""
	@echo "$(GREEN)Quick Start:$(RESET)"
	@echo "  run                         Build and run go-playground (show help)"

# Build targets
.PHONY: build
build: tidy fmt vet lint osv
	@echo "  >  Building go-playground..."
	@mkdir -p bin
	go build -ldflags="-s -w" -o ./bin/$(BINARY_NAME) .
	@echo "  >  go-playground built successfully."

.PHONY: build-all
build-all: build-cellular-automaton build-conway-game-of-life build-mandelbrot-set build-random-walk build-digital-rain

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

.PHONY: build-random-walk
build-random-walk: tidy fmt vet lint osv 
	@echo "  >  Building random walk..."
	@mkdir -p bin
	go build -ldflags="-s -w" -o ./bin/random-walk ./random-walk
	@echo "  >  Random walk built successfully."

.PHONY: build-digital-rain
build-digital-rain: tidy fmt vet lint osv 
	@echo "  >  Building digital rain..."
	@mkdir -p bin
	go build -ldflags="-s -w" -o ./bin/digital-rain ./digital-rain
	@echo "  >  Digital rain built successfully."

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
	@echo "  >  Cleaning go-playground..."
	git clean -xdf
	@echo "  >  go-playground cleaned successfully."

# Quick start
.PHONY: run
run: build
	@echo "  >  Running go-playground..."
	./bin/$(BINARY_NAME) --help

# Cellular Automaton demos
.PHONY: cellular-automaton
cellular-automaton: build
	@echo "Demo Cellular Automaton: Basic Rule 30..."
	./bin/$(BINARY_NAME) cellular-automaton

# Conway Game of Life demos
.PHONY: conway-game-of-life
conway-game-of-life: build-conway-game-of-life
	@echo "Demo Conway Game of Life: Default Settings..."
	./bin/conway-game-of-life
	@echo "$(YELLOW)Note: This demo is still using standalone binary. Run 'make build' and use './bin/$(BINARY_NAME) conwayGameOfLife' when migrated to cobra.$(RESET)"

# Mandelbrot Set demos
.PHONY: mandelbrot-set
mandelbrot-set: build-mandelbrot-set
	@echo "Demo Mandelbrot Set: Fractal Visualization..."
	./bin/mandelbrot-set
	@echo "$(YELLOW)Note: This demo is still using standalone binary. Run 'make build' and use './bin/$(BINARY_NAME) mandelbrotSet' when migrated to cobra.$(RESET)"

# Random Walk demos
.PHONY: random-walk
random-walk: build-random-walk
	@echo "Demo Random Walk: Various random walk algorithms..."
	./bin/random-walk
	@echo "$(YELLOW)Note: This demo is still using standalone binary. Run 'make build' and use './bin/$(BINARY_NAME) randomWalk' when migrated to cobra.$(RESET)"

# Digital Rain demos
.PHONY: digital-rain
digital-rain: build-digital-rain
	@echo "Demo Digital Rain: Matrix-style falling characters..."
	./bin/digital-rain
	@echo "$(YELLOW)Note: This demo is still using standalone binary. Run 'make build' and use './bin/$(BINARY_NAME) digitalRain' when migrated to cobra.$(RESET)"
