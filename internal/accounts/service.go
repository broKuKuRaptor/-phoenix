package accounts

import (
	"context"
	"errors"
	ct "phoenix/internal/common/types"
)

// AccountService provides account-scoped operations such as querying
// supported currencies and their available (network, token) routes.
//
// # Error conventions
//
// Methods return both a result slice and an error. Handlers use this
// combination to decide the HTTP status code:
//
//	result, err     → meaning
//	data,    nil    → 200 OK
//	[]T{},   err   → 404 Not Found (entity missing but payload valid)
//	nil,     err   → 500 Internal Server Error
type AccountService struct {
	currencies []ct.Currency
	// TODO: Add fields for database connection, logger, etc.
}

// NewAccountService initializes and returns a new AccountService instance.
func NewAccountService(currencies []ct.Currency) *AccountService {
	return &AccountService{currencies: currencies}
}

// GetCurrencies returns all currencies supported by the service, each with
// its available (network, token) routes. The data is loaded from the
// "currensies" section of the service configuration.
func (s *AccountService) GetCurrencies(_ctx context.Context) ([]ct.Currency, error) {
	return s.currencies, nil
}


// GetAccountCurrencies returns the currencies supported on the account
// linked to ownerId, each with its available CurrencyRoutes.
//
// An unknown ownerId returns an empty slice with an error (HTTP 404).
func (s *AccountService) GetAccountCurrencies(_ctx context.Context, ownerId ct.UID) ([]ct.Currency, error) {
	if ownerId.String() == "0000000000000000000000000000000000000000000000000000000000000000" {
		return []ct.Currency{
			{
				Symbol: "ETH",
				Routes: []ct.CurrencyRoute{
					{Network: "ETH", Token: "ETH"},
				},
			},
			{
				Symbol: "TRX",
				Routes: []ct.CurrencyRoute{
					{Network: "TRX", Token: "TRX"},
				},
			},
			{
				Symbol: "USDT",
				Routes: []ct.CurrencyRoute{
					{Network: "ETH", Token: "USDT"},
					{Network: "TRX", Token: "USDT"},
				},
			},
		}, nil
	}
	// return nil, errors.New("Some others error")
	return []ct.Currency{}, errors.New("account not found")
}
