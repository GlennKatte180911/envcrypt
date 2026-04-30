package env

import "fmt"

// SchemaField describes an expected environment variable.
type SchemaField struct {
	Key      string
	Required bool
	Default  string // used only when Required is false
}

// Schema is an ordered list of SchemaFields that describe the expected shape
// of an env file.
type Schema []SchemaField

// MissingKeyError is returned when a required key is absent from the entries.
type MissingKeyError struct {
	Key string
}

func (e *MissingKeyError) Error() string {
	return fmt.Sprintf("required key %q is missing", e.Key)
}

// ApplySchema validates entries against the schema and returns a normalised
// slice that includes default values for optional absent keys.
// It returns all missing-required-key errors encountered.
func ApplySchema(entries []Entry, schema Schema) ([]Entry, []error) {
	m := ToMap(entries)
	var errs []error
	var result []Entry

	for _, field := range schema {
		val, ok := m[field.Key]
		switch {
		case ok:
			result = append(result, Entry{Key: field.Key, Value: val})
		case field.Required:
			errs = append(errs, &MissingKeyError{Key: field.Key})
		default:
			result = append(result, Entry{Key: field.Key, Value: field.Default})
		}
	}
	return result, errs
}
