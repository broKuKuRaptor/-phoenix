package main

import (
	"fmt"
	"log"

	"phoenix/internal/config"
)

// AccountsConfig mirrors the "accounts" section of config.yaml.
// koanf tags map struct fields to YAML keys.
type AccountsConfig struct {
	Database struct {
		URL string `koanf:"url"`
	} `koanf:"database"`
	Address struct {
		Host string `koanf:"host"`
		Port int    `koanf:"port"`
	} `koanf:"address"`
}

func main() {
	// Load config, scoped to the "accounts" service section.
	// CLI overrides like --address.port=9090 are applied automatically.
	cfg, err := config.Load("accounts")
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	var ac AccountsConfig
	if err := cfg.Unmarshal("", &ac); err != nil {
		log.Fatalf("unmarshal: %v", err)
	}

	fmt.Printf("Database URL: %s\n", ac.Database.URL)
	fmt.Printf("Address:      %s:%d\n", ac.Address.Host, ac.Address.Port)
}
