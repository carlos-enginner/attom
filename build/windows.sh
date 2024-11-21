#!/bin/bash

# GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o ./tmp/atom.exe ./main.go
# go-msi make --msi attom.msi --version 1.0.0 -s ./wix/*.wsx
# goreleaser release --clean --skip=publish
# GOOS=windows GOARCH=amd64 goreleaser build -buildvcs --single-target 
# goreleaser build --snapshot --clean
# goreleaser release

# GITHUB_TOKEN=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
goreleaser release --clean

