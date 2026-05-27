package accounts

import (
	"time"

	"phoenix/internal/common"

	"github.com/shopspring/decimal"
)

// AccountKey — уникальный идентификатор счета, состоящий из идентификатора владельца и валюты.
type AccountKey struct {
	OwnerId common.UID `db:"owner_id" json:"owner_id"`
	Symbol  string     `db:"symbol"   json:"symbol"`
}

func (a AccountKey) String() string {
	return a.OwnerId.String() + "/" + a.Symbol
}

// Account — модель данных финансовых счетов.
type Account struct {
	AccountKey
	CreatedAt time.Time       `db:"created_at" json:"created_at"`
	Balance   decimal.Decimal `db:"balance"    json:"balance"`
}
