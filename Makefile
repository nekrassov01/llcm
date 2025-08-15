NAME := llcm

CMD_PATH := ./cmd/$(NAME)/
GOBIN ?= $(shell go env GOPATH)/bin

VERSION := $$(make -s show-version)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := "-s -w -X main.version=$(VERSION) -X main.revision=$(REVISION)"

HAS_LINT := $(shell command -v $(GOBIN)/golangci-lint 2> /dev/null)
HAS_VULN := $(shell command -v $(GOBIN)/govulncheck 2> /dev/null)
HAS_BUMP := $(shell command -v $(GOBIN)/gobump 2> /dev/null)

BIN_LINT := github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
BIN_VULN := golang.org/x/vuln/cmd/govulncheck@latest
BIN_BUMP := github.com/x-motemen/gobump/cmd/gobump@latest

export GO111MODULE=on

.PHONY: deps deps-lint deps-vuln deps-bump clean build check test cover bench lint vuln show-version check-git publish release

# -------
#  deps
# -------

deps: deps-lint deps-vuln deps-bump

deps-lint:
ifndef HAS_LINT
	go install $(BIN_LINT)
endif

deps-vuln:
ifndef HAS_VULN
	go install $(BIN_VULN)
endif

deps-bump:
ifndef HAS_BUMP
	go install $(BIN_BUMP)
endif

# --------
#  build
# --------

clean:
	go clean
	rm -f $(NAME) cover.out cover.html cpu.prof mem.prof $(NAME).test
	find . -maxdepth 1 -type f -regextype posix-extended -regex '\./$(NAME)[0-9]*\.html' -exec rm {} \;

build: clean
	go mod tidy
	go build -ldflags $(LDFLAGS) -o $(NAME) $(CMD_PATH)

# --------
#  check
# --------

check: test cover bench lint vuln

test:
	go test -race -cover -v ./... -coverprofile=cover.out -covermode=atomic

cover:
	go tool cover -html=cover.out -o cover.html

bench:
	go test -run=^$$ -bench=. -benchmem -count 10 -benchtime=10000x -cpuprofile=cpu.prof -memprofile=mem.prof

lint: deps-lint
	golangci-lint run ./... -v

vuln: deps-vuln
	$(GOBIN)/govulncheck -test ./...

# ----------
#  release
# ----------

show-version: deps-bump
	$(GOBIN)/gobump show -r $(CMD_PATH)

check-git:
ifneq ($(shell git status --porcelain),)
	$(error git workspace is dirty)
endif
ifneq ($(shell git rev-parse --abbrev-ref HEAD),main)
	$(error current branch is not main)
endif

publish: check-git deps-bump
	$(GOBIN)/gobump up -w $(CMD_PATH)
	git commit -am "bump up version to $(VERSION)"
	git push origin main

release: check-git
	git tag "v$(VERSION)"
	git push origin "refs/tags/v$(VERSION)"
