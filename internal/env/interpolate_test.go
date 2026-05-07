package env

import (
	"testing"
)

func TestInterpolateNoRefs(t *testing.T) {
	entries := []Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "PORT", Value: "5432"},
	}
	out, err := Interpolate(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "localhost" || out[1].Value != "5432" {
		t.Errorf("values should be unchanged, got %v", out)
	}
}

func TestInterpolateBraceForm(t *testing.T) {
	entries := []Entry{
		{Key: "HOST", Value: "db.example.com"},
		{Key: "DSN", Value: "postgres://${HOST}/mydb"},
	}
	out, err := Interpolate(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "postgres://db.example.com/mydb"
	if out[1].Value != want {
		t.Errorf("want %q, got %q", want, out[1].Value)
	}
}

func TestInterpolateBareForm(t *testing.T) {
	entries := []Entry{
		{Key: "USER", Value: "admin"},
		{Key: "GREETING", Value: "hello $USER!"},
	}
	out, err := Interpolate(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "hello admin!"
	if out[1].Value != want {
		t.Errorf("want %q, got %q", want, out[1].Value)
	}
}

func TestInterpolateUnresolvedRef(t *testing.T) {
	entries := []Entry{
		{Key: "DSN", Value: "postgres://${HOST}/mydb"},
	}
	_, err := Interpolate(entries)
	if err == nil {
		t.Fatal("expected error for unresolved reference, got nil")
	}
	ie, ok := err.(*InterpolateError)
	if !ok {
		t.Fatalf("expected *InterpolateError, got %T", err)
	}
	if ie.Key != "HOST" {
		t.Errorf("expected key %q, got %q", "HOST", ie.Key)
	}
}

func TestInterpolateChained(t *testing.T) {
	entries := []Entry{
		{Key: "SCHEME", Value: "https"},
		{Key: "HOST", Value: "example.com"},
		{Key: "BASE_URL", Value: "${SCHEME}://${HOST}"},
		{Key: "API_URL", Value: "${BASE_URL}/api/v1"},
	}
	out, err := Interpolate(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[2].Value != "https://example.com" {
		t.Errorf("BASE_URL: want %q, got %q", "https://example.com", out[2].Value)
	}
	if out[3].Value != "https://example.com/api/v1" {
		t.Errorf("API_URL: want %q, got %q", "https://example.com/api/v1", out[3].Value)
	}
}

func TestInterpolatePreservesComment(t *testing.T) {
	entries := []Entry{
		{Key: "ENV", Value: "prod"},
		{Key: "TAG", Value: "$ENV-v1", Comment: "deployment tag"},
	}
	out, err := Interpolate(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[1].Comment != "deployment tag" {
		t.Errorf("comment not preserved, got %q", out[1].Comment)
	}
}

func TestInterpolateErrorMessage(t *testing.T) {
	e := &InterpolateError{Key: "SECRET"}
	want := `interpolate: unresolved variable reference: "SECRET"`
	if e.Error() != want {
		t.Errorf("want %q, got %q", want, e.Error())
	}
}
