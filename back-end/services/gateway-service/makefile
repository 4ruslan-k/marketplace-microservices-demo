BIN_DIR := $(shell pwd)/bin
ENV_FILE_PATH := $(shell pwd)/.env
binDir := $(shell pwd)/bin

.PHONY: test
test:
	@echo "Running tests..."
	@go test -cover -v ./...


.PHONY: install-tools
install-tools:
	@echo "Installing linter"
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(binDir)
	@echo "Installing air"
	@curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(binDir)

.PHONY: run
run:
	@ENV_FILE_PATH=$(ENV_FILE_PATH) $(binDir)/air

.PHONY: lint
lint:
	@$(binDir)/golangci-lint run -v

.PHONY: install
install:
	@go mod tidy -v && \
	go mod download