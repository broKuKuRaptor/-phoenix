package accounts

import "context"

// AccountService provides account-scoped operations such as querying
// supported currencies and their available (network, token) routes.
type AccountService struct {
	// TODO: Add fields for database connection, logger, etc.
}

// NewAccountService initializes and returns a new AccountService instance.
func NewAccountService() *AccountService {
	return &AccountService{}
}

// GetCurrencyRoutes returns the list of (network, token) pairs through which
// the service supports operations for the given currency symbol.
// Returns an empty slice if the symbol is unknown.
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
	return []CurrencyRoute{}, nil
}

// GetAccountCurrencies returns the currencies supported on the account
// linked to ownerId, each with its available CurrencyRoutes.
//
// When symbols are provided, the result is filtered to only those currencies.
// An empty symbols list returns all supported currencies.
func (s *AccountService) GetAccountCurrencies(_ctx context.Context, _ownerId string, symbols ...string) ([]AccountCurrency, error) {
	all := []AccountCurrency{
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
	}

	if len(symbols) == 0 {
		return all, nil
	}

	filtered := make([]AccountCurrency, 0, len(symbols))
	for _, symbol := range symbols {
		for _, ac := range all {
			if ac.Symbol == symbol {
				filtered = append(filtered, ac)
				break
			}
		}
	}
	return filtered, nil
}
