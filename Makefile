.PHONY: accounts

BIN_DIR := bin
ACCOUNTS_BIN := $(BIN_DIR)/accounts

accounts:
	@mkdir -p $(BIN_DIR)
	go build -o $(ACCOUNTS_BIN) ./cmd/accounts
