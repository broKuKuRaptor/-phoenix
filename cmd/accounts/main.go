package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"phoenix/internal/accounts"
	"phoenix/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// main — точка входа в сервис учётных записей.
func main() {
	cfg, err := config.GetAccountsConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	router := chi.NewRouter()
	// Middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(30 * time.Second))

	// Монтирование маршрутов сервиса
	service := accounts.NewService(cfg)
	router.Mount("/api/accounts", service.Routes())

	// Запуск сервера
	addr := fmt.Sprintf("%s:%d", cfg.Address.Host, cfg.Address.Port)
	log.Printf("service started on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to start server: %v", err)
	}	
	log.Print("server stopped")
}
