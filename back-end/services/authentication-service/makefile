BIN_DIR := $(shell pwd)/bin
ENV_FILE_PATH := $(shell pwd)/.env
binDir := $(shell pwd)/bin
MOCKS_DESTINATION=mocks

.PHONY: test
test:
	@echo "Running tests..."
	@go test -cover -v ./...


.PHONY: install
install:
	@go mod tidy -v && \
	go mod download

.PHONY: install-tools
install-tools:
	@echo "Installing linter"
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(binDir)
	@echo "Installing air"
	@curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(binDir)
	@echo "Installing mockgen"
	@GOBIN=$(BIN_DIR) go install github.com/golang/mock/mockgen

# .PHONY: run
# run:
# 	@ENV_FILE_PATH=$(ENV_FILE_PATH) $(binDir)/air -c .air_debug.toml

.PHONY: run
run:
	@ENV_FILE_PATH=$(ENV_FILE_PATH) go run $(binDir)/../cmd/app

.PHONY: lint
lint:
	@$(binDir)/golangci-lint run -v

.PHONY: mocks
# put the files with interfaces you'd like to mock in prerequisites
# wildcards are allowed
mocks: pkg/messaging/nats/nats.go
	@echo "Generating mocks..."
	@rm -rf $(MOCKS_DESTINATION)
	@for file in $^; do ${BIN_DIR}/mockgen -source=$$file -destination=$(MOCKS_DESTINATION)/$$file; done
	@echo "Mocks are generated at /$(MOCKS_DESTINATION)"
