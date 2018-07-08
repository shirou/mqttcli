GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)


.PHONY: all

GO_PROJECT = github.com/shirou/mqttcli
BUILD_DEST = build
COMMIT_HASH=`git rev-parse --short HEAD`
GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_SHA    = $(shell git rev-parse --short HEAD)
GIT_DIRTY  = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")

LDFLAGS += -w -s -extldflags -static

ifndef VERSION
	VERSION = DEV
endif

GOFLAGS := -ldflags "$(LDFLAGS)"

## Download dependencies and the run unit test and build the binary
all: install clean build

## Clean the dist directory
clean:
	@rm -rf $(BUILD_DEST)

## download dependencies to run this project
install:
	@which gox > /dev/null || go get github.com/mitchellh/gox
	@which dep > /dev/null || go get github.com/golang/dep/cmd/dep
	dep ensure -vendor-only

## Run for local development
start:
	DATA_DIRECTORY="$$PWD/data" \
	go run *.go

## Build the linux binary
build:
	@rm -rf $(BUILD_DEST)
	@mkdir -p $(BUILD_DEST) > /dev/null
	@CGO_ENABLED=0 \
	gox \
	-output "$(BUILD_DEST)/{{.Dir}}_{{.OS}}_{{.Arch}}" \
	$(GOFLAGS) \
	.

## Prints the version info about the project
info:
	 @echo "Version:           ${VERSION}"
	 @echo "Git Commit:        ${GIT_COMMIT}"
	 @echo "Git Tree State:    ${GIT_DIRTY}"

## Print the dependency graph and open in MAC
dependencygraph:
	dep status -dot | dot -T png | open -f -a /Applications/Preview.app

## Prints this help command
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET}: ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)