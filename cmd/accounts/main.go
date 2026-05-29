package main

import (
	"fmt"
	"phoenix/internal/config"
)


func main() {
	accountsCfg, err := config.GetAccountsConfig()
	if err != nil {
		panic(err)
	}

	fmt.Printf("БД: %s\n", accountsCfg.Database.Url)
	fmt.Printf("Адрес: %s:%d\n", accountsCfg.Address.Host, accountsCfg.Address.Port)
}

