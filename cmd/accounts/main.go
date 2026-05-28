package main

import (
	"fmt"
	"log"
	"net/http"

	"phoenix/internal/accounts"
	"phoenix/internal/config"

	"github.com/go-chi/chi/v5"
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

// main is the entry point for the accounts service. It loads configuration,
// initializes the AccountService with currencies from config, mounts HTTP
// routes, and starts the server.
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

	router := chi.NewRouter()
	currencies, err := cfg.Currencies()
	if err != nil {
		log.Fatalf("currencies: %v", err)
	}
	accountsService := accounts.NewAccountService(currencies)
	router.Mount("/api/accounts", accountsService.Routes())

	addr := fmt.Sprintf("%s:%d", ac.Address.Host, ac.Address.Port)
	log.Printf("Accounts server running at http://%s (db: %s)", addr, ac.Database.URL)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
