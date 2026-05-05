package env

import (
	"testing"
)

var transformEntries = []Entry{
	{Key: "db_host", Value: "localhost"},
	{Key: "DB_PORT", Value: "5432"},
	{Key: "api_key", Value: "secret"},
	{Key: "EMPTY_VAL", Value: ""},
}

func TestUppercaseKeys(t *testing.T) {
	out := Transform(transformEntries, UppercaseKeys())
	for _, e := range out {
		if e.Key != toUpper(e.Key) {
			t.Errorf("expected uppercase key, got %q", e.Key)
		}
	}
}

func TestLowercaseKeys(t *testing.T) {
	out := Transform(transformEntries, LowercaseKeys())
	for _, e := range out {
		if e.Key != toLower(e.Key) {
			t.Errorf("expected lowercase key, got %q", e.Key)
		}
	}
}

func TestRenameKey(t *testing.T) {
	out := Transform(transformEntries, RenameKey("db_host", "DATABASE_HOST"))
	found := false
	for _, e := range out {
		if e.Key == "DATABASE_HOST" {
			found = true
		}
		if e.Key == "db_host" {
			t.Error("old key still present after rename")
		}
	}
	if !found {
		t.Error("renamed key not found in output")
	}
}

func TestMaskValues(t *testing.T) {
	out := Transform(transformEntries, MaskValues("***"))
	for _, e := range out {
		if e.Value != "***" {
			t.Errorf("expected masked value, got %q", e.Value)
		}
	}
}

func TestDropEmpty(t *testing.T) {
	out := Transform(transformEntries, DropEmpty())
	for _, e := range out {
		if e.Value == "" {
			t.Errorf("expected no empty values, found key %q with empty value", e.Key)
		}
	}
	if len(out) != len(transformEntries)-1 {
		t.Errorf("expected %d entries, got %d", len(transformEntries)-1, len(out))
	}
}

func TestChain(t *testing.T) {
	chained := Chain(DropEmpty(), UppercaseKeys(), MaskValues("REDACTED"))
	out := Transform(transformEntries, chained)
	for _, e := range out {
		if e.Value == "" {
			t.Error("chain should have dropped empty values")
		}
		if e.Key != toUpper(e.Key) {
			t.Errorf("chain should have uppercased key, got %q", e.Key)
		}
		if e.Value != "REDACTED" {
			t.Errorf("chain should have masked value, got %q", e.Value)
		}
	}
}

// helpers to avoid importing strings in test
func toUpper(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'a' && c <= 'z' {
			b[i] = c - 32
		}
	}
	return string(b)
}

func toLower(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + 32
		}
	}
	return string(b)
}
