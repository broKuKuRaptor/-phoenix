package accounts

import (
	"phoenix/internal/common"
)

type AccountService struct {
	// Здесь могут быть зависимости, например, репозитории для доступа к данным
}

// NewService создает новый экземпляр AccountService
func NewService() *AccountService {
	return &AccountService{}
}

// ListAccounts возвращает список всех аккаунтов
func (s *AccountService) ListAccounts() ([]Account, error) {
	var accounts []Account = sampleAccounts
	return accounts, nil
}

// GetAccountsByOwner возвращает список аккаунтов, принадлежащих указанному владельцу
func (s *AccountService) GetAccountsByOwner(ownerId common.UID) ([]Account, error) {
	// var accounts []Account = ACCOUNTS
	// return accounts, nil
	return nil, common.NotFound("Accounts not found for owner_id: " + ownerId.String())
}

// GetAccountByOwnerAndSymbol возвращает аккаунт, принадлежащий указанному владельцу по заданной валюте
func (s *AccountService) GetAccountByOwnerAndSymbol(ownerId common.UID, symbol string) (*Account, error) {
	// var account *Account = &ACCOUNTS[0]
	// return account, nil
	return nil, common.NotFound("Account not found for owner_id: " + ownerId.String() + " and symbol: `" + symbol + "`")
}
