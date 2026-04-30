// Package env provides utilities for parsing, serialising, filtering,
// diffing, and merging .env files.
//
// # Validation
//
// The validate.go file exposes a rule-based validation system.
// Rules are plain functions of type ValidationRule that accept an Entry and
// return an error when the entry violates a constraint.
//
// Built-in rules:
//
//	- RequireNonEmpty  – rejects entries with blank values.
//	- RequirePrefix    – rejects entries whose key lacks a given prefix.
//	- ForbidPrefix     – rejects entries whose key carries a forbidden prefix.
//
// # Schema
//
// The schema.go file provides a lightweight schema mechanism.  A Schema is a
// slice of SchemaField values that declare which keys are required and what
// default values optional keys should receive when absent.
//
// ApplySchema normalises an entry slice according to the schema, filling in
// defaults and collecting MissingKeyError values for absent required keys.
package env
