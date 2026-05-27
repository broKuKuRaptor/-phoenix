package accounts

import (
	"context"
	"errors"
	"testing"

	"phoenix/internal/common"
)

// testStore открывает SQLite :memory: и возвращает Store + cleanup.
func testStore(t *testing.T) Store {
	t.Helper()
	store, err := Open("sqlite://:memory:")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}

// testUID возвращает константный UID для тестов.
func testUID() common.UID {
	return common.UID{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	}
}

func TestStore_Create(t *testing.T) {
	store := testStore(t)
	ctx := context.Background()
	key := AccountKey{OwnerId: testUID(), Symbol: "RUB"}

	account, err := store.Create(ctx, key)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if account.Symbol != "RUB" {
		t.Errorf("Symbol = %q, want %q", account.Symbol, "RUB")
	}
	if !account.Balance.IsZero() {
		t.Errorf("Balance = %s, want 0", account.Balance)
	}
	if account.CreatedAt.IsZero() {
		t.Error("CreatedAt is zero")
	}
}

func TestStore_CreateDuplicate(t *testing.T) {
	store := testStore(t)
	ctx := context.Background()
	key := AccountKey{OwnerId: testUID(), Symbol: "USD"}

	if _, err := store.Create(ctx, key); err != nil {
		t.Fatalf("first Create: %v", err)
	}
	_, err := store.Create(ctx, key)
	if err == nil {
		t.Fatal("expected error on duplicate")
	}
	var appErr *common.AppError
	if !errors.As(err, &appErr) || appErr.Status != 409 {
		t.Errorf("expected 409 AlreadyExists, got %v", err)
	}
}

func TestStore_List(t *testing.T) {
	store := testStore(t)
	ctx := context.Background()

	t.Run("empty", func(t *testing.T) {
		accounts, err := store.List(ctx, 0, 100)
		if err != nil {
			t.Fatalf("List empty: %v", err)
		}
		if len(accounts) != 0 {
			t.Errorf("len = %d, want 0", len(accounts))
		}
	})

	// Создаём 3 счёта
	for _, sym := range []string{"RUB", "USD", "EUR"} {
		store.Create(ctx, AccountKey{OwnerId: testUID(), Symbol: sym})
	}

	t.Run("all", func(t *testing.T) {
		accounts, err := store.List(ctx, 0, 100)
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		if len(accounts) != 3 {
			t.Errorf("len = %d, want 3", len(accounts))
		}
	})

	t.Run("limit", func(t *testing.T) {
		accounts, err := store.List(ctx, 0, 2)
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		if len(accounts) != 2 {
			t.Errorf("len = %d, want 2", len(accounts))
		}
	})

	t.Run("offset", func(t *testing.T) {
		accounts, err := store.List(ctx, 1, 100)
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		if len(accounts) != 2 {
			t.Errorf("len = %d, want 2", len(accounts))
		}
	})
}

func TestStore_ByOwner(t *testing.T) {
	store := testStore(t)
	ctx := context.Background()
	uid := testUID()

	store.Create(ctx, AccountKey{OwnerId: uid, Symbol: "RUB"})
	store.Create(ctx, AccountKey{OwnerId: uid, Symbol: "USD"})

	accounts, err := store.ByOwner(ctx, uid)
	if err != nil {
		t.Fatalf("ByOwner: %v", err)
	}
	if len(accounts) != 2 {
		t.Errorf("len = %d, want 2", len(accounts))
	}
}

func TestStore_ByOwner_NotFound(t *testing.T) {
	store := testStore(t)
	ctx := context.Background()

	_, err := store.ByOwner(ctx, testUID())
	if err == nil {
		t.Fatal("expected error for owner with no accounts")
	}
	var appErr *common.AppError
	if !errors.As(err, &appErr) || appErr.Status != 404 {
		t.Errorf("expected 404 NotFound, got %v", err)
	}
}

func TestStore_ByKey(t *testing.T) {
	store := testStore(t)
	ctx := context.Background()
	uid := testUID()

	store.Create(ctx, AccountKey{OwnerId: uid, Symbol: "RUB"})

	t.Run("found", func(t *testing.T) {
		account, err := store.ByKey(ctx, uid, "RUB")
		if err != nil {
			t.Fatalf("ByKey: %v", err)
		}
		if account.Symbol != "RUB" {
			t.Errorf("Symbol = %q, want RUB", account.Symbol)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		_, err := store.ByKey(ctx, uid, "XXX")
		if err == nil {
			t.Fatal("expected error")
		}
		var appErr *common.AppError
		if !errors.As(err, &appErr) || appErr.Status != 404 {
			t.Errorf("expected 404 NotFound, got %v", err)
		}
	})
}
