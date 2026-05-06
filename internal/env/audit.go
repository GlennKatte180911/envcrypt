package env

import "fmt"

// AuditSeverity represents the severity level of an audit finding.
type AuditSeverity string

const (
	SeverityInfo    AuditSeverity = "INFO"
	SeverityWarning AuditSeverity = "WARNING"
	SeverityError   AuditSeverity = "ERROR"
)

// AuditFinding describes a single issue found during an audit.
type AuditFinding struct {
	Key      string
	Message  string
	Severity AuditSeverity
}

func (f AuditFinding) Error() string {
	return fmt.Sprintf("[%s] %s: %s", f.Severity, f.Key, f.Message)
}

// AuditRule is a function that inspects entries and returns findings.
type AuditRule func(entries []Entry) []AuditFinding

// Audit runs all provided rules against the entries and returns all findings.
func Audit(entries []Entry, rules ...AuditRule) []AuditFinding {
	var findings []AuditFinding
	for _, rule := range rules {
		findings = append(findings, rule(entries)...)
	}
	return findings
}

// AuditNoEmptyValues returns a rule that flags entries with empty values.
func AuditNoEmptyValues() AuditRule {
	return func(entries []Entry) []AuditFinding {
		var findings []AuditFinding
		for _, e := range entries {
			if e.Value == "" {
				findings = append(findings, AuditFinding{
					Key:      e.Key,
					Message:  "value is empty",
					Severity: SeverityWarning,
				})
			}
		}
		return findings
	}
}

// AuditNoDuplicateKeys returns a rule that flags duplicate keys.
func AuditNoDuplicateKeys() AuditRule {
	return func(entries []Entry) []AuditFinding {
		seen := make(map[string]int)
		for _, e := range entries {
			seen[e.Key]++
		}
		var findings []AuditFinding
		for key, count := range seen {
			if count > 1 {
				findings = append(findings, AuditFinding{
					Key:      key,
					Message:  fmt.Sprintf("key appears %d times", count),
					Severity: SeverityError,
				})
			}
		}
		return findings
	}
}

// AuditRequireUppercaseKeys returns a rule that flags keys that are not fully uppercase.
func AuditRequireUppercaseKeys() AuditRule {
	return func(entries []Entry) []AuditFinding {
		var findings []AuditFinding
		for _, e := range entries {
			for _, ch := range e.Key {
				if ch >= 'a' && ch <= 'z' {
					findings = append(findings, AuditFinding{
						Key:      e.Key,
						Message:  "key contains lowercase characters",
						Severity: SeverityInfo,
					})
					break
				}
			}
		}
		return findings
	}
}
