package accounts

import (
	"context"
	"database/sql"
	"fmt"

	"phoenix/internal/common"

	"github.com/shopspring/decimal"
)

// Store defines data-access operations for accounts.
type Store interface {
	Create(ctx context.Context, key AccountKey) (*Account, error)
	List(ctx context.Context, offset, limit int) ([]Account, error)
	ByOwner(ctx context.Context, ownerId common.UID) ([]Account, error)
	ByKey(ctx context.Context, ownerId common.UID, symbol string) (*Account, error)
	Close() error
}

// SQLStore implements Store backed by *common.DB.
type SQLStore struct {
	*common.DB
}

// Open opens a database connection and runs accounts migration.
func Open(databaseURL string) (*SQLStore, error) {
	db, err := common.Open(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("accounts.Open: %w", err)
	}
	store := &SQLStore{DB: db}
	if err := store.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("accounts.Open: migrate: %w", err)
	}
	return store, nil
}

// migrate создаёт таблицу, если её ещё нет.
func (s *SQLStore) migrate() error {
	balanceType := "REAL"
	if s.IsPostgres() {
		balanceType = "NUMERIC"
	}
	query := `
	CREATE TABLE IF NOT EXISTS accounts (
		owner_id   TEXT NOT NULL,
		symbol     TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		balance    ` + balanceType + ` NOT NULL DEFAULT 0,
		PRIMARY KEY (owner_id, symbol)
	);`
	_, err := s.Exec(query)
	return err
}

// Create вставляет новый счёт. Возвращает AlreadyExists, если ключ занят.
func (s *SQLStore) Create(ctx context.Context, key AccountKey) (*Account, error) {
	account := &Account{
		AccountKey: key,
		Balance:    decimal.Zero,
	}
	err := s.GetContext(ctx, account, s.Rebind(`
		INSERT INTO accounts (owner_id, symbol, balance)
		VALUES (?, ?, ?)
		RETURNING owner_id, symbol, created_at, balance`),
		key.OwnerId, key.Symbol, account.Balance,
	)
	if err != nil {
		if common.IsUniqueViolation(err) {
			return nil, common.AlreadyExists("Account already exists: " + key.String())
		}
		return nil, fmt.Errorf("accounts.Create: %w", err)
	}
	return account, nil
}

// List возвращает список счетов с пагинацией.
func (s *SQLStore) List(ctx context.Context, offset, limit int) ([]Account, error) {
	var accounts []Account
	err := s.SelectContext(ctx, &accounts, s.Rebind(`
		SELECT owner_id, symbol, created_at, balance
		FROM accounts
		ORDER BY owner_id, symbol
		LIMIT ? OFFSET ?`), limit, offset)
	if err != nil {
		return nil, fmt.Errorf("accounts.List: %w", err)
	}
	return accounts, nil
}

// ByOwner возвращает счета конкретного владельца.
func (s *SQLStore) ByOwner(ctx context.Context, ownerId common.UID) ([]Account, error) {
	var accounts []Account
	err := s.SelectContext(ctx, &accounts, s.Rebind(`
		SELECT owner_id, symbol, created_at, balance
		FROM accounts
		WHERE owner_id = ?
		ORDER BY symbol`), ownerId)
	if err != nil {
		return nil, fmt.Errorf("accounts.ByOwner: %w", err)
	}
	if len(accounts) == 0 {
		return nil, common.NotFound("No accounts for owner: " + ownerId.String())
	}
	return accounts, nil
}

// ByKey возвращает конкретный счёт по ключу.
func (s *SQLStore) ByKey(ctx context.Context, ownerId common.UID, symbol string) (*Account, error) {
	var a Account
	err := s.GetContext(ctx, &a, s.Rebind(`
		SELECT owner_id, symbol, created_at, balance
		FROM accounts
		WHERE owner_id = ? AND symbol = ?`), ownerId, symbol)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, common.NotFound("Account not found: " + ownerId.String() + "/" + symbol)
		}
		return nil, fmt.Errorf("accounts.ByKey: %w", err)
	}
	return &a, nil
}

// Close закрывает соединение с БД.
func (s *SQLStore) Close() error {
	return s.DB.Close()
}
