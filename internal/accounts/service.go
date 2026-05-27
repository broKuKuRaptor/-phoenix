package accounts

import (
	"context"

	"phoenix/internal/common"
)

// AccountService — сервисный слой для работы со счетами.
type AccountService struct {
	store Store
}

// NewService создает новый экземпляр AccountService.
func NewService(store Store) *AccountService {
	return &AccountService{store: store}
}

// Create создаёт новый счёт.
func (s *AccountService) Create(ctx context.Context, key AccountKey) (*Account, error) {
	return s.store.Create(ctx, key)
}

// ListAccounts возвращает список всех аккаунтов с пагинацией.
func (s *AccountService) ListAccounts(ctx context.Context, offset, limit int) ([]Account, error) {
	return s.store.List(ctx, offset, limit)
}

// GetAccountsByOwner возвращает список аккаунтов, принадлежащих указанному владельцу.
func (s *AccountService) GetAccountsByOwner(ctx context.Context, ownerId common.UID) ([]Account, error) {
	return s.store.ByOwner(ctx, ownerId)
}

// GetAccountByOwnerAndSymbol возвращает аккаунт по owner_id и символу.
func (s *AccountService) GetAccountByOwnerAndSymbol(ctx context.Context, ownerId common.UID, symbol string) (*Account, error) {
	return s.store.ByKey(ctx, ownerId, symbol)
}

// GetAccountByKey возвращает аккаунт по ключу (owner_id + symbol).
func (s *AccountService) GetAccountByKey(ctx context.Context, key AccountKey) (*Account, error) {
	return s.GetAccountByOwnerAndSymbol(ctx, key.OwnerId, key.Symbol)
}
