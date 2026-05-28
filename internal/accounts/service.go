package accounts

import (
	"context"
	"errors"
	"phoenix/internal/common"
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
	// TODO: Add fields for database connection, logger, etc.
}

// NewAccountService initializes and returns a new AccountService instance.
func NewAccountService() *AccountService {
	return &AccountService{}
}

// GetCurrencyRoutes returns the list of (network, token) pairs through which
// the service supports operations for the given currency symbol.
//
// An unknown symbol returns an empty slice with no error (HTTP 200 with []).
func (s *AccountService) GetCurrencyRoutes(_ctx context.Context, symbol string) ([]CurrencyRoute, error) {
	switch symbol {
	case "USDT":
		return []CurrencyRoute{
			{Coin: "ETH", Token: "USDT"},
			{Coin: "TRX", Token: "USDT"},
		}, nil
	case "ETH":
		return []CurrencyRoute{
			{Coin: "ETH", Token: "ETH"},
		}, nil
	case "TRX":
		return []CurrencyRoute{
			{Coin: "TRX", Token: "TRX"},
		}, nil
	}
	// return nil, errors.New("Some others error")
	return []CurrencyRoute{}, nil
}

// GetAccountCurrencies returns the currencies supported on the account
// linked to ownerId, each with its available CurrencyRoutes.
//
// An unknown ownerId returns an empty slice with an error (HTTP 404).
func (s *AccountService) GetAccountCurrencies(_ctx context.Context, ownerId common.UID) ([]AccountCurrency, error) {
	if ownerId.String() == "0000000000000000000000000000000000" {
		return []AccountCurrency{
			{
				Symbol: "ETH",
				Routes: []CurrencyRoute{
					{Coin: "ETH", Token: "ETH"},
				},
			},
			{
				Symbol: "TRX",
				Routes: []CurrencyRoute{
					{Coin: "TRX", Token: "TRX"},
				},
			},
			{
				Symbol: "USDT",
				Routes: []CurrencyRoute{
					{Coin: "ETH", Token: "USDT"},
					{Coin: "TRX", Token: "USDT"},
				},
			},
		}, nil
	}
	// return nil, errors.New("Some others error")
	return []AccountCurrency{}, errors.New("account not found")
}
