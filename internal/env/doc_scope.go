// Package env provides the Scope type for working with prefixed environment
// variable namespaces.
//
// # Overview
//
// Many teams store configuration for multiple deployment targets in a single
// .env file by using a naming convention such as:
//
//	PROD_DB_HOST=db.prod.example.com
//	STAGING_DB_HOST=db.staging.example.com
//
// A Scope extracts and presents all keys that share a common prefix as a
// clean, prefix-free view.
//
// # Usage
//
//	entries := []env.Entry{
//		{Key: "PROD_DB_HOST", Value: "db.prod.example.com"},
//		{Key: "PROD_DB_PORT", Value: "5432"},
//	}
//
//	s := env.NewScope("prod", entries)
//	v, ok := s.Lookup("DB_HOST") // "db.prod.example.com", true
//
// Promote() reverses the operation, re-attaching the prefix so the entries
// can be merged back into a full env set.
//
// ScopeNames() inspects a flat entry list and returns prefix names that
// appear on at least two keys, acting as an autodiscovery helper.
package env
