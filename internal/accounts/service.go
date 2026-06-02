package accounts

import (
	"slices"
	"time"

	"phoenix/internal/broker"
	"phoenix/internal/common"
)

// AccountsService — сервис для работы с учётными записями.
type AccountsService struct {
	*common.BaseService
	broker *broker.Broker
}

// NewAccountsService создаёт и запускает экземпляр AccountsService.
func NewAccountsService(cfg *AccountsConfig) (*AccountsService, error) {
	as := &AccountsService{BaseService: common.NewBaseService()}

	if err := as.startAMQP(cfg); err != nil {
		return nil, err
	}

	as.StartPeriodicTask(as.updateCurrenciesSupportStatus, currenciesSupportPingInterval*time.Second)
	return as, nil
}

// GetCurrenciesSupportStatus возвращает список заполненных структур CurrenciesSupportStatus.
func (as *AccountsService) GetCurrenciesSupportStatus(symbols ...string) ([]CurrenciesSupportStatus, error) {
	result := []CurrenciesSupportStatus{}
	if len(symbols) > 0 {
		filtered := []CurrenciesSupportStatus{}
		for _, css := range result {
			if slices.Contains(symbols, css.Symbol) {
				filtered = append(filtered, css)
			}
		}
		return filtered, nil
	}
	return result, nil
}
