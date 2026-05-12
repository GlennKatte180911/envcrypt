package env

import (
	"testing"
)

func TestFlattenNestedSimple(t *testing.T) {
	nested := map[string]interface{}{
		"DB": map[string]interface{}{
			"HOST": "localhost",
			"PORT": "5432",
		},
		"APP_NAME": "envcrypt",
	}

	entries := FlattenNested(nested)
	m := ToMap(entries)

	if m["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", m["DB_HOST"])
	}
	if m["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", m["DB_PORT"])
	}
	if m["APP_NAME"] != "envcrypt" {
		t.Errorf("expected APP_NAME=envcrypt, got %q", m["APP_NAME"])
	}
}

func TestFlattenNestedCustomSeparator(t *testing.T) {
	nested := map[string]interface{}{
		"REDIS": map[string]interface{}{
			"HOST": "127.0.0.1",
		},
	}

	entries := FlattenNested(nested, WithSeparator("."))
	m := ToMap(entries)

	if m["REDIS.HOST"] != "127.0.0.1" {
		t.Errorf("expected REDIS.HOST=127.0.0.1, got %q", m["REDIS.HOST"])
	}
}

func TestFlattenNestedWithPrefix(t *testing.T) {
	nested := map[string]interface{}{
		"NAME": "test",
	}

	entries := FlattenNested(nested, WithFlattenPrefix("APP"))
	m := ToMap(entries)

	if m["APP_NAME"] != "test" {
		t.Errorf("expected APP_NAME=test, got %q", m["APP_NAME"])
	}
}

func TestFlattenNestedNonStringValue(t *testing.T) {
	nested := map[string]interface{}{
		"MAX_CONN": 10,
	}

	entries := FlattenNested(nested)
	m := ToMap(entries)

	if m["MAX_CONN"] != "10" {
		t.Errorf("expected MAX_CONN=10, got %q", m["MAX_CONN"])
	}
}

func TestFlattenNestedEmpty(t *testing.T) {
	entries := FlattenNested(map[string]interface{}{})
	if len(entries) != 0 {
		t.Errorf("expected empty entries, got %d", len(entries))
	}
}

func TestUnflattenToMap(t *testing.T) {
	entries := []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "APP_NAME", Value: "envcrypt"},
	}

	result := UnflattenToMap(entries, "_")

	db, ok := result["DB"].(map[string]interface{})
	if !ok {
		t.Fatal("expected DB to be a nested map")
	}
	if db["HOST"] != "localhost" {
		t.Errorf("expected DB.HOST=localhost, got %v", db["HOST"])
	}
	if db["PORT"] != "5432" {
		t.Errorf("expected DB.PORT=5432, got %v", db["PORT"])
	}
	if result["APP"] == nil {
		t.Error("expected APP key to be present")
	}
}

func TestUnflattenDefaultSeparator(t *testing.T) {
	entries := []Entry{{Key: "A_B", Value: "v"}}
	result := UnflattenToMap(entries, "")
	a, ok := result["A"].(map[string]interface{})
	if !ok {
		t.Fatal("expected A to be nested map with default separator")
	}
	if a["B"] != "v" {
		t.Errorf("expected A.B=v, got %v", a["B"])
	}
}
