package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"go.yaml.in/yaml/v3"
)

// loadConfigFile загружает в viper конфигурацию из указанного файла.
// Если в корне есть ключ с именем sectionKey, загружает его содержимое, иначе загружает файл как есть.
func loadConfigFromFile(cfgPath, sectionKey string) error {
	var raw map[string]any
	if data, err := os.ReadFile(cfgPath); err == nil {
		if err := yaml.Unmarshal(data, &raw); err != nil {
			return fmt.Errorf("Failed to parse config file: %v", err)
		}
	} else {
		return fmt.Errorf("Failed to read config file: %v", err)
	}
	// Если в корне есть ключ с указанным именем — сохраняем его содержимое во временный файл и загружаем его через viper.
	if section, ok := raw[sectionKey]; ok {
		sectionData, err := yaml.Marshal(section)
		if err != nil {
			return fmt.Errorf("Failed to marshal accounts section: %v", err)
		}
		tmpFile, err := os.CreateTemp("", "*.yaml")
		if err != nil {
			return fmt.Errorf("Failed to create temp config file: %v", err)
		}
		if _, err := tmpFile.Write(sectionData); err != nil {
			tmpFile.Close()
			os.Remove(tmpFile.Name())
			return fmt.Errorf("Failed to write temp config file: %v", err)
		}
		tmpFile.Close()
		cfgPath = tmpFile.Name()
		defer os.Remove(tmpFile.Name())
	}
	//  Если ключа нет — загружаем файл как есть.
	viper.SetConfigFile(cfgPath)
	return viper.ReadInConfig()
}
