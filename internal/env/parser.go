// Package env provides utilities for parsing and serializing .env files.
package env

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Entry represents a single key-value pair from a .env file.
type Entry struct {
	Key   string
	Value string
}

// Parse reads a .env file from the given reader and returns a slice of entries.
// Lines beginning with '#' are treated as comments and skipped.
// Blank lines are also skipped.
func Parse(r io.Reader) ([]Entry, error) {
	var entries []Entry
	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("line %d: invalid format %q", lineNum, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "" {
			return nil, fmt.Errorf("line %d: empty key", lineNum)
		}

		// Strip optional surrounding quotes from value.
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		entries = append(entries, Entry{Key: key, Value: value})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning input: %w", err)
	}

	return entries, nil
}

// Serialize writes a slice of entries to the given writer in .env format.
func Serialize(w io.Writer, entries []Entry) error {
	for _, e := range entries {
		if _, err := fmt.Fprintf(w, "%s=%s\n", e.Key, e.Value); err != nil {
			return fmt.Errorf("writing entry %q: %w", e.Key, err)
		}
	}
	return nil
}
