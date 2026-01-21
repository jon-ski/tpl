SHELL := /usr/bin/env bash

APP := $(notdir $(CURDIR))
BIN_DIR ?= bin
OUT ?= $(BIN_DIR)/$(APP)

GO ?= go
GOFLAGS ?=
CGO_ENABLED ?= 0

MAIN ?= .

MODULE := $(shell $(GO) list -m -f '{{.Path}} {{.Dir}}' all 2>/dev/null | awk '$$2=="$(CURDIR)" {print $$1; exit}')
VERSION_PKG ?= $(MODULE)/internal/version

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE    := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS ?= -s -w \
	-X '$(VERSION_PKG).Version=$(VERSION)' \
	-X '$(VERSION_PKG).Commit=$(COMMIT)' \
	-X '$(VERSION_PKG).Date=$(DATE)'

.PHONY: all tidy fmt vet test build run clean install

all: test build

tidy:
	$(GO) mod tidy

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

test:
	$(GO) test $(GOFLAGS) ./...

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

build: $(BIN_DIR)
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build $(GOFLAGS) -trimpath -ldflags "$(LDFLAGS)" -o "$(OUT)" $(MAIN)

run:
	CGO_ENABLED=$(CGO_ENABLED) $(GO) run $(GOFLAGS) -trimpath -ldflags "$(LDFLAGS)" $(MAIN)

clean:
	rm -rf "$(BIN_DIR)"

install:
	CGO_ENABLED=$(CGO_ENABLED) $(GO) install $(GOFLAGS) -trimpath -ldflags "$(LDFLAGS)"
