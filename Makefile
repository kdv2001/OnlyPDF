SHELL := /bin/bash

GOLINT = $(shell which golangci-lint)
GOIMPORTS = $(shell which goimports)

lint-fix:
	@echo Running goimports: &&
	goimports -w -local "astra-jsc.com/" . &&
	${GOLINT} run -c .golangci.yml &&
	@echo Running staticcheck: &&
	staticcheck ./...

lint:
	@echo Running golangci-lint: &&
	${GOLINT} run -c .golangci.yml &&
	@echo Running staticcheck: &&
	staticcheck ./...

.ONESHELL:
