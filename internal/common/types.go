// Package common provides shared types and helpers used across all services.
package common

import (
	"encoding/hex"
	"fmt"
)

// UID is a fixed-size 32-byte unique identifier.
//
// In JSON it serializes as a 64-character hex string (without 0x prefix).
// Use ParseUID to parse a hex string, or rely on JSON unmarshaling.
//
// The zero value (all bytes zero) is reserved as a sentinel for
// "root" or "system" account and must not be assigned to real entities.
type UID [32]byte

// MarshalJSON encodes the UID as a hex string in JSON.
func (u UID) MarshalJSON() ([]byte, error) {
	dst := make([]byte, hex.EncodedLen(len(u)))
	hex.Encode(dst, u[:])
	return []byte(`"` + string(dst) + `"`), nil
}

// UnmarshalJSON decodes a hex string into the UID byte array. It accepts both quoted and unquoted hex strings.
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

// ParseUID parses a hex string into a UID.
func ParseUID(s string) (UID, error) {
	var u UID
	if err := (&u).UnmarshalJSON([]byte(`"` + s + `"`)); err != nil {
		return u, err
	}
	return u, nil
}

// String returns the UID as a hex string (without quotes).
func (u UID) String() string {
	return hex.EncodeToString(u[:])
}
