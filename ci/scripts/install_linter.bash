#!/usr/bin/env bash

set -e

LINTER=$(which golangci-lint)
LINTER_VERSION=1.16.0

if [[ "$LINTER" == "" ]] ; then
    wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s v${LINTER_VERSION}
fi
