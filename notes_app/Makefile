GO_PATH ?= $(shell go env GOPATH 2>/dev/null)
ifeq ($(GO_PATH),)
  GO_BIN_DIR := $(HOME)/go/bin
else
  GO_BIN_DIR := $(GO_PATH)/bin
endif
export PATH := $(GO_BIN_DIR):$(PATH)
GOOSE ?= goose
SWAG ?= $(GO_BIN_DIR)/swag
ENV_FILE ?= .env
DIR ?= db/migrations

migrate-up:
	set -a; . ./$(ENV_FILE); set +a; $(GOOSE) mysql "$$DBUSER:$$DBPASS@tcp($$DBADDRESS)/$$DBNAME?parseTime=true" -dir $(DIR) up

migrate-down:
	set -a; . ./$(ENV_FILE); set +a; $(GOOSE) mysql "$$DBUSER:$$DBPASS@tcp($$DBADDRESS)/$$DBNAME?parseTime=true" -dir $(DIR) down 1

migrate-status:
	set -a; . ./$(ENV_FILE); set +a; $(GOOSE) mysql "$$DBUSER:$$DBPASS@tcp($$DBADDRESS)/$$DBNAME?parseTime=true" -dir $(DIR) status

migrate-create:
	@read -p "name: " n; $(GOOSE) create $$n sql -dir $(DIR)

swagger:
	$(SWAG) init -g main.go -o docs

.PHONY: migrate-up migrate-down migrate-status migrate-create swagger
