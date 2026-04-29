// Package env provides utilities for reading, parsing, and serializing
// .env files used by the envcrypt tool.
//
// A .env file is a plain-text file containing KEY=VALUE pairs, one per line.
// Lines beginning with '#' are treated as comments and ignored during parsing.
// Values may optionally be wrapped in single or double quotes, which are
// stripped during parsing but not preserved on serialization.
//
// Typical usage:
//
//	// Parse an existing .env file.
//	f, err := os.Open(".env")
//	if err != nil { ... }
//	defer f.Close()
//
//	entries, err := env.Parse(f)
//	if err != nil { ... }
//
//	// Modify entries, then serialize back.
//	var buf bytes.Buffer
//	if err := env.Serialize(&buf, entries); err != nil { ... }
//
// The parsed entries are order-preserving, which ensures that serialized
// output maintains the same key ordering as the original file.
package env
