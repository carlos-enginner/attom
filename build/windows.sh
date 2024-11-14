#!/bin/bash

# GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o ./tmp/push-relay.exe ./main.go
# go-msi make --msi attom.msi --version 1.0.0 -s ./wix/*.wsx
# goreleaser release --clean --skip=publish
goreleaser build --single-target 

