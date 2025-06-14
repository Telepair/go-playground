# Cellular Automaton Makefile
# Author: liys87x@gmail.com
# Description: Quick build and run cellular automaton

# Help, Default target
.PHONY: help
help:
	@echo "Cellular Automaton:"
	@echo "  cellular-automaton-basic       Basic Rule 30 pattern"
	@echo "  cellular-automaton-sierpinski  Sierpinski Triangle (Rule 90)"
	@echo "  cellular-automaton-turing      Turing Complete (Rule 110)"
	@echo "  cellular-automaton-traffic     Traffic Flow (Rule 184)"
	@echo "  cellular-automaton-infinite    Infinite mode"
	@echo "  cellular-automaton-colorful    Colorful pattern"
	@echo "  cellular-automaton-fixed       Fixed mode"
	@echo "  cellular-automaton-periodic    Periodic mode"
	@echo "  cellular-automaton-reflect     Reflect mode"
	@echo ""
	@echo "Other:"
	@echo "  help          Show this help message"
	@echo "  build         Build the cellular automaton"
	@echo "  check         Run fmt, vet and tidy"
	@echo "  clean         Clean build artifacts and output files"

# Clean and check
.PHONY: clean, check, build
clean:
	@echo "Cleaning..."
	git clean -xdf
	@echo "Clean complete!"
check:
	@echo "Formatting Go code..."
	go fmt ./...
	@echo "Vetting Go code..."
	go vet ./...
	@echo "Tidying Go modules..."
	go mod tidy
	@echo "Installing dependencies..."
	go mod download
build:
	@echo "Building cellular automaton..."
	go build -o bin/cellular-automaton .

# Cellular Automaton demos
.PHONY: cellular-automaton-basic
cellular-automaton-basic: build
	@echo "Demo Cellular Automaton: Basic Rule 30..."
	./bin/cellular-automaton -rule 30 -size 60x30 -steps 30 -lang en

.PHONY: cellular-automaton-sierpinski
cellular-automaton-sierpinski: build
	@echo "Demo Cellular Automaton: Sierpinski Triangle (Rule 90)..."
	./bin/cellular-automaton -rule 90 -size 60x30 -steps 30 -cellsize 1 -lang en

.PHONY: cellular-automaton-turing
cellular-automaton-turing: build
	@echo "Demo Cellular Automaton: Turing Complete (Rule 110)..."
	./bin/cellular-automaton -rule 110 -size 60x30 -steps 30 -lang en

.PHONY: cellular-automaton-traffic
cellular-automaton-traffic: build
	@echo "Demo Cellular Automaton: Traffic Flow (Rule 184)..."
	./bin/cellular-automaton -rule 184 -size 60x30 -steps 30 -alive-char "🚗" -dead-char "░" -lang en

.PHONY: cellular-automaton-infinite
cellular-automaton-infinite: build
	@echo "Demo Cellular Automaton: Infinite Mode (Rule 30) - Press q to quit..."
	./bin/cellular-automaton -rule 30 -size 60x30 -steps 0 -lang en -refresh 0.1

.PHONY: cellular-automaton-colorful
cellular-automaton-colorful: build
	@echo "Demo Cellular Automaton: Colorful Pattern (Rule 90)..."
	./bin/cellular-automaton -rule 90 -size 60x30 -steps 30 -cellsize 1 -alive-color "#00FF00" -dead-color "#FF0000" -lang en

.PHONY: cellular-automaton-fixed
cellular-automaton-fixed: build
	@echo "Demo Cellular Automaton: Fixed Mode (Rule 30)..."
	./bin/cellular-automaton -rule 30 -size 60x30 -steps 50 -lang en -boundary fixed

.PHONY: cellular-automaton-periodic
cellular-automaton-periodic: build
	@echo "Demo Cellular Automaton: Periodic Mode (Rule 30)..."
	./bin/cellular-automaton -rule 30 -size 60x30 -steps 50 -lang en -boundary periodic

.PHONY: cellular-automaton-reflect
cellular-automaton-reflect: build
	@echo "Demo Cellular Automaton: Reflect Mode (Rule 30)..."
	./bin/cellular-automaton -rule 30 -size 60x30 -steps 50 -lang en -boundary reflect