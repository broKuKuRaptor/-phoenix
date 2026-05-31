package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// LoadConfigFromFile загружает в Viper конфигурацию из указанного YAML-файла.
// Если в корне YAML-документа есть ключ sectionKey (например, "accounts"),
// загружается только содержимое этой секции; иначе файл загружается целиком.
func LoadConfigFromFile(cfgPath, sectionKey string) error {
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Если в корне есть ключ sectionKey — извлекаем только эту секцию
	if section, ok := raw[sectionKey]; ok {

		// Если в корне есть ключ common, то объединяем его с секцией sectionKey
		if common, ok := raw["common"]; ok {
			section = mergeMaps(common, section)
		}

		sectionData, err := yaml.Marshal(section)
		if err != nil {
			return fmt.Errorf("failed to serialize section %s: %w", sectionKey, err)
		}

		tmpFile, err := os.CreateTemp("", "*.yaml")
		if err != nil {
			return fmt.Errorf("failed to create temp file: %w", err)
		}
		defer os.Remove(tmpFile.Name())

		if _, err := tmpFile.Write(sectionData); err != nil {
			tmpFile.Close()
			return fmt.Errorf("failed to write temp file: %w", err)
		}
		tmpFile.Close()

		cfgPath = tmpFile.Name()
	}

	viper.SetConfigFile(cfgPath)
	return viper.ReadInConfig()
}

// mergeMaps объединяет base и overlay; значения из overlay перекрывают base.
func mergeMaps(base, overlay any) map[string]any {
	baseMap, _ := base.(map[string]any)
	overlayMap, _ := overlay.(map[string]any)
	if baseMap == nil {
		baseMap = map[string]any{}
	}
	if overlayMap == nil {
		return baseMap
	}

	out := make(map[string]any, len(baseMap)+len(overlayMap))
	for k, v := range baseMap {
		out[k] = v
	}
	for k, v := range overlayMap {
		if baseVal, ok := out[k]; ok {
			if baseSub, ok := baseVal.(map[string]any); ok {
				if overlaySub, ok := v.(map[string]any); ok {
					out[k] = mergeMaps(baseSub, overlaySub)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}
