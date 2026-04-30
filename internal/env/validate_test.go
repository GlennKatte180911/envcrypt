package env

import (
	"testing"
)

func validationEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "envcrypt"},
		{Key: "APP_SECRET", Value: ""},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PASSWORD", Value: ""},
	}
}

func TestValidateNoRules(t *testing.T) {
	errs := Validate(validationEntries())
	if len(errs) != 0 {
		t.Fatalf("expected no errors with no rules, got %d", len(errs))
	}
}

func TestRequireNonEmpty(t *testing.T) {
	errs := Validate(validationEntries(), RequireNonEmpty())
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(errs))
	}
}

func TestRequirePrefix(t *testing.T) {
	entries := []Entry{
		{Key: "APP_NAME", Value: "x"},
		{Key: "DB_HOST", Value: "y"},
	}
	errs := Validate(entries, RequirePrefix("APP_"))
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	ve, ok := errs[0].(*ValidationError)
	if !ok {
		t.Fatal("expected *ValidationError")
	}
	if ve.Key != "DB_HOST" {
		t.Errorf("expected error for DB_HOST, got %s", ve.Key)
	}
}

func TestForbidPrefix(t *testing.T) {
	entries := []Entry{
		{Key: "SECRET_TOKEN", Value: "abc"},
		{Key: "APP_NAME", Value: "envcrypt"},
	}
	errs := Validate(entries, ForbidPrefix("SECRET_"))
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestValidateMultipleRules(t *testing.T) {
	entries := []Entry{
		{Key: "APP_NAME", Value: ""},
		{Key: "DB_HOST", Value: "localhost"},
	}
	errs := Validate(entries, RequireNonEmpty(), RequirePrefix("APP_"))
	// APP_NAME: empty value → 1 error
	// DB_HOST: missing prefix → 1 error
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(errs))
	}
}

func TestValidationErrorMessage(t *testing.T) {
	ve := &ValidationError{Key: "FOO", Message: "must not be empty"}
	want := `validation error for key "FOO": must not be empty`
	if ve.Error() != want {
		t.Errorf("got %q, want %q", ve.Error(), want)
	}
}
