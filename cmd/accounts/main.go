package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"phoenix/internal/accounts"

	"github.com/go-chi/chi/v5"
)

func main() {
	dbURL := flag.String("db", "sqlite://:memory:", "Database URL (sqlite://path or postgres://...)")
	addr := flag.String("addr", ":8080", "Listen address")
	flag.Parse()

	store, err := accounts.Open(*dbURL)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
		os.Exit(1)
	}
	defer store.Close()

	router := chi.NewRouter()
	accountsService := accounts.NewService(store)
	router.Mount("/api/accounts", accountsService.Routes())

	log.Printf("Сервер запущен на http://localhost%s (db: %s)", *addr, *dbURL)
	if err := http.ListenAndServe(*addr, router); err != nil {
		log.Fatal(err)
	}
}
