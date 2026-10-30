#!/usr/bin/env bash
set -e

# Create bin output directory
mkdir -p bin

# Compile main.go into an executable inside bin/
go build -o bin/interpreter ./cmd/interpreter
