package env

import (
	"testing"
)

var auditEntries = []Entry{
	{Key: "HOST", Value: "localhost"},
	{Key: "PORT", Value: ""},
	{Key: "db_password", Value: "secret"},
	{Key: "HOST", Value: "remotehost"},
}

func TestAuditNoFindings(t *testing.T) {
	entries := []Entry{{Key: "HOST", Value: "localhost"}}
	findings := Audit(entries, AuditNoEmptyValues(), AuditNoDuplicateKeys(), AuditRequireUppercaseKeys())
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

func TestAuditNoEmptyValues(t *testing.T) {
	findings := Audit(auditEntries, AuditNoEmptyValues())
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Key != "PORT" {
		t.Errorf("expected finding for PORT, got %s", findings[0].Key)
	}
	if findings[0].Severity != SeverityWarning {
		t.Errorf("expected WARNING severity, got %s", findings[0].Severity)
	}
}

func TestAuditNoDuplicateKeys(t *testing.T) {
	findings := Audit(auditEntries, AuditNoDuplicateKeys())
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Key != "HOST" {
		t.Errorf("expected finding for HOST, got %s", findings[0].Key)
	}
	if findings[0].Severity != SeverityError {
		t.Errorf("expected ERROR severity, got %s", findings[0].Severity)
	}
}

func TestAuditRequireUppercaseKeys(t *testing.T) {
	findings := Audit(auditEntries, AuditRequireUppercaseKeys())
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Key != "db_password" {
		t.Errorf("expected finding for db_password, got %s", findings[0].Key)
	}
	if findings[0].Severity != SeverityInfo {
		t.Errorf("expected INFO severity, got %s", findings[0].Severity)
	}
}

func TestAuditMultipleRules(t *testing.T) {
	findings := Audit(auditEntries, AuditNoEmptyValues(), AuditNoDuplicateKeys(), AuditRequireUppercaseKeys())
	if len(findings) != 3 {
		t.Fatalf("expected 3 findings, got %d", len(findings))
	}
}

func TestAuditFindingErrorString(t *testing.T) {
	f := AuditFinding{Key: "PORT", Message: "value is empty", Severity: SeverityWarning}
	want := "[WARNING] PORT: value is empty"
	if f.Error() != want {
		t.Errorf("expected %q, got %q", want, f.Error())
	}
}

func TestAuditNoRules(t *testing.T) {
	findings := Audit(auditEntries)
	if findings != nil {
		t.Errorf("expected nil findings with no rules, got %v", findings)
	}
}
