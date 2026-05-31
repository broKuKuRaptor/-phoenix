package accounts

import (
	"fmt"

	"phoenix/internal/common"
	"phoenix/internal/utils"

	"github.com/spf13/viper"
)

// AccountsConfig описывает конфигурацию для секции accounts.
type AccountsConfig struct {
	common.ServiceConfig `mapstructure:",squash"`
}

// GetAccountsConfig загружает и возвращает конфигурацию для секции accounts.
func GetAccountsConfig() (*AccountsConfig, error) {
	if cfgPath := viper.GetString("config"); cfgPath != "" {
		if err := utils.LoadConfigFromFile(cfgPath, "accounts"); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var accountsConfig AccountsConfig
	if err := viper.Unmarshal(&accountsConfig); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &accountsConfig, nil
}
