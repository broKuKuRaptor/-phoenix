package broker

// Имена routing key для topic exchange phoenix.events (вариант A).
const (
	RoutingPingCurrenciesSupportStatus = "currencies.ping.support_status"
	RoutingPongCurrenciesSupportStatus = "currencies.pong.support_status"

	BindingCurrencyPing = "currencies.ping.*"
	BindingAccountsPong = "currencies.pong.*"
)

// Типы сообщений в JSON-поле type.
const (
	TypePingCurrenciesSupportStatus = "ping:CurrenciesSupportStatus"
	TypePongCurrenciesSupportStatus = "pong:CurrenciesSupportStatus"
)
