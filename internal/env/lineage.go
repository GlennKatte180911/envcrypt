package env

import "time"

// ChangeKind describes the type of change recorded in a Lineage entry.
type ChangeKind string

const (
	ChangeKindSet    ChangeKind = "set"
	ChangeKindDelete ChangeKind = "delete"
	ChangeKindMerge  ChangeKind = "merge"
)

// LineageRecord captures a single mutation event applied to a set of entries.
type LineageRecord struct {
	Timestamp time.Time
	Kind      ChangeKind
	Label     string
	Diff      []DiffEntry
}

// Lineage tracks a sequence of changes over time, providing a simple
// audit trail of how a set of env entries has evolved.
type Lineage struct {
	records []LineageRecord
}

// NewLineage returns an empty Lineage ready to record changes.
func NewLineage() *Lineage {
	return &Lineage{}
}

// Record appends a new LineageRecord built from the diff between base and
// updated, labelled with the supplied kind and label.
func (l *Lineage) Record(kind ChangeKind, label string, base, updated []Entry) {
	d := Diff(base, updated)
	if len(d) == 0 {
		return
	}
	l.records = append(l.records, LineageRecord{
		Timestamp: time.Now().UTC(),
		Kind:      kind,
		Label:     label,
		Diff:      d,
	})
}

// Records returns a read-only copy of all recorded LineageRecords.
func (l *Lineage) Records() []LineageRecord {
	out := make([]LineageRecord, len(l.records))
	copy(out, l.records)
	return out
}

// Len returns the number of records stored in the lineage.
func (l *Lineage) Len() int {
	return len(l.records)
}

// ChangedKeys returns the deduplicated set of keys that appear in any diff
// across all recorded changes.
func (l *Lineage) ChangedKeys() []string {
	seen := map[string]struct{}{}
	for _, r := range l.records {
		for _, d := range r.Diff {
			seen[d.Key] = struct{}{}
		}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	return keys
}
