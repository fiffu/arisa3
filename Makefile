default: help

BINARY := arisa3
COVERAGE_FILE := coverage.out

build: FORCE  ## build app
	go build -o "${BINARY}"


test: ## Run unit tests
	go test -failfast -covermode=count -coverprofile ${COVERAGE_FILE} ./...
	@go tool cover -func=${COVERAGE_FILE} | grep 'total' | sed -e 's/\t\+/ /g'
	@echo ✓ [make test-unit] Done

install-dev: install tooling  ## install for dev environments

install: FORCE  ## install build dependencies
	go get

tooling: FORCE  ## install development tooling
	@echo Installing golangci-lint
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.45.2

	@echo ''
	pip install pre-commit -q
	pre-commit install

	golangci-lint --version

help:  ## Show this help
	@echo 'usage: make [target] ...'
	@echo ''
	@echo 'targets:'
	@egrep '^(.+)\:\ .*##\ (.+)' ${MAKEFILE_LIST} | sed 's/:.*##/#/' | column -t -c 2 -s '#'

.PHONY: FORCE
FORCE:
