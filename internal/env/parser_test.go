package env_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envcrypt/internal/env"
)

func TestParseBasic(t *testing.T) {
	input := `
# This is a comment
DB_HOST=localhost
DB_PORT=5432
SECRET_KEY="my secret"
API_TOKEN='token123'
`
	entries, err := env.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(entries))
	}

	cases := []env.Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "SECRET_KEY", Value: "my secret"},
		{Key: "API_TOKEN", Value: "token123"},
	}
	for i, want := range cases {
		if entries[i] != want {
			t.Errorf("entry %d: got %+v, want %+v", i, entries[i], want)
		}
	}
}

func TestParseEmptyInput(t *testing.T) {
	entries, err := env.Parse(strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestParseInvalidLine(t *testing.T) {
	_, err := env.Parse(strings.NewReader("NOEQUALSIGN\n"))
	if err == nil {
		t.Fatal("expected error for invalid line, got nil")
	}
}

func TestParseEmptyKey(t *testing.T) {
	_, err := env.Parse(strings.NewReader("=value\n"))
	if err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestSerializeRoundTrip(t *testing.T) {
	original := []env.Entry{
		{Key: "FOO", Value: "bar"},
		{Key: "BAZ", Value: "qux"},
	}

	var buf strings.Builder
	if err := env.Serialize(&buf, original); err != nil {
		t.Fatalf("serialize error: %v", err)
	}

	parsed, err := env.Parse(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatalf("parse error after serialize: %v", err)
	}
	if len(parsed) != len(original) {
		t.Fatalf("expected %d entries, got %d", len(original), len(parsed))
	}
	for i := range original {
		if parsed[i] != original[i] {
			t.Errorf("entry %d mismatch: got %+v, want %+v", i, parsed[i], original[i])
		}
	}
}
