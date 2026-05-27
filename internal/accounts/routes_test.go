package accounts

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

// testService создаёт AccountService на SQLite :memory:.
func testService(t *testing.T) AccountService {
	t.Helper()
	store, err := Open("sqlite://:memory:")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	return *NewService(store)
}

// testRouter создаёт chi.Router с маршрутами аккаунтов.
func testRouter(t *testing.T) chi.Router {
	t.Helper()
	svc := testService(t)
	r := chi.NewRouter()
	r.Mount("/api/accounts", svc.Routes())
	return r
}

func TestCreateAccountHandler(t *testing.T) {
	t.Run("json", func(t *testing.T) {
		r := testRouter(t)
		body := `{"owner_id":"000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","symbol":"rub"}`
		req := httptest.NewRequest("POST", "/api/accounts/v1/", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("status = %d, want %d", w.Code, http.StatusCreated)
		}
		var acc Account
		json.NewDecoder(w.Body).Decode(&acc)
		if acc.Symbol != "RUB" {
			t.Errorf("Symbol = %q, want RUB", acc.Symbol)
		}
	})

	t.Run("duplicate", func(t *testing.T) {
		r := testRouter(t)
		body := `{"owner_id":"000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","symbol":"usd"}`
		// Первый запрос
		req := httptest.NewRequest("POST", "/api/accounts/v1/", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(httptest.NewRecorder(), req)
		// Дубликат
		req = httptest.NewRequest("POST", "/api/accounts/v1/", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusConflict {
			t.Errorf("status = %d, want %d", w.Code, http.StatusConflict)
		}
	})

	t.Run("bad_content_type", func(t *testing.T) {
		r := testRouter(t)
		req := httptest.NewRequest("POST", "/api/accounts/v1/", strings.NewReader("x"))
		req.Header.Set("Content-Type", "text/plain")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})
}

func TestListAccountsHandler(t *testing.T) {
	r := testRouter(t)

	// Создаём 3 счёта
	for _, sym := range []string{"RUB", "USD", "EUR"} {
		body := `{"owner_id":"000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","symbol":"` + sym + `"}`
		req := httptest.NewRequest("POST", "/api/accounts/v1/", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(httptest.NewRecorder(), req)
	}

	t.Run("all", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/accounts/v1/", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var accounts []Account
		json.NewDecoder(w.Body).Decode(&accounts)
		if len(accounts) != 3 {
			t.Errorf("len = %d, want 3", len(accounts))
		}
	})

	t.Run("limit", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/accounts/v1/?limit=1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		var accounts []Account
		json.NewDecoder(w.Body).Decode(&accounts)
		if len(accounts) != 1 {
			t.Errorf("len = %d, want 1", len(accounts))
		}
	})
}

func TestGetAccountHandler(t *testing.T) {
	r := testRouter(t)

	// Создаём счёт
	body := `{"owner_id":"000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","symbol":"rub"}`
	req := httptest.NewRequest("POST", "/api/accounts/v1/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(httptest.NewRecorder(), req)

	t.Run("found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/accounts/v1/000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f/RUB", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/accounts/v1/000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f/XXX", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
		}
	})
}

func TestListAccountsHandler_Empty(t *testing.T) {
	r := testRouter(t)
	req := httptest.NewRequest("GET", "/api/accounts/v1/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	body := strings.TrimSpace(w.Body.String())
	if body != "[]" {
		t.Errorf("body = %s, want []", body)
	}
}

func TestGetAccountsByOwnerHandler_NotFound(t *testing.T) {
	r := testRouter(t)

	req := httptest.NewRequest("GET",
		"/api/accounts/v1/000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestGetAccountsByOwnerHandler_BadUID(t *testing.T) {
	r := testRouter(t)

	req := httptest.NewRequest("GET", "/api/accounts/v1/xyz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGetAccountsByOwnerHandler(t *testing.T) {
	r := testRouter(t)

	// Создаём счета
	for _, sym := range []string{"RUB", "USD"} {
		body := `{"owner_id":"000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","symbol":"` + sym + `"}`
		req := httptest.NewRequest("POST", "/api/accounts/v1/", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(httptest.NewRecorder(), req)
	}

	req := httptest.NewRequest("GET", "/api/accounts/v1/000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	var accounts []Account
	json.NewDecoder(w.Body).Decode(&accounts)
	if len(accounts) != 2 {
		t.Errorf("len = %d, want 2", len(accounts))
	}
}
