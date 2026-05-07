package env

import (
	"strings"
	"testing"
)

func templateEntries() []Entry {
	return []Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "PORT", Value: "8080"},
		{Key: "APP_NAME", Value: "envcrypt"},
	}
}

func TestRenderTemplateBasic(t *testing.T) {
	entries := templateEntries()
	out, err := RenderTemplate("{{.HOST}}:{{.PORT}}", entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "localhost:8080" {
		t.Errorf("expected 'localhost:8080', got %q", out)
	}
}

func TestRenderTemplateNoRefs(t *testing.T) {
	out, err := RenderTemplate("static string", templateEntries())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "static string" {
		t.Errorf("expected 'static string', got %q", out)
	}
}

func TestRenderTemplateMissingKey(t *testing.T) {
	_, err := RenderTemplate("{{.MISSING}}", templateEntries())
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
	var te *TemplateError
	if !isTemplateError(err, &te) {
		t.Errorf("expected TemplateError, got %T", err)
	}
}

func TestRenderTemplateInvalidSyntax(t *testing.T) {
	_, err := RenderTemplate("{{.HOST", templateEntries())
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
}

func TestRenderEntriesBasic(t *testing.T) {
	entries := []Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "PORT", Value: "8080"},
		{Key: "ADDR", Value: "{{.HOST}}:{{.PORT}}"},
	}

	result, err := RenderEntries(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	m := ToMap(result)
	if m["ADDR"] != "localhost:8080" {
		t.Errorf("expected 'localhost:8080', got %q", m["ADDR"])
	}
	if m["HOST"] != "localhost" {
		t.Errorf("HOST should be unchanged, got %q", m["HOST"])
	}
}

func TestRenderEntriesPreservesComment(t *testing.T) {
	entries := []Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "ADDR", Value: "{{.HOST}}", Comment: "the address"},
	}

	result, err := RenderEntries(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result[1].Comment != "the address" {
		t.Errorf("expected comment preserved, got %q", result[1].Comment)
	}
}

func TestRenderEntriesMissingKey(t *testing.T) {
	entries := []Entry{
		{Key: "ADDR", Value: "{{.MISSING}}"},
	}

	_, err := RenderEntries(entries)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "ADDR") {
		t.Errorf("expected error to mention key 'ADDR', got: %v", err)
	}
}

func TestTemplateErrorMessage(t *testing.T) {
	e := &TemplateError{Key: "FOO", Cause: fmt.Errorf("bad")}
	if !strings.Contains(e.Error(), "FOO") {
		t.Errorf("expected key in error message, got: %s", e.Error())
	}
}

// isTemplateError is a helper to avoid importing errors package in test.
func isTemplateError(err error, out **TemplateError) bool {
	if te, ok := err.(*TemplateError); ok {
		*out = te
		return true
	}
	return false
}
