package currency

import (
	"fmt"

	"phoenix/internal/common"
	"phoenix/internal/utils"

	"github.com/spf13/viper"
)

// CurrencyConfig описывает конфигурацию для секции currency.
type CurrencyConfig struct {
	common.ServiceConfig `mapstructure:",squash"`
}

// GetCurrencyConfig загружает и возвращает конфигурацию для секции currency.
func GetCurrencyConfig() (*CurrencyConfig, error) {
	if cfgPath := viper.GetString("config"); cfgPath != "" {
		if err := utils.LoadConfigFromFile(cfgPath, "currency"); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var currencyConfig CurrencyConfig
	if err := viper.Unmarshal(&currencyConfig); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &currencyConfig, nil
}
