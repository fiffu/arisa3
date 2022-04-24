default: help

build: FORCE  ## Build app
	go build -o arisa3
	@echo [make build] Done ✓


test:  ## Run unit tests
	go test -failfast -covermode=count -coverprofile coverage.out ./...
	@go tool cover -func=coverage.out | grep 'total' | sed -e 's/\t\+/ /g'
	@echo [make test] Done ✓

install-dev: install tooling  ## Install for dev environments

install: FORCE  ## Install build dependencies
	go get
	@echo [make install] Done ✓

tooling: FORCE  ## Install development tooling
	@echo Installing golangci-lint
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.45.2

	@echo ''
	pip install pre-commit -q
	pre-commit install

	golangci-lint --version
	@echo [make tooling] Done ✓

help:  ## Show this help
	@echo 'usage: make [target] ...'
	@echo ''
	@echo 'targets:'
	@egrep '^(.+)\:\ .*##\ (.+)' ${MAKEFILE_LIST} | sed 's/:.*##/#/' | column -t -c 2 -s '#'

.PHONY: FORCE
FORCE:
