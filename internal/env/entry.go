package env

// Entry represents a single key-value pair from a .env file.
// Comments and blank lines are preserved via the Comment field.
type Entry struct {
	// Key is the environment variable name. Empty for comment/blank lines.
	Key string

	// Value is the raw (unquoted) value of the variable.
	Value string

	// Comment is an optional inline comment (text after ' #' on the same line).
	Comment string

	// IsComment is true when the entire line is a comment or blank line.
	IsComment bool
}

// IsBlank reports whether this entry carries no meaningful data.
func (e Entry) IsBlank() bool {
	return e.Key == "" && !e.IsComment
}

// String returns a human-readable representation of the entry,
// primarily useful for debugging.
func (e Entry) String() string {
	if e.IsComment {
		return "# " + e.Comment
	}
	if e.Comment != "" {
		return e.Key + "=" + e.Value + " # " + e.Comment
	}
	return e.Key + "=" + e.Value
}
