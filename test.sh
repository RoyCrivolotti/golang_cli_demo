#!/bin/bash

echo "Running tests in both modules"
find . -name go.mod -execdir go test ./... \