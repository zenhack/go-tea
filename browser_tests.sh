#!/usr/bin/env sh
#
# This script is used to run the wasm-only tests; execute it
# and then go to the displayed URL, and open up the browser
# console to see the results.

set -euo pipefail

GOOS=js GOARCH=wasm go test -c
go run internal/testrunner/main.go
