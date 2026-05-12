// Package env provides utilities for parsing, manipulating, and exporting
// environment variable entries.
//
// # Flatten / Unflatten
//
// FlattenNested converts a nested map[string]interface{} (such as one decoded
// from a JSON or YAML config file) into a flat []Entry slice by joining key
// segments with a separator (default "_").
//
//	nested := map[string]interface{}{
//	    "DB": map[string]interface{}{
//	        "HOST": "localhost",
//	        "PORT": "5432",
//	    },
//	}
//	entries := env.FlattenNested(nested)
//	// → [{Key:"DB_HOST", Value:"localhost"}, {Key:"DB_PORT", Value:"5432"}]
//
// Options:
//
//	env.WithSeparator(".")         — use "." instead of "_"
//	env.WithFlattenPrefix("APP")   — prepend "APP" to every key
//
// UnflattenToMap is the inverse operation: it reconstructs a nested
// map[string]interface{} from a flat []Entry by splitting keys on the
// given separator.
//
//	m := env.UnflattenToMap(entries, "_")
//	// m["DB"].(map[string]interface{})["HOST"] == "localhost"
package env
