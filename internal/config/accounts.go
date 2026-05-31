package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// AccountsConfig описывает конфигурацию для секции accounts.
type AccountsConfig struct {
	// Database содержит настройки подключения к базе данных.
	Database struct {
		Url string `mapstructure:"url"`
	} `mapstructure:"database"`

	// Address содержит настройки сетевого адреса сервиса.
	Address struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"address"`

	// Параметры HTTP сервера
	HTTP struct {
		ReadTimeout  int `mapstructure:"read_timeout"`
		WriteTimeout int `mapstructure:"write_timeout"`
	} `mapstructure:"http"`
}

// init регистрирует флаги командной строки, связывает их с Viper,
// включает автоматическое чтение переменных окружения.
// Приоритет: флаг CLI -> переменная окружения -> значение по умолчанию.
func init() {
	pflag.StringP("config", "c", "", "путь к файлу конфигурации")

	pflag.String("database.url", "sql:///:memory:", "Database URL")
	pflag.String("address.host", "accounts", "Host of service address")
	pflag.Int("address.port", 9000, "Port of service address")

	pflag.Parse()

	// Связывание pflag с Viper
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		fmt.Printf("failed to bind flags: %v\n", err)
		os.Exit(1)
	}
	// Включает автоматическое чтение переменных окружения
	viper.AutomaticEnv()
	// Заменяет разделители для переменных окружения: дефисы на подчёркивания, точки на двойные подчёркивания
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "__"))
}

// GetAccountsConfig загружает и возвращает конфигурацию для секции accounts.
func GetAccountsConfig() (*AccountsConfig, error) {
	if cfgPath := viper.GetString("config"); cfgPath != "" {
		if err := loadConfigFromFile(cfgPath, "accounts"); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var accountsConfig AccountsConfig
	if err := viper.Unmarshal(&accountsConfig); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &accountsConfig, nil
}
