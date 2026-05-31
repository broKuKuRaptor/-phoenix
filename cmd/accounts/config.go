package main

import (
	"fmt"
	"os"
	"strings"

	"phoenix/internal/config"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// AccountsConfig описывает конфигурацию для секции accounts.
type AccountsConfig struct {
	config.CommonConfig `mapstructure:",squash"`

	// Address содержит настройки сетевого адреса сервиса.
	Address struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"address"`
}

// initConfig инициализирует конфигурацию для сервиса accounts.
func initConfig() {
	pflag.StringP("config", "c", "", "Config file path")

	pflag.String("database.url", "sql:///:memory:", "Database URL")
	pflag.String("address.host", "accounts", "Host of service address")
	pflag.Int("address.port", 9000, "Port of service address")
	pflag.Int("http.read_timeout", 10, "HTTP Read timeout")
	pflag.Int("http.write_timeout", 10, "HTTP Write timeout")
	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		fmt.Printf("failed to bind flags: %v\n", err)
		os.Exit(1)
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "__"))
}

// GetAccountsConfig загружает и возвращает конфигурацию для секции accounts.
func GetAccountsConfig() (*AccountsConfig, error) {
	if cfgPath := viper.GetString("config"); cfgPath != "" {
		if err := config.LoadConfigFromFile(cfgPath, "accounts"); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var accountsConfig AccountsConfig
	if err := viper.Unmarshal(&accountsConfig); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &accountsConfig, nil
}
