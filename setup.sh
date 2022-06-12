#!/bin/bash

cd ./notifier && go get ./... && cd ../executable && go get ./...
go build -o bin/exe ./cmd/cli/main.go
./bin/exe -i=1000 -url=http://url.com < test_files/test.txt
./executable/cmd/cli
go build -o bin/exe ./cmd/cli/main.go
./bin/exe -i=1000 -url=http://url.com < test_files/test.txt
