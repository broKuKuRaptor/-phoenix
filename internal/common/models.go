package common

import (
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
