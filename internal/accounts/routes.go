package accounts

import (
	"net/http"
	"strings"

	"phoenix/internal/common"

	"github.com/go-chi/chi/v5"
)

func listAccountsHandler(w http.ResponseWriter, _ *http.Request, s AccountService) {
	if accounts, err := s.ListAccounts(); err == nil {
		common.RespondJSON(w, http.StatusOK, accounts)
	} else {
		common.RespondError(w, err)
	}
}

func getAccountsHandler(w http.ResponseWriter, r *http.Request, s AccountService) {
	ownerId, err := common.ParseUID(chi.URLParam(r, "owner_id"))
	if err != nil {
		common.RespondError(w, common.BadRequest("Invalid owner_id"))
		return
	}
	if accounts, err := s.GetAccountsByOwner(ownerId); err == nil {
		common.RespondJSON(w, http.StatusOK, accounts)
		return
	} else {
		common.RespondError(w, err)
	}
}

func getAccountHandler(w http.ResponseWriter, r *http.Request, s AccountService) {
	ownerId, err := common.ParseUID(chi.URLParam(r, "owner_id"))
	if err != nil {
		common.RespondError(w, common.BadRequest("Invalid owner_id"))
		return
	}
	symbol := strings.ToUpper(chi.URLParam(r, "symbol"))
	if account, err := s.GetAccountByOwnerAndSymbol(ownerId, symbol); err == nil {
		common.RespondJSON(w, http.StatusOK, account)
		return
	} else {
		common.RespondError(w, err)
	}

}

// Routes возвращает маршруты для AccountService
func (s AccountService) Routes() chi.Router {
	r := chi.NewRouter()
	r.Route("/v1", func(r chi.Router) {
		r.Get("/list", func(w http.ResponseWriter, r *http.Request) {
			listAccountsHandler(w, r, s)
		})
		r.Get("/{owner_id}", func(w http.ResponseWriter, r *http.Request) {
			getAccountsHandler(w, r, s)
		})
		r.Get("/{owner_id}/{symbol}", func(w http.ResponseWriter, r *http.Request) {
			getAccountHandler(w, r, s)
		})
	})
	return r
}
