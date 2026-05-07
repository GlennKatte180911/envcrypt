// Package env provides utilities for parsing, transforming, and managing
// environment variable entries.
//
// # Template Rendering
//
// The template sub-feature allows Go text/template expressions to be embedded
// inside .env values or arbitrary strings, resolved against a set of entries.
//
// Use [RenderTemplate] to render a standalone template string:
//
//	result, err := env.RenderTemplate("{{.HOST}}:{{.PORT}}", entries)
//
// Use [RenderEntries] to render each entry value as a template, enabling
// cross-references between entries:
//
//	entries := []env.Entry{
//		{Key: "HOST", Value: "localhost"},
//		{Key: "PORT", Value: "8080"},
//		{Key: "ADDR", Value: "{{.HOST}}:{{.PORT}}"},
//	}
//	rendered, err := env.RenderEntries(entries)
//	// rendered[2].Value == "localhost:8080"
//
// Missing keys produce a [TemplateError] identifying the offending entry key.
package env
