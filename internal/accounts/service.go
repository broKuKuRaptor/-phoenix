package accounts

import (
	"phoenix/internal/config"
	"phoenix/internal/models"
	"phoenix/internal/types"
)

// AccountsService — сервис для работы с учётными записями.
// Содержит бизнес-логику и взаимодействует с конфигурацией.
type AccountsService struct {
	cfg *config.AccountsConfig
}

// NewService создаёт и возвращает новый экземпляр AccountsService.
func NewService(cfg *config.AccountsConfig) *AccountsService {
	return &AccountsService{cfg: cfg}
}

// getAccounts возвращает список всех учётных записей.
// В текущей реализации — заглушка.
func (ac *AccountsService) getAccounts() ([]models.Account, error) {
	return []models.Account{}, nil
}

// getAccount возвращает учётную запись по её идентификатору.
// В текущей реализации — заглушка.
func (ac *AccountsService) getAccount(_ownerID types.UID) (*models.Account, error) {
	return nil, nil
}
