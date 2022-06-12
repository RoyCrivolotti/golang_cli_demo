#!/bin/bash

echo "Downloading dependencies"
cd ./notifier && go get ./... && cd ../executable && go get ./...

echo "Building executable"
go build -o bin/exe ./cmd/cli/main.go

