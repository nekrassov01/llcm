NAME := llcm
CMD_PATH := ./cmd/$(NAME)/
GOBIN ?= $(shell go env GOPATH)/bin
VERSION := $$(make -s show-version)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := "-s -w -X main.Version=$(VERSION) -X main.Revision=$(REVISION)"

HAS_LINT := $(shell command -v $(GOBIN)/golangci-lint 2> /dev/null)
HAS_VULNCHECK := $(shell command -v $(GOBIN)/govulncheck 2> /dev/null)
HAS_GOBUMP := $(shell command -v $(GOBIN)/gobump 2> /dev/null)

BIN_LINT := github.com/golangci/golangci-lint/cmd/golangci-lint@latest
BIN_GOVULNCHECK := golang.org/x/vuln/cmd/govulncheck@latest
BIN_GOBUMP := github.com/x-motemen/gobump/cmd/gobump@latest

export GO111MODULE=on

.PHONY: build
build: clean
	go mod tidy
	go build -ldflags $(LDFLAGS) -o $(NAME) $(CMD_PATH)

.PHONY: check
check: test cover bench vet golangci-lint govulncheck

.PHONY: deps
deps: deps-lint deps-govulncheck deps-gobump

.PHONY: deps-lint
deps-lint:
ifndef HAS_LINT
	go install $(BIN_LINT)
endif

.PHONY: deps-govulncheck
deps-govulncheck:
ifndef HAS_VULNCHECK
	go install $(BIN_GOVULNCHECK)
endif

.PHONY: deps-gobump
deps-gobump:
ifndef HAS_GOBUMP
	go install $(BIN_GOBUMP)
endif

.PHONY: test
test:
	go test -race -cover -v ./... -coverprofile=cover.out -covermode=atomic

.PHONY: cover
cover:
	go tool cover -html=cover.out -o cover.html

.PHONY: bench
bench:
	go test -run=^$$ -bench=. -benchmem -count 5 -cpuprofile=cpu.prof -memprofile=mem.prof

.PHONY: vet
vet:
	go vet

.PHONY: golangci-lint
golangci-lint: deps-lint
	golangci-lint run ./... -v

.PHONY: govulncheck
govulncheck: deps-govulncheck
	$(GOBIN)/govulncheck -test ./...

.PHONY: show-version
show-version: deps-gobump
	$(GOBIN)/gobump show -r $(CMD_PATH)

.PHONY: check-git
ifneq ($(shell git status --porcelain),)
	$(error git workspace is dirty)
endif
ifneq ($(shell git rev-parse --abbrev-ref HEAD),main)
	$(error current branch is not main)
endif

.PHONY: publish
publish: deps-gobump check-git
	$(GOBIN)/gobump up -w $(CMD_PATH)
	git commit -am "bump up version to $(VERSION)"
	git push origin main

.PHONY: release
release: check-git
	git tag "v$(VERSION)"
	git push origin "refs/tags/v$(VERSION)"

.PHONY: clean
clean:
	go clean
	rm -f $(NAME) cover.out cover.html cpu.prof mem.prof $(NAME).test
