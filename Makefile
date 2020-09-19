BASEPATH = $(shell pwd)
export PATH := $(BASEPATH)/bin:$(PATH)

# Basic go commands
GOCMD      = go
GOBUILD    = $(GOCMD) build
GOINSTALL  = $(GOCMD) install
GORUN      = $(GOCMD) run
GOTEST     = $(GOCMD) test
GOGET      = $(GOCMD) get
GOFMT      = $(GOCMD) fmt
GOGENERATE = $(GOCMD) generate
GOTYPE     = $(GOCMD) type

BUILD_DIR = $(BASEPATH)
COVERAGE_DIR  = $(BUILD_DIR)/coverage
SUBCOV_DIR    = $(COVERAGE_DIR)/packages

BINARY = payments

DOCKERCOMPOSE = docker-compose

# all src packages without vendor and generated code
PKGS = $(shell go list ./... | grep -v /vendor)

# Colors
GREEN_COLOR   = \033[0;32m
PURPLE_COLOR  = \033[0;35m
DEFAULT_COLOR = \033[m

all: fmt build generate test lint

help:
	@echo 'Usage: make <TARGETS> ... <OPTIONS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@echo '    help               Show this help screen.'
	@echo '    test               Run unit tests.'
	@echo '    lint               Run all linters including vet and gosec and others'
	@echo '    coverage           Report code tests coverage.'
	@echo '    fmt                Run gofmt on package sources.'
	@echo '    build              Compile packages and dependencies.'
	@echo '    version            Print Go version.'
	@echo '    version            Print Go version.'

	@echo '    start              Start local env.'
	@echo '    stop               Stop local env.'
	@echo '    generate           Generate mock.'
	@echo ''
	@echo 'Targets run by default are: fmt build generate test lint.'
	@echo ''

coverage:
	@echo " [$(GREEN_COLOR)[coverage]$(DEFAULT_COLOR)]"
	@-mkdir -p $(SUBCOV_DIR)/
	@for package in $(PKGS); do $(GOTEST) -covermode=count -coverprofile $(SUBCOV_DIR)/`basename "$$package"`.cov "$$package"; done
	@echo 'mode: count' > $(COVERAGE_DIR)/coverage.cov ;
	@tail -q -n +2 $(SUBCOV_DIR)/*.cov >> $(COVERAGE_DIR)/coverage.cov ;
	@go tool cover -func=$(COVERAGE_DIR)/coverage.cov ;
	@if [ $(html) ]; then go tool cover -html=$(COVERAGE_DIR)/coverage.cov -o coverage.html ; fi
	@rm -rf $(COVERAGE_DIR);

lint:
	@echo " [$(GREEN_COLOR)lint$(DEFAULT_COLOR)]"
	@$(GORUN) ./vendor/github.com/golangci/golangci-lint/cmd/golangci-lint/main.go run \
	--no-config --enable=gosec --enable=gocyclo --enable=nakedret \
	--enable=bodyclose --enable=golint --enable=rowserrcheck --enable=stylecheck --enable=interfacer \
	--enable=unconvert --enable=goconst --enable=maligned --enable=depguard \
	--enable=unparam --enable=dogsled --enable=nakedret --enable=scopelint --enable=gocritic \
	--enable=whitespace --enable=goprintffuncname --enable=prealloc --enable=gofmt --enable=goimports \
	--enable=megacheck ./...

test:
	@echo " $(GREEN_COLOR)[test]$(DEFAULT_COLOR)"
	@$(GOTEST) -race -count=1 $(PKGS)

fmt:
	@echo " $(GREEN_COLOR)[format]$(DEFAULT_COLOR)"
	@$(GOFMT) $(PKGS)

build:
	@echo " $(GREEN_COLOR)[build]$(DEFAULT_COLOR)"
	@$(GOBUILD) --tags static -o $(BINARY)

version:
	@echo " $(GREEN_COLOR)[version]$(DEFAULT_COLOR)"
	@$(GOCMD) version

generate:
	@mkdir -p ./bin
ifeq ("$(wildcard ./bin/mockery)","")
	@echo " $(PURPLE_COLOR)[build mockery]$(DEFAULT_COLOR)"
	@$(GOBUILD) -o ./bin/mockery ./vendor/github.com/vektra/mockery/cmd/mockery/
endif
	@echo " $(GREEN_COLOR)[generate]$(DEFAULT_COLOR)"
	@$(GOGENERATE) $(PKGS)

update-deps:
	@echo " $(GREEN_COLOR)[update all deps]$(DEFAULT_COLOR)"
	@$(GOGET) -u
	@$(GOCMD) mod tidy
	@$(GOCMD) mod vendor

start:
	@echo " $(GREEN_COLOR)[start local env]$(DEFAULT_COLOR)"
	@${DOCKERCOMPOSE} -f ./local.yml build
	@${DOCKERCOMPOSE} -f ./local.yml up -d

stop:
	@echo " $(GREEN_COLOR)[stop local env]$(DEFAULT_COLOR)"
	@${DOCKERCOMPOSE} -f local.yml down

restart: stop start

