package env

import (
	"fmt"
	"strings"
	"text/template"
)

// TemplateError is returned when template rendering fails.
type TemplateError struct {
	Key string
	Cause error
}

func (e *TemplateError) Error() string {
	return fmt.Sprintf("template error for key %q: %v", e.Key, e.Cause)
}

// RenderTemplate executes a Go text/template string using the provided entries
// as the data context. The entries are converted to a map[string]string before
// rendering, so template expressions like {{.MY_VAR}} resolve to entry values.
//
// Example template: "{{.APP_HOST}}:{{.APP_PORT}}"
func RenderTemplate(tmpl string, entries []Entry) (string, error) {
	data := ToMap(entries)

	t, err := template.New("").Option("missingkey=error").Parse(tmpl)
	if err != nil {
		return "", &TemplateError{Key: "<template>", Cause: err}
	}

	var buf strings.Builder
	if err := t.Execute(&buf, data); err != nil {
		return "", &TemplateError{Key: "<template>", Cause: err}
	}

	return buf.String(), nil
}

// RenderEntries applies RenderTemplate to each entry value, treating the
// values as templates. Entries are resolved independently using the full
// entry set as context, allowing cross-references like VALUE={{.OTHER}}.
// Returns a new slice with rendered values; original entries are unchanged.
func RenderEntries(entries []Entry) ([]Entry, error) {
	data := ToMap(entries)

	result := make([]Entry, 0, len(entries))
	for _, e := range entries {
		t, err := template.New(e.Key).Option("missingkey=error").Parse(e.Value)
		if err != nil {
			return nil, &TemplateError{Key: e.Key, Cause: err}
		}

		var buf strings.Builder
		if err := t.Execute(&buf, data); err != nil {
			return nil, &TemplateError{Key: e.Key, Cause: err}
		}

		result = append(result, Entry{Key: e.Key, Value: buf.String(), Comment: e.Comment})
	}

	return result, nil
}
