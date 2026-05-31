package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"phoenix/internal/accounts"
	"phoenix/internal/apiserver"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
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

func main() {
	// Загрузка конфигурации
	cfg, err := accounts.GetAccountsConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	// Инициализация маршрутизатора
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(time.Duration(cfg.Http.WriteTimeout) * time.Second))

	// Запуск основного сервиса
	service, err := accounts.NewAccountsService(cfg)
	if err != nil {
		log.Fatalf("Failed to start accounts server: %v", err)
	}
	log.Printf("Accounts service started successfully")

	// Монтирование маршрутов сервиса
	router.Mount("/api/v1/accounts", service.AccountsRouterV1())
	router.Mount("/api/v1/currencies", service.CurrenciesRouterV1())

	// Запуск API сервера
	apiServer, err := apiserver.CreateAndStart(cfg.ServiceConfig, router)
	if err != nil {
		log.Fatalf("failed to start API server: %v", err)
	}
	log.Printf("API service started on %s", apiServer.Address())

	// Ожидаю системных сигналов завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Блокируемся здесь, пока не придет сигнал завершения

	// Создаю контекст с таймаутом 10 сек для завершения API сервера и основного сервиса
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Error while API server stopping: %v", err)
	}
	log.Print("API server stopped")

	if err := service.Shutdown(ctx); err != nil {
		log.Printf("Error while accounts server stopping: %v", err)
	}
	log.Print("Account server stopped")

}
