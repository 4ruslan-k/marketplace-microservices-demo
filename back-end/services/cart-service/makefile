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
	@ENV_FILE_PATH=$(ENV_FILE_PATH) $(binDir)/air -c .air_debug.toml

.PHONY: lint
lint:
	@$(binDir)/golangci-lint run -v

.PHONY: install
install:
	@go mod tidy -v && \
	go mod download

.PHONY: migrate_init
migrate_init: # create migration tables
	@cd ./migrate && \
	BUNDEBUG=2 ENV_FILE_PATH=$(ENV_FILE_PATH) go run . db init

.PHONY: migrate_up
migrate_up:
	@cd ./migrate && \
	BUNDEBUG=2 ENV_FILE_PATH=$(ENV_FILE_PATH) go run . db migrate

.PHONY: migrate_down
migrate_down: # rollback the last migration group
	@cd ./migrate && \
	BUNDEBUG=2 ENV_FILE_PATH=$(ENV_FILE_PATH) go run . db rollback

.PHONY: migrate_create
migrate_create: # create up and down SQL migrations. Example: migrate_create_sql name=new_migration_name
	@cd ./migrate && \
	BUNDEBUG=2 ENV_FILE_PATH=$(ENV_FILE_PATH) go run . db create_sql $(name)

.PHONY: migrate_status
migrate_status: # print migrations status
	@cd ./migrate && \
	BUNDEBUG=2 ENV_FILE_PATH=$(ENV_FILE_PATH) go run . db status
