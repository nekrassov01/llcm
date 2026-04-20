NAME := llcm
CMD_PATH := ./cmd/llcm/

GOBIN ?= $(shell go env GOPATH)/bin

VERSION := $$(make version)
REVISION := $$(make revision)
LDFLAGS := "-s -w -X github.com/nekrassov01/llcm.version=$(VERSION) -X github.com/nekrassov01/llcm.revision=$(REVISION)"

HAS_LINT := $(shell command -v $(GOBIN)/golangci-lint 2> /dev/null)
HAS_VULN := $(shell command -v $(GOBIN)/govulncheck 2> /dev/null)
HAS_BUMP := $(shell command -v $(GOBIN)/gobump 2> /dev/null)

BIN_LINT := github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
BIN_VULN := golang.org/x/vuln/cmd/govulncheck@latest
BIN_BUMP := github.com/x-motemen/gobump/cmd/gobump@latest

export GO111MODULE=on

.PHONY: deps deps-lint deps-vuln deps-bump clean build check test cover bench lint vuln version revision check-git check-branch bump

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
	rm -f $(NAME) coverage.out coverage.html cpu.prof mem.prof $(NAME).test
	find . -maxdepth 1 -type f -regextype posix-extended -regex '\./$(NAME)[0-9]*\.html' -exec rm {} \;

build: clean
	go mod tidy
	go build -ldflags $(LDFLAGS) -o $(NAME) $(CMD_PATH)

# --------
#  check
# --------

check: test cover bench lint vuln

test:
	go test -race -cover -v -coverprofile coverage.out -covermode atomic ./...

cover:
	go tool cover -html coverage.out -o coverage.html

bench:
	go test -bench -benchmem -count 5 -benchtime 10000x -cpuprofile cpu.prof -memprofile mem.prof .

lint: deps-lint
	golangci-lint run --verbose ./...

vuln: deps-vuln
	govulncheck -test -show verbose ./...

# ----------
#  version
# ----------

version: deps-bump
	@echo $(shell gobump show -r .)

revision: deps-bump
	@echo $(shell git rev-parse --short HEAD)

check-git:
ifneq ($(shell git status --porcelain),)
	$(error git workspace is dirty)
endif

check-branch:
ifndef branch
	$(error branch is undefined)
endif

bump: check-branch deps-bump
	gobump up -w .
	git commit -am "bump up version to $(VERSION)"
	git push origin $(branch)
