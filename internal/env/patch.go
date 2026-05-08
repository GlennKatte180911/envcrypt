package env

// PatchOp represents a single patch operation to apply to a set of entries.
type PatchOp struct {
	// Op is the operation type: "set", "delete", or "rename".
	Op string
	// Key is the target key for the operation.
	Key string
	// Value is used by the "set" operation.
	Value string
	// NewKey is used by the "rename" operation.
	NewKey string
}

// PatchError records an operation that could not be applied.
type PatchError struct {
	Op  PatchOp
	Msg string
}

func (e *PatchError) Error() string {
	return "patch op " + e.Op.Op + " on key " + e.Op.Key + ": " + e.Msg
}

// Patch applies a sequence of PatchOps to entries and returns the resulting
// slice. The original slice is never modified. Operations are applied in order;
// an unknown Op returns a PatchError immediately.
func Patch(entries []Entry, ops []PatchOp) ([]Entry, error) {
	out := Clone(entries)
	for _, op := range ops {
		var err error
		switch op.Op {
		case "set":
			out = patchSet(out, op.Key, op.Value)
		case "delete":
			out, err = patchDelete(out, op.Key)
		case "rename":
			out, err = patchRename(out, op.Key, op.NewKey)
		default:
			return nil, &PatchError{Op: op, Msg: "unknown operation"}
		}
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func patchSet(entries []Entry, key, value string) []Entry {
	for i, e := range entries {
		if e.Key == key {
			entries[i].Value = value
			return entries
		}
	}
	return append(entries, Entry{Key: key, Value: value})
}

func patchDelete(entries []Entry, key string) ([]Entry, error) {
	for i, e := range entries {
		if e.Key == key {
			return append(entries[:i], entries[i+1:]...), nil
		}
	}
	return entries, &PatchError{Op: PatchOp{Op: "delete", Key: key}, Msg: "key not found"}
}

func patchRename(entries []Entry, key, newKey string) ([]Entry, error) {
	for i, e := range entries {
		if e.Key == key {
			entries[i].Key = newKey
			return entries, nil
		}
	}
	return entries, &PatchError{Op: PatchOp{Op: "rename", Key: key, NewKey: newKey}, Msg: "key not found"}
}
