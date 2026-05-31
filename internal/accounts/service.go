package accounts

import (
	"fmt"
	"slices"
	"time"

	"phoenix/internal/config"
	"phoenix/pkg/bases"
)

// AccountsService — сервис для работы с учётными записями.
type AccountsService struct {
	*bases.BaseService // Наследование логики базового сервиса
}

// CreateAndStart — создаёт и запускает экземпляр AccountsService.
func NewAccountsService(cfg *config.AccountsConfig) (*AccountsService, error) {
	as := &AccountsService{BaseService: bases.NewBaseService()}
	// Фоновые задачи
	as.StartPeriodicTask(as.updateCurrenciesSupportStatus, 3*time.Second)
	return as, nil
}

func (as *AccountsService) updateCurrenciesSupportStatus() {
	fmt.Print("Currencies support status updated!")
}

// GetCurrenciesSupportStatus — возвращает список заполненных структур CurrenciesSupportStatus. 
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
	} else {
		return result, nil
	}
}
