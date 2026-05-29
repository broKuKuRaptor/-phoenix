package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// AccountsConfig описывает секциию accounts
type AccountsConfig struct {
	Database struct {
		Url string `mapstructure:"url"`
	} `mapstructure:"database"`
	Address struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"address"`
}

func init() {
	pflag.StringP("config", "c", "", "Path to config file")

	pflag.String("database.url", "sql:///:memory:", "Database URL")
	pflag.String("address.host", "accounts", "Hostname for service")
	pflag.Int("address.port", 9000, "Port for service")

	pflag.Parse()

	// Связывание pflag с Viper
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		fmt.Printf("Error binding flags: %v\n", err)
		os.Exit(1)
	}
	// Включает автоматическое чтение ENV
	viper.AutomaticEnv()
	// Заменяет дефисы на подчеркивания и точки на двойные подчеркивания для переменных окружения
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "__"))
	// Приоритет: флаг CLI -> переменная окружения -> значение по умолчанию
}

// GetAccountsConfig загружает и возвращает конфигурацию для секции accounts.
func GetAccountsConfig() (*AccountsConfig, error) {

	if cfgPath := viper.GetString("config"); cfgPath != "" {
		if err := loadConfigFromFile(cfgPath, "accounts"); err != nil {
			return nil, fmt.Errorf("Failed to read config file: %v", err)
		}
	}

	// Если config-файл не указан — используются только флаги, ENV и значения по умолчанию.
	var accountsConfig AccountsConfig
	if err := viper.Unmarshal(&accountsConfig); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal config: %v", err)
	}

	return &accountsConfig, nil
}
