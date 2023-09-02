BIN_DIR := $(shell pwd)/bin
ENV_FILE_PATH := $(shell pwd)/.env
binDir := $(shell pwd)/bin
MOCKS_DESTINATION=internal/mocks

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

.PHONY: run-air
run:
	@ENV_FILE_PATH=$(ENV_FILE_PATH) $(binDir)/air -c .air_debug.toml

.PHONY: run
run:
	@ENV_FILE_PATH=$(ENV_FILE_PATH) go run $(shell pwd)/cmd/app


.PHONY: lint
lint:
	@$(binDir)/golangci-lint run -v

.PHONY: mocks
mocks:  
	${BIN_DIR}/mockgen -source=internal/services/user.go -destination=$(MOCKS_DESTINATION)/services/user.go
	${BIN_DIR}/mockgen -source=internal/repositories/authentication/interface.go -destination=$(MOCKS_DESTINATION)/repositories/authentication/interface.go
	${BIN_DIR}/mockgen -source=pkg/nats/interface.go -destination=$(MOCKS_DESTINATION)/nats/nats.go
	${BIN_DIR}/mockgen -source=internal/transport/http/controllers/user_controller.go -destination=$(MOCKS_DESTINATION)/sessions/session.go
	@echo "Mocks are generated at /$(MOCKS_DESTINATION)"
