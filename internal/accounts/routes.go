package accounts

import (
	"net/http"
	types "phoenix/internal/common/types"
	"phoenix/internal/common"

	"github.com/go-chi/chi/v5"
)


// getCurrencies handles GET /v1/currencies.
//
// Error discrimination:
//   - err != nil && currencies == nil → 500 (unexpected internal error)
//   - err != nil && currencies != nil → 404 (no currencies configured)
func getCurrencies(w http.ResponseWriter, r *http.Request, s AccountService) error {
	if currencies, err := s.GetCurrencies(r.Context()); err != nil {
		if currencies == nil {
			return common.RespondError(w, http.StatusInternalServerError, err)
		}
		return common.RespondError(w, http.StatusNotFound, err)
	} else {
		return common.RespondJSON(w, http.StatusOK, currencies)
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
	if uid, err := types.ParseUID(ownerId); err != nil {
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
//	GET /v1/currencies                      – all supported currencies and their routes
//	GET /v1/accounts/{owner_id}/currencies  – currencies enabled on an account
//
func (s AccountService) Routes() chi.Router {
	r := chi.NewRouter()
	r.Route("/v1", func(r chi.Router) {
		r.Get("/currencies", func(w http.ResponseWriter, r *http.Request) {
			getCurrencies(w, r, s)	
		})
		r.Get("/accounts/{owner_id}/currencies", func(w http.ResponseWriter, r *http.Request) {
			getAccountCurrencies(w, r, s)
		})
	})
	return r
}

