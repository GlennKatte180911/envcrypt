package env

import (
	"fmt"
	"io"
	"strings"
)

// ExportFormat defines the output format for exported entries.
type ExportFormat int

const (
	// FormatDotEnv exports entries as KEY=VALUE lines (default .env format).
	FormatDotEnv ExportFormat = iota
	// FormatShell exports entries as shell export statements.
	FormatShell
	// FormatJSON exports entries as a JSON object.
	FormatJSON
	// FormatDockerEnv exports entries as Docker --env-file compatible lines.
	FormatDockerEnv
)

// ExportOptions configures the behaviour of Export.
type ExportOptions struct {
	Format  ExportFormat
	Comment string // optional header comment
}

// Export writes entries to w in the requested format.
func Export(w io.Writer, entries []Entry, opts ExportOptions) error {
	if opts.Comment != "" {
		if opts.Format == FormatJSON {
			// JSON does not support comments; skip silently.
		} else {
			for _, line := range strings.Split(opts.Comment, "\n") {
				if _, err := fmt.Fprintf(w, "# %s\n", line); err != nil {
					return err
				}
			}
		}
	}

	switch opts.Format {
	case FormatDotEnv, FormatDockerEnv:
		return exportDotEnv(w, entries)
	case FormatShell:
		return exportShell(w, entries)
	case FormatJSON:
		return exportJSON(w, entries)
	default:
		return fmt.Errorf("env: unknown export format %d", opts.Format)
	}
}

func exportDotEnv(w io.Writer, entries []Entry) error {
	for _, e := range entries {
		if _, err := fmt.Fprintf(w, "%s=%s\n", e.Key, e.Value); err != nil {
			return err
		}
	}
	return nil
}

func exportShell(w io.Writer, entries []Entry) error {
	for _, e := range entries {
		quoted := strings.ReplaceAll(e.Value, "'", "'\\''")
		if _, err := fmt.Fprintf(w, "export %s='%s'\n", e.Key, quoted); err != nil {
			return err
		}
	}
	return nil
}

func exportJSON(w io.Writer, entries []Entry) error {
	if _, err := fmt.Fprint(w, "{\n"); err != nil {
		return err
	}
	for i, e := range entries {
		comma := ","
		if i == len(entries)-1 {
			comma = ""
		}
		escaped := strings.ReplaceAll(e.Value, `"`, `\"`)
		if _, err := fmt.Fprintf(w, "  %q: %q%s\n", e.Key, escaped, comma); err != nil {
			return err
		}
	}
	_, err := fmt.Fprint(w, "}\n")
	return err
}
