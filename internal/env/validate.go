package env

import (
	"fmt"
	"strings"
)

// ValidationError describes a single validation failure for an entry.
type ValidationError struct {
	Key     string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for key %q: %s", e.Key, e.Message)
}

// ValidationRule is a function that validates a single Entry.
// It returns a non-nil error if the entry is invalid.
type ValidationRule func(entry Entry) error

// RequireNonEmpty returns a rule that rejects entries whose value is empty.
func RequireNonEmpty() ValidationRule {
	return func(e Entry) error {
		if strings.TrimSpace(e.Value) == "" {
			return &ValidationError{Key: e.Key, Message: "value must not be empty"}
		}
		return nil
	}
}

// RequirePrefix returns a rule that rejects entries whose key does not start
// with the given prefix.
func RequirePrefix(prefix string) ValidationRule {
	return func(e Entry) error {
		if !strings.HasPrefix(e.Key, prefix) {
			return &ValidationError{Key: e.Key, Message: fmt.Sprintf("key must start with %q", prefix)}
		}
		return nil
	}
}

// ForbidPrefix returns a rule that rejects entries whose key starts with the
// given prefix.
func ForbidPrefix(prefix string) ValidationRule {
	return func(e Entry) error {
		if strings.HasPrefix(e.Key, prefix) {
			return &ValidationError{Key: e.Key, Message: fmt.Sprintf("key must not start with %q", prefix)}
		}
		return nil
	}
}

// Validate applies all rules to every entry and returns all validation errors
// encountered. If no errors are found it returns nil.
func Validate(entries []Entry, rules ...ValidationRule) []error {
	var errs []error
	for _, entry := range entries {
		for _, rule := range rules {
			if err := rule(entry); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errs
}
