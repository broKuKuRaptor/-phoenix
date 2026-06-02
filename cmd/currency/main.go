package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"phoenix/internal/currency"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	pflag.StringP("config", "c", "", "Config file path")

	pflag.String("database.url", "sql:///:memory:", "Database URL")
	pflag.String("address.host", "currency", "Host of service address")
	pflag.Int("address.port", 9001, "Port of service address")
	pflag.Int("http.read_timeout", 10, "HTTP Read timeout")
	pflag.Int("http.write_timeout", 10, "HTTP Write timeout")
	pflag.String("amqp.url", "amqp://guest:guest@localhost:5672/", "AMQP broker URL")
	pflag.String("amqp.exchange", "phoenix.events", "AMQP topic exchange name")
	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		fmt.Printf("failed to bind flags: %v\n", err)
		os.Exit(1)
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "__"))
}

func main() {
	cfg, err := currency.GetCurrencyConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	service, err := currency.NewCurrencyService(cfg)
	if err != nil {
		log.Fatalf("Failed to start currency service: %v", err)
	}
	log.Printf("Currency service started (address %s:%d reserved for future HTTP)",
		cfg.Address.Host, cfg.Address.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := service.Shutdown(ctx); err != nil {
		log.Printf("Error while currency service stopping: %v", err)
	}
	log.Print("Currency service stopped")
}
