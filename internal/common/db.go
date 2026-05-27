package common

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

// DB — обёртка над sqlx.DB с поддержкой разных диалектов.
type DB struct {
	*sqlx.DB
	dialect string // "sqlite" or "postgres"
}

// ParseDSN разбирает databaseURL и возвращает driver + data source name.
//
// Поддерживаемые схемы:
//
//	sqlite:///path/to/db.sqlite
//	sqlite://:memory:
//	postgres://user:pass@host:port/dbname?sslmode=disable
func ParseDSN(databaseURL string) (driver, dsn string, err error) {
	u, err := url.Parse(databaseURL)
	if err != nil {
		return "", "", fmt.Errorf("ParseDSN: %w", err)
	}

	switch u.Scheme {
	case "sqlite", "sqlite3":
		driver = "sqlite"
		switch u.Host {
		case "", "localhost":
			dsn = u.Path
		default:
			dsn = u.Host + u.Path
		}
		if dsn == "" || dsn == "/" {
			return "", "", fmt.Errorf("ParseDSN: sqlite requires a file path or :memory:")
		}
		q := u.Query()
		for k, vs := range q {
			for _, v := range vs {
				dsn += "&" + k + "=" + v
			}
		}
		if strings.Contains(dsn, "?") {
			dsn = strings.Replace(dsn, "&", "?", 1)
		}
	case "postgres", "postgresql":
		driver = "postgres"
		dsn = databaseURL
	default:
		return "", "", fmt.Errorf("ParseDSN: unsupported scheme %q (use sqlite:// or postgres://)", u.Scheme)
	}
	return driver, dsn, nil
}

// Open открывает БД по URL, пингует и возвращает *DB.
func Open(databaseURL string) (*DB, error) {
	driver, dsn, err := ParseDSN(databaseURL)
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("Open: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("Open: ping: %w", err)
	}
	return &DB{DB: db, dialect: driver}, nil
}

// IsPostgres возвращает true для PostgreSQL.
func (db *DB) IsPostgres() bool { return db.dialect == "postgres" }

// IsUniqueViolation проверяет, является ли ошибка нарушением уникальности.
func IsUniqueViolation(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "UNIQUE constraint") ||
		strings.Contains(msg, "duplicate key")
}
