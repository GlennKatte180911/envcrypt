package env

import (
	"testing"
)

func redactEntries() []Entry {
	return []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: "supersecret"},
		{Key: "API_TOKEN", Value: "tok_abc123"},
		{Key: "AWS_SECRET_KEY", Value: "aws/secret"},
		{Key: "DEBUG", Value: "true"},
	}
}

func TestRedactNoOptions(t *testing.T) {
	entries := redactEntries()
	result := Redact(entries)
	for i, e := range result {
		if e.Value != entries[i].Value {
			t.Errorf("expected value %q for key %q, got %q", entries[i].Value, e.Key, e.Value)
		}
	}
}

func TestRedactExactKeys(t *testing.T) {
	result := Redact(redactEntries(), WithRedactKeys("DB_PASSWORD", "DEBUG"))
	want := map[string]string{
		"APP_NAME":      "myapp",
		"DB_PASSWORD":   "***",
		"API_TOKEN":     "tok_abc123",
		"AWS_SECRET_KEY": "aws/secret",
		"DEBUG":         "***",
	}
	for _, e := range result {
		if e.Value != want[e.Key] {
			t.Errorf("key %q: expected %q, got %q", e.Key, want[e.Key], e.Value)
		}
	}
}

func TestRedactSuffixes(t *testing.T) {
	result := Redact(redactEntries(), WithRedactSuffixes("_PASSWORD", "_TOKEN", "_KEY"))
	redacted := map[string]bool{
		"DB_PASSWORD":    true,
		"API_TOKEN":      true,
		"AWS_SECRET_KEY": true,
	}
	for _, e := range result {
		if redacted[e.Key] && e.Value != "***" {
			t.Errorf("key %q should be redacted, got %q", e.Key, e.Value)
		}
		if !redacted[e.Key] && e.Value == "***" {
			t.Errorf("key %q should not be redacted", e.Key)
		}
	}
}

func TestRedactCustomPlaceholder(t *testing.T) {
	result := Redact(redactEntries(),
		WithRedactKeys("API_TOKEN"),
		WithRedactPlaceholder("<REDACTED>"),
	)
	for _, e := range result {
		if e.Key == "API_TOKEN" && e.Value != "<REDACTED>" {
			t.Errorf("expected <REDACTED>, got %q", e.Value)
		}
	}
}

func TestRedactDoesNotMutateOriginal(t *testing.T) {
	original := redactEntries()
	Redact(original, WithRedactKeys("DB_PASSWORD"))
	for _, e := range original {
		if e.Key == "DB_PASSWORD" && e.Value != "supersecret" {
			t.Error("original entries were mutated")
		}
	}
}

func TestRedactSuffixCaseInsensitive(t *testing.T) {
	entries := []Entry{
		{Key: "my_secret", Value: "val"},
		{Key: "MY_SECRET", Value: "val"},
	}
	result := Redact(entries, WithRedactSuffixes("_SECRET"))
	for _, e := range result {
		if e.Value != "***" {
			t.Errorf("key %q should be redacted", e.Key)
		}
	}
}
