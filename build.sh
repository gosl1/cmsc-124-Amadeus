#!/usr/bin/env bash
set -e

# Create bin output directory
mkdir -p bin

# Go compiles a package from a single directory. Stage the scanner and parser
# sources together because they currently share package main.
build_dir="$(mktemp -d)"
trap 'rm -rf "$build_dir"' EXIT
cp cmd/scanner/*.go cmd/parser/*.go "$build_dir"/

go build -o bin/scanner "$build_dir"/*.go
