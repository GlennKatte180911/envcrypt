// Package env provides utilities for parsing, serializing, and manipulating
// environment variable entries.
//
// # Audit
//
// The audit sub-feature allows teams to run configurable rule-based checks
// against a set of [Entry] values before encrypting or committing them.
//
// Built-in rules:
//
//   - [AuditNoEmptyValues] — warns when an entry has an empty value.
//   - [AuditNoDuplicateKeys] — errors when the same key appears more than once.
//   - [AuditRequireUppercaseKeys] — informs when a key contains lowercase letters.
//
// Use [Audit] to run one or more [AuditRule] functions and collect
// [AuditFinding] results, each tagged with a [AuditSeverity] level:
// INFO, WARNING, or ERROR.
//
// Example:
//
//	findings := env.Audit(entries,
//	    env.AuditNoEmptyValues(),
//	    env.AuditNoDuplicateKeys(),
//	)
package env
