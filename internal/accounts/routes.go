package accounts

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"phoenix/internal/common"

	"github.com/go-chi/chi/v5"
)

// listAccountsHandler возвращает список всех аккаунтов.
// Поддерживает query-параметры: offset (default 0), limit (default 20).
func listAccountsHandler(w http.ResponseWriter, r *http.Request, s AccountService) {
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 || limit > 100 {
		limit = 20
	}
	accounts, err := s.ListAccounts(r.Context(), offset, limit)
	if err != nil {
		common.RespondError(w, err)
		return
	}
	if accounts == nil {
		accounts = []Account{}
	}
	common.RespondJSON(w, http.StatusOK, accounts)
}

// getAccountsHandler возвращает список аккаунтов для указанного владельца.
func getAccountsHandler(w http.ResponseWriter, r *http.Request, s AccountService) {
	ownerId, err := common.ParseUID(chi.URLParam(r, "owner_id"))
	if err != nil {
		common.RespondError(w, common.BadRequest("Invalid owner_id"))
		return
	}
	accounts, err := s.GetAccountsByOwner(r.Context(), ownerId)
	if err != nil {
		common.RespondError(w, err)
		return
	}
	if accounts == nil {
		accounts = []Account{}
	}
	common.RespondJSON(w, http.StatusOK, accounts)
}

// getAccountHandler обрабатывает GET-запрос на получение аккаунта по owner_id и symbol.
func getAccountHandler(w http.ResponseWriter, r *http.Request, s AccountService) {
	ownerId, err := common.ParseUID(chi.URLParam(r, "owner_id"))
	if err != nil {
		common.RespondError(w, common.BadRequest("Invalid owner_id"))
		return
	}
	symbol := strings.ToUpper(chi.URLParam(r, "symbol"))
	if account, err := s.GetAccountByOwnerAndSymbol(r.Context(), ownerId, symbol); err == nil {
		common.RespondJSON(w, http.StatusOK, account)
		return
	} else {
		common.RespondError(w, err)
	}

}

// createAccountHandler обрабатывает POST-запрос на создание аккаунта.
func createAccountHandler(w http.ResponseWriter, r *http.Request, s AccountService) {
	key, err := parseAccountKey(r)
	if err != nil {
		common.RespondError(w, err)
		return
	}
	account, err := s.Create(r.Context(), key)
	if err != nil {
		common.RespondError(w, err)
		return
	}
	common.RespondJSON(w, http.StatusCreated, account)
}

// parseAccountKey извлекает AccountKey из тела запроса.
// Поддерживает JSON (application/json) и форму (application/x-www-form-urlencoded).
func parseAccountKey(r *http.Request) (AccountKey, error) {
	ct := r.Header.Get("Content-Type")

	// JSON
	if strings.HasPrefix(ct, "application/json") {
		var key AccountKey
		if err := json.NewDecoder(r.Body).Decode(&key); err != nil {
			return key, common.BadRequest("Invalid JSON: " + err.Error())
		}
		key.Symbol = strings.ToUpper(key.Symbol)
		return key, nil
	}

	// Form
	if strings.HasPrefix(ct, "application/x-www-form-urlencoded") {
		if err := r.ParseForm(); err != nil {
			return AccountKey{}, common.BadRequest("Invalid form data")
		}
		ownerId, err := common.ParseUID(r.FormValue("owner_id"))
		if err != nil {
			return AccountKey{}, common.BadRequest("Invalid owner_id")
		}
		return AccountKey{
			OwnerId: ownerId,
			Symbol:  strings.ToUpper(r.FormValue("symbol")),
		}, nil
	}

	return AccountKey{}, common.BadRequest("Expected application/json or application/x-www-form-urlencoded")
}

// Routes возвращает маршруты для AccountService
func (s AccountService) Routes() chi.Router {
	r := chi.NewRouter()
	r.Route("/v1", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			listAccountsHandler(w, r, s)
		})
		r.Get("/{owner_id}", func(w http.ResponseWriter, r *http.Request) {
			getAccountsHandler(w, r, s)
		})
		r.Get("/{owner_id}/{symbol}", func(w http.ResponseWriter, r *http.Request) {
			getAccountHandler(w, r, s)
		})
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			createAccountHandler(w, r, s)
		})
	})
	return r
}
