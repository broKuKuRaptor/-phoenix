.PHONY: accounts currency all compose-up compose-down compose-test

BIN_DIR := bin
ACCOUNTS_BIN := $(BIN_DIR)/accounts
CURRENCY_BIN := $(BIN_DIR)/currency
CURRENCY_REPLICAS ?= 3

accounts:
	@mkdir -p $(BIN_DIR)
	go build -o $(ACCOUNTS_BIN) ./cmd/accounts

currency:
	@mkdir -p $(BIN_DIR)
	go build -o $(CURRENCY_BIN) ./cmd/currency

all: accounts currency

compose-up:
	docker compose up --build -d --scale currency=$(CURRENCY_REPLICAS)

compose-down:
	docker compose down -v --remove-orphans

compose-test:
	./scripts/test-docker-compose.sh
