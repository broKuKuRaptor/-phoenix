package accounts

import (
	"errors"
	"fmt"
	"net/http"

	"phoenix/internal/utils"

	"github.com/go-chi/chi/v5"
)

// AccountsRouterV1 возвращает маршруты CRUD учётных записей.
// Определяет конечные точки REST API с префиксом /v1:
//
//	GET    /         — список всех учётных записей
//	GET    /{id}     — получение учётной записи по owner_id
//	POST   /         — создание новой учётной записи
func (s *AccountsService) AccountsRouterV1() chi.Router {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		utils.RespondError(w, http.StatusNotImplemented, errors.New("get accounts list: not yet implemented"))
	})
	r.Get("/{owner_id}", func(w http.ResponseWriter, r *http.Request) {
		utils.RespondError(w, http.StatusNotImplemented, errors.New("get account info: not yet implemented"))
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		utils.RespondError(w, http.StatusNotImplemented, errors.New("create account: not yet implemented"))
	})

	return r
}

// CurrenciesRouterV1 возвращает маршруты получения информация о типе и статусе поддержки сервисом валют.
// Определяет конечные точки REST API с префиксом /v1:
//
//	GET    /             — информация по всем валютам
//	GET    /{symbol}     — информация по заданной валюте
func (s *AccountsService) CurrenciesRouterV1() chi.Router {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if result, err := s.GetCurrenciesSupportStatus(); err == nil {
			utils.RespondJSON(w, http.StatusOK, result)
		} else {
			utils.RespondError(w, http.StatusInternalServerError, err)
		}
	})

	r.Get("/{symbol}", func(w http.ResponseWriter, r *http.Request) {
		symbol := chi.URLParam(r, "symbol")
		if result, err := s.GetCurrenciesSupportStatus(symbol); err == nil {
			if len(result) > 0 {
				utils.RespondJSON(w, http.StatusOK, result[0])
			} else {
				utils.RespondError(w, http.StatusNotFound, fmt.Errorf("information for currency %s not found", symbol))
			}
		} else {
			utils.RespondError(w, http.StatusInternalServerError, err)
		}
	})

	return r
}
