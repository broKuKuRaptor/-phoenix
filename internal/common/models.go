package common

import (
	"database/sql/driver"
	"encoding/hex"
	"fmt"
)

// UID — уникальный идентификатор (32 байта), сериализуется в JSON как hex-строка.
type UID [32]byte

// MarshalJSON реализует json.Marshaler — [32]byte → hex-строка.
func (u UID) MarshalJSON() ([]byte, error) {
	dst := make([]byte, hex.EncodedLen(len(u)))
	hex.Encode(dst, u[:])
	return []byte(`"` + string(dst) + `"`), nil
}

// UnmarshalJSON реализует json.Unmarshaler — hex-строка → [32]byte.
func (u *UID) UnmarshalJSON(data []byte) error {
	s := string(data)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}
	if len(s) != hex.EncodedLen(len(u)) {
		return fmt.Errorf("invalid UID length: expected %d hex chars, got %d", hex.EncodedLen(len(u)), len(s))
	}
	_, err := hex.Decode(u[:], []byte(s))
	return err
}

// Scan реализует sql.Scanner — TEXT → UID.
func (u *UID) Scan(value any) error {
	if value == nil {
		return fmt.Errorf("cannot scan nil into UID")
	}
	var s string
	switch v := value.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return fmt.Errorf("cannot scan %T into UID", value)
	}
	return u.UnmarshalJSON([]byte(`"` + s + `"`))
}

// Value реализует driver.Valuer — UID → TEXT.
func (u UID) Value() (driver.Value, error) {
	return u.String(), nil
}

// ParseUID преобразует hex-строку в UID.
func ParseUID(s string) (UID, error) {
	var u UID
	if err := (&u).UnmarshalJSON([]byte(`"` + s + `"`)); err != nil {
		return u, err
	}
	return u, nil
}

// String возвращает UID в виде hex-строки (без кавычек).
func (u UID) String() string {
	return hex.EncodeToString(u[:])
}
