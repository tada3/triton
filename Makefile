PACKAGE  = github.com/tada3/triton
EXENAME  = triton
ADMIN_TOOL = att
DATE    ?= $(shell date +%FT%T%z)
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || \
			cat $(CURDIR)/.version 2> /dev/null || echo v0)
BIN      = $(GOPATH)/bin
PATH	:= $(BIN):$(PATH)
BASE     = $(GOPATH)/src/$(PACKAGE)

GO       = go
GOFMT    = go fmt
GOTEST   = go test -v
GOLIST   = go list ./... | grep -v "^$(PACKAGE)/custom"
TIMEOUT  = 15

V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

GIT_HEAD = $(shell git rev-parse HEAD)

.PHONY: all
all: fmt build



.PHONY: build
build: ; $(info $(M) building executable to bin/ …) @ ## Build cic-server executable
	$Q $(GO) build \
		-i -v \
		-tags release \
		-ldflags "-X $(PACKAGE).Version=$(VERSION) -X $(PACKAGE).BuildDate=$(DATE)" \
		-o bin/$(EXENAME) cmd/main.go

.PHONY: admin
admin: ; $(info $(M) building admin to bin/ …) @ ## Build cic-server executable
	$Q $(GO) build \
		-i -v \
		-tags release \
		-ldflags "-X $(PACKAGE).Version=$(VERSION) -X $(PACKAGE).BuildDate=$(DATE)" \
		-o bin/$(ADMIN_TOOL) cmd/admin.go


# Tools

.PHONY: lint
lint: ; $(info $(M) running golint…) @ ## Run golint
	$Q cd $(BASE) && $(GOLINT) $$($(GOLIST)) | grep -v 'be unexported' | grep -v '.g.go' || exit 0

.PHONY: strict-lint
strict-lint: ; $(info $(M) running golint…) @ ## Run golint
	$Q cd $(BASE) && $(GOLINT) $$($(GOLIST))

.PHONY: fmt
fmt: ; $(info $(M) running gofmt…) @ ## Run gofmt on all source files
	$Q cd $(BASE) && $(GOFMT) $$($(GOLIST))

.PHONY: test
test: ; $(info $(M) running go test...) @ ## Run go test
	$Q cd $(BASE) && GORACE="halt_on_error=1" $(GOTEST) -race $$($(GOLIST)) | grep -v "^$(PACKAGE)/integration"



# Misc

.PHONY: run
run: build ; @ ## Build and run a server
	$Q bin/$(EXENAME)

.PHONY: clean
clean: ; $(info $(M) cleaning…)	@ ## Cleanup everything
	@rm -rf bin
	@rm -rf test/tests.* test/coverage.*
	@rm -rf **/*.g.go

.PHONY: help
help:
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: version
version:
	@echo $(VERSION)

