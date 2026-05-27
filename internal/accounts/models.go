package accounts

import (
	"time"

	"phoenix/internal/common"

	"github.com/shopspring/decimal"
)

// AccountKey — уникальный идентификатор счета, состоящий из идентификатора владельца и валюты.
type AccountKey struct {
	OwnerId common.UID `json:"owner_id"`
	Symbol  string     `json:"symbol"`
}

// Account — модель данных финансовых счетов.
type Account struct {
	AccountKey
	CreatedAt time.Time       `json:"created_at"`
	Balance   decimal.Decimal `json:"balance"`
}
