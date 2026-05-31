package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	service "phoenix/internal/accounts"
	"phoenix/internal/config"
	"phoenix/pkg/apiserver"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// main — точка входа в сервис учётных записей.
func main() {
	// Загрузка конфигурации
	cfg, err := config.GetAccountsConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Запуск основного сервиса
	service, err := service.NewAccountsService(cfg)
	if err != nil {
		log.Fatalf("Failed to start accounts server: %v", err)
	}
	log.Printf("Accounts service started successful")

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(30 * time.Second))

	// Монтирование маршрутов сервиса
	router.Mount("/api/accounts", service.AccountsRouterV1())
	router.Mount("/api/currencies", service.CurrenciesRouterV1())

	// Запуск сервера
	apiServer, err := apiserver.CreateAndStart(*cfg, router)
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
