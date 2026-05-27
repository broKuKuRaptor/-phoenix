package common

import (
	"encoding/json"
	"testing"
)

func TestUID_MarshalUnmarshal(t *testing.T) {
	uid := UID{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	}

	// Marshal
	data, err := json.Marshal(uid)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	expected := `"000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"`
	if string(data) != expected {
		t.Errorf("Marshal = %s, want %s", data, expected)
	}

	// Unmarshal
	var uid2 UID
	if err := json.Unmarshal(data, &uid2); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if uid != uid2 {
		t.Errorf("roundtrip failed: %x != %x", uid, uid2)
	}
}

func TestUID_UnmarshalInvalid(t *testing.T) {
	tests := []string{
		`""`,
		`"too_short"`,
		`"GGGG"`, // не hex
	}
	for _, input := range tests {
		var u UID
		err := json.Unmarshal([]byte(input), &u)
		if err == nil {
			t.Errorf("Unmarshal(%s) = nil, want error", input)
		}
	}
}

func TestParseUID(t *testing.T) {
	hex := "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"
	uid, err := ParseUID(hex)
	if err != nil {
		t.Fatalf("ParseUID: %v", err)
	}
	if uid.String() != hex {
		t.Errorf("String() = %s, want %s", uid.String(), hex)
	}
}

func TestParseUID_Invalid(t *testing.T) {
	_, err := ParseUID("bad")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUID_String(t *testing.T) {
	uid := UID{}
	if len(uid.String()) != 64 {
		t.Errorf("len(String()) = %d, want 64", len(uid.String()))
	}
}
