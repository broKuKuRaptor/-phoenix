// Package accounts manages user accounts and their supported currencies.
//
// # Domain model
//
// Symbol — currency code (e.g. "USDT", "ETH", "TRX").
// Coin  — blockchain network identified by its native coin (e.g. "ETH" for Ethereum, "TRX" for Tron).
// Token — contract token on that network; equals Coin for native coins.
//
//	Example: {Coin: "ETH", Token: "USDT"} means "USDT as an ERC-20 token on Ethereum".
//
// CurrencyRoute defines one (network, token) pair through which operations
// with a currency are supported.
//
// AccountCurrency groups a currency (Symbol) with the list of CurrencyRoutes
// available for a specific user account.
package accounts

// CurrencyRoute describes a supported (network, token) pair for operations
// with a currency.
//
// Coin identifies the blockchain network; Token identifies the contract
// on that network. For native coins Coin and Token are equal.
//
// Examples:
//
//	{Coin: "ETH", Token: "ETH"}   – native ETH on Ethereum
//	{Coin: "ETH", Token: "USDT"}  – USDT as ERC-20 on Ethereum
//	{Coin: "TRX", Token: "USDT"}  – USDT as TRC-20 on Tron
type CurrencyRoute struct {
	Coin  string `json:"coin"`  // Blockchain network (native coin symbol)
	Token string `json:"token"` // Token contract on that network
}

// AccountCurrency defines which currency is supported on a user account
// and through which CurrencyRoutes operations can be performed.
type AccountCurrency struct {
	Symbol string          `json:"symbol"` // Currency code (e.g. "USDT", "ETH")
	Routes []CurrencyRoute `json:"routes"` // Supported (network, token) pairs
}
