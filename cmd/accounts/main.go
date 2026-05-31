package main

import (
	"fmt"
	"log"
)

func main() {
	initConfig()

	cfg, err := GetAccountsConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Println("Address: ", cfg.Address.Host, ":", cfg.Address.Port)
	fmt.Println("Database URL: ", cfg.Database.Url)
	fmt.Println("HTTP Read timeout: ", cfg.Http.ReadTimeout)
	fmt.Println("HTTP Write timeout: ", cfg.Http.WriteTimeout)
}
