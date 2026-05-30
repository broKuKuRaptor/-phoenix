package accounts

import (
	"errors"
	"fmt"
	"net/http"

	"phoenix/internal/types"
	"phoenix/internal/utils"

	"github.com/go-chi/chi/v5"
)

// Routes возвращает маршруты для сервиса учётных записей.
// Определяет конечные точки REST API с префиксом /v1:
//   GET    /         — список всех учётных записей
//   GET    /{id}     — получение учётной записи по owner_id
//   POST   /         — создание новой учётной записи
func (s *AccountsService) Routes() chi.Router {
	r := chi.NewRouter()
	r.Route("/v1", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			listAccountsHandler(w, r, s)
		})
		r.Get("/{owner_id}", func(w http.ResponseWriter, r *http.Request) {
			getAccountsHandler(w, r, s)
		})
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			createAccountHandler(w, r, s)
		})
	})
	return r
}

// listAccountsHandler обрабатывает GET-запрос на получение списка всех учётных записей.
func listAccountsHandler(w http.ResponseWriter, r *http.Request, ac *AccountsService) {
	accounts, err := ac.getAccounts()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err)
		return
	}
	utils.RespondJSON(w, http.StatusOK, accounts)
}

// getAccountsHandler обрабатывает GET-запрос на получение учётной записи по owner_id.
func getAccountsHandler(w http.ResponseWriter, r *http.Request, ac *AccountsService) {
	ownerIDStr := chi.URLParam(r, "owner_id")
	ownerID, err := types.ParseUID(ownerIDStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, errors.New("wrong owner_id parameter"))
		return
	}

	account, err := ac.getAccount(ownerID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err)
		return
	}
	if account == nil {
		utils.RespondError(w, http.StatusNotFound, fmt.Errorf("account %s not found", ownerID))
		return
	}
	utils.RespondJSON(w, http.StatusOK, account)
}

// createAccountHandler обрабатывает POST-запрос на создание новой учётной записи.
// В текущей реализации — заглушка.
func createAccountHandler(w http.ResponseWriter, r *http.Request, ac *AccountsService) {
	utils.RespondError(w, http.StatusNotImplemented, errors.New("create account: not yet implemented"))
}