PACKAGE  = github.com/tadaweb/gonsx
BINARYNAME = gonsx-cli
DATE    ?= $(shell date +%FT%T%z)
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || \
			cat $(CURDIR)/VERSION 2> /dev/null || echo v0)

GOPATH   = $(CURDIR)/.gopath~
BIN      = $(GOPATH)/bin
BASE     = $(GOPATH)/src/$(PACKAGE)
PKGS     = $(or $(PKG),$(shell cd $(BASE) && env GOPATH=$(GOPATH) $(GO) list ./... | grep -v "^$(PACKAGE)/vendor/"))
TESTPKGS = $(shell env GOPATH=$(GOPATH) $(GO) list -f '{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' $(PKGS))


GO      = go
GODOC   = godoc
GOFMT   = gofmt
GLIDE   = glide
TIMEOUT = 15
V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

.PHONY: all
all: fmt lint test vendor | $(BASE) ; $(info $(M) building executable…) @ ## Build program binary
	$Q cd $(BASE)/cli && $(GO) build -tags release -o ../$(BINARYNAME) ./*.go

$(BASE): ; $(info $(M) setting GOPATH…)
	@mkdir -p $(dir $@)
	@ln -sf $(CURDIR) $@

# Tools

GOLINT = $(BIN)/golint
$(BIN)/golint: | $(BASE) ; $(info $(M) building golint…)
	$Q go get github.com/golang/lint/golint

TEST_TARGETS := test-default
.PHONY: $(TEST_TARGETS) check test tests
$(TEST_TARGETS): NAME=$(MAKECMDGOALS:test-%=%)
$(TEST_TARGETS): test
check test tests: fmt lint vendor | $(BASE) ; $(info $(M) running $(NAME:%=% )tests…) @ ## Run tests
	$Q cd $(BASE) && $(GO) test -timeout $(TIMEOUT)s $(ARGS) $(TESTPKGS)

.PHONY: lint
lint: vendor | $(BASE) $(GOLINT) ; $(info $(M) running golint…) @ ## Run golint
	$Q cd $(BASE) && ret=0 && for pkg in $(PKGS); do \
		test -z "$$($(GOLINT) $$pkg | tee /dev/stderr)" || ret=1 ; \
	 done ; exit $$ret

.PHONY: fmt
fmt: ; $(info $(M) running gofmt…) @ ## Run gofmt on all source files
	@ret=0 && for d in $$($(GO) list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		$(GOFMT) -l -w $$d/*.go || ret=$$? ; \
	 done ; exit $$ret

.PHONY: clean
clean: ; $(info $(M) cleaning…)	@ ## Cleanup everything
	@rm -rf $(GOPATH)
	@rm -rf bin
	@rm -rf test/tests.* test/coverage.*
	@rm -rf gonsx-cli

.PHONY: help
help:
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: version
version:
	@echo $(VERSION)
