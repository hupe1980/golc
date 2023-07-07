PROJECTNAME=$(shell basename "$(PWD)")
VET_DIRS := $(shell find ./_examples -mindepth 1 -type d)

# Go related variables.
# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

.PHONY: setup
## setup: Setup installes dependencies
setup:
	@go mod tidy

.PHONY: lint
## test: Runs the linter
lint:
	golangci-lint run --color=always --sort-results ./...

.PHONY: test
## test: Runs go test with default values
test: 
	@go test -v -race -count=1 -coverprofile=coverage.out ./...

.PHONY: vet
## vet: Runs go vet for all exampels
vet:
	@for dir in $(VET_DIRS); do \
		echo "Running go vet in $$dir..."; \
		output=$$(go vet $$dir/*.go 2>&1); \
		if [ -n "$$output" ]; then \
            echo "$$output"; \
        fi \
	done

.PHONY: help
## help: Prints this help message
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo