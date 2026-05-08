// Package env — patch module
//
// Patch applies a declarative sequence of mutations to a slice of [Entry]
// values without modifying the original slice.
//
// Supported operations:
//
//	"set"    – update an existing key's value, or append a new entry
//	"delete" – remove an entry by key (returns an error if the key is absent)
//	"rename" – change a key's name in-place (returns an error if absent)
//
// Example:
//
//	entries := []env.Entry{
//		{Key: "HOST", Value: "localhost"},
//		{Key: "PORT", Value: "8080"},
//	}
//
//	ops := []env.PatchOp{
//		{Op: "set",    Key: "PORT",  Value: "9090"},
//		{Op: "rename", Key: "HOST",  NewKey: "HOSTNAME"},
//		{Op: "set",    Key: "DEBUG", Value: "true"},
//	}
//
//	result, err := env.Patch(entries, ops)
//
// Patch integrates naturally with [Pipeline] — wrap it in a [PipeFunc] to
// compose it with other transformation stages.
package env
