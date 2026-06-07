# Build helpers — no changes to source files
TAG     := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "0.0.0")
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
USER    ?= $(shell git config user.github 2>/dev/null || echo "unknown")

LDFLAGS := -ldflags="-X 'github.com/bscott/pm-cli/internal/cli.Version=${TAG}+${USER}.${COMMIT}'"

.PHONY: build test

build:
	go build ${LDFLAGS} -o pm-cli ./cmd/pm-cli/

test:
	go test -count=1 -timeout 60s ./...

all: test build