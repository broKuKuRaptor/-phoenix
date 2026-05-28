package accounts

import (
	"net/http"
	"phoenix/internal/common"
	"strings"

	"github.com/go-chi/chi/v5"
)

// getCurrencyRoutes handles GET /v1/currencies/{symbol}.
//
// Error discrimination:
//   - err != nil && routes == nil → 500 (unexpected internal error)
//   - err != nil && routes != nil → 404 (symbol not supported)
func getCurrencyRoutes(w http.ResponseWriter, r *http.Request, s AccountService) error {
	symbol := strings.ToUpper(chi.URLParam(r, "symbol"))
	if routes, err := s.GetCurrencyRoutes(r.Context(), symbol); err != nil {
		if routes == nil {
			return common.RespondError(w, http.StatusInternalServerError, err)
		}
		return common.RespondError(w, http.StatusNotFound, err)
	} else {
		return common.RespondJSON(w, http.StatusOK, routes)
	}
}

// getAccountCurrencies handles GET /v1/accounts/{owner_id}/currencies.
//
// Error discrimination:
//   - Invalid owner_id hex   → 400 Bad Request
//   - err != nil && currencies == nil → 500 (unexpected internal error)
//   - err != nil && currencies != nil → 404 (account not found)
func getAccountCurrencies(w http.ResponseWriter, r *http.Request, s AccountService) error {
	ownerId := chi.URLParam(r, "owner_id")
	if uid, err := common.ParseUID(ownerId); err != nil {
		return common.RespondError(w, http.StatusBadRequest, err)
	} else {
		if currencies, err := s.GetAccountCurrencies(r.Context(), uid); err != nil {
			if currencies == nil {
				return common.RespondError(w, http.StatusInternalServerError, err)
			}
			return common.RespondError(w, http.StatusNotFound, err)
		} else {
			return common.RespondJSON(w, http.StatusOK, currencies)
		}
	}
}

// Routes returns an HTTP handler for the accounts API.
//
// Endpoints:
//
//	GET /v1/currencies/{symbol}            – supported routes for a currency
//	GET /v1/accounts/{owner_id}/currencies  – currencies enabled on an account
//
// Mount under a service prefix, e.g. r.Mount("/api/accounts", svc.Routes()).
func (s AccountService) Routes() chi.Router {
	r := chi.NewRouter()
	r.Route("/v1", func(r chi.Router) {
		r.Get("/currencies/{symbol}", func(w http.ResponseWriter, r *http.Request) {
			getCurrencyRoutes(w, r, s)
		})
		r.Get("/accounts/{owner_id}/currencies", func(w http.ResponseWriter, r *http.Request) {
			getAccountCurrencies(w, r, s)
		})
	})
	return r
}
