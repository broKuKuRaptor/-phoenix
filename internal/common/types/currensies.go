package common

// CurrencyRoute describes a supported (network, token) pair for operations
// with a currency.
//
// Network identifies the blockchain network; Token identifies the contract
// on that network. For native coins Network and Token are equal.
//
// Examples:
//
//	{Network: "ETH", Token: "ETH"}   – native ETH on Ethereum
//	{Network: "ETH", Token: "USDT"}  – USDT as ERC-20 on Ethereum
//	{Network: "TRX", Token: "USDT"}  – USDT as TRC-20 on Tron
type CurrencyRoute struct {
	Network string `json:"network"` // Blockchain network (native coin symbol)
	Token   string `json:"token"`   // Token contract on that network
}

// Currency describes a supported currency and through which CurrencyRoutes operations can be performed.
//
// Example:
//
//	{
//	  "symbol": "USD(T)",
//	  "routes": [
//	    {"network": "ETH", "token": "USDT"},
//	    {"network": "TRX", "token": "USDT"}
//	  ]
//	}		
// means "USD(T) supported as USDT on Ethereum and Tron".	
// Note: the symbol may include parentheses, e.g. "USD(T)", to indicate that the currency is a tokenized version of another currency (USD in this case).
//
// {
//   "symbol": "ETH",
//   "routes": [
//     {"network": "ETH", "token": "ETH"}	
//   ]
// }
// means "ETH supported as native coin on Ethereum".			
type Currency struct {
	Symbol string          `json:"symbol"` // Currency symbol (e.g. "USD(T)", "ETH")
	Routes []CurrencyRoute `json:"routes"` // Supported (network, token) pairs
}
