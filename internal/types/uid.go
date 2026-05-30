package types

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// UID — уникальный идентификатор, представленный 32 байтами (256 бит).
type UID [32]byte

// NewUID генерирует новый случайный UID (256 бит) с использованием криптостойкого
// генератора псевдослучайных чисел.
func NewUID() (UID, error) {
	var u UID
	if _, err := rand.Read(u[:]); err != nil {
		return u, fmt.Errorf("failed to generate UID: %w", err)
	}
	return u, nil
}

// MarshalJSON кодирует UID в hex-строку для JSON.
func (u UID) MarshalJSON() ([]byte, error) {
	dst := make([]byte, hex.EncodedLen(len(u)))
	hex.Encode(dst, u[:])
	return []byte(`"` + string(dst) + `"`), nil
}

// UnmarshalJSON декодирует hex-строку в UID. Принимает как строку в кавычках, так и без них.
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

// ParseUID парсит hex-строку в UID.
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

// IsZero возвращает true, если UID равен нулевому значению (все байты — 0).
func (u UID) IsZero() bool {
	return u == UID{}
}
