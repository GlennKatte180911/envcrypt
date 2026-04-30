package env

import (
	"testing"
)

func TestApplySchemaAllPresent(t *testing.T) {
	entries := []Entry{
		{Key: "APP_NAME", Value: "envcrypt"},
		{Key: "APP_PORT", Value: "8080"},
	}
	schema := Schema{
		{Key: "APP_NAME", Required: true},
		{Key: "APP_PORT", Required: true},
	}
	result, errs := ApplySchema(entries, schema)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestApplySchemaRequiredMissing(t *testing.T) {
	entries := []Entry{
		{Key: "APP_NAME", Value: "envcrypt"},
	}
	schema := Schema{
		{Key: "APP_NAME", Required: true},
		{Key: "APP_SECRET", Required: true},
	}
	_, errs := ApplySchema(entries, schema)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	me, ok := errs[0].(*MissingKeyError)
	if !ok {
		t.Fatal("expected *MissingKeyError")
	}
	if me.Key != "APP_SECRET" {
		t.Errorf("expected APP_SECRET, got %s", me.Key)
	}
}

func TestApplySchemaDefaultApplied(t *testing.T) {
	entries := []Entry{
		{Key: "APP_NAME", Value: "envcrypt"},
	}
	schema := Schema{
		{Key: "APP_NAME", Required: true},
		{Key: "APP_PORT", Required: false, Default: "3000"},
	}
	result, errs := ApplySchema(entries, schema)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	m := ToMap(result)
	if m["APP_PORT"] != "3000" {
		t.Errorf("expected default 3000, got %s", m["APP_PORT"])
	}
}

func TestApplySchemaOrderPreserved(t *testing.T) {
	entries := []Entry{
		{Key: "Z_KEY", Value: "z"},
		{Key: "A_KEY", Value: "a"},
	}
	schema := Schema{
		{Key: "A_KEY", Required: true},
		{Key: "Z_KEY", Required: true},
	}
	result, _ := ApplySchema(entries, schema)
	if result[0].Key != "A_KEY" || result[1].Key != "Z_KEY" {
		t.Errorf("schema order not preserved: %v", result)
	}
}

func TestMissingKeyErrorMessage(t *testing.T) {
	e := &MissingKeyError{Key: "SECRET"}
	want := `required key "SECRET" is missing`
	if e.Error() != want {
		t.Errorf("got %q, want %q", e.Error(), want)
	}
}
