package env

import (
	"testing"
)

func TestCensusEmpty(t *testing.T) {
	result := TakeCensus(nil)
	if result.Total != 0 {
		t.Errorf("expected Total=0, got %d", result.Total)
	}
	if result.AverageValueLen != 0 {
		t.Errorf("expected AverageValueLen=0, got %f", result.AverageValueLen)
	}
}

func TestCensusTotal(t *testing.T) {
	entries := []Entry{
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
		{Key: "C", Value: "3"},
	}
	r := TakeCensus(entries)
	if r.Total != 3 {
		t.Errorf("expected Total=3, got %d", r.Total)
	}
}

func TestCensusEmptyValues(t *testing.T) {
	entries := []Entry{
		{Key: "A", Value: ""},
		{Key: "B", Value: "hello"},
		{Key: "C", Value: ""},
	}
	r := TakeCensus(entries)
	if r.EmptyValues != 2 {
		t.Errorf("expected EmptyValues=2, got %d", r.EmptyValues)
	}
}

func TestCensusUniqueAndDuplicateKeys(t *testing.T) {
	entries := []Entry{
		{Key: "A", Value: "1"},
		{Key: "A", Value: "2"},
		{Key: "B", Value: "3"},
	}
	r := TakeCensus(entries)
	if r.UniqueKeys != 2 {
		t.Errorf("expected UniqueKeys=2, got %d", r.UniqueKeys)
	}
	if r.DuplicateKeys != 1 {
		t.Errorf("expected DuplicateKeys=1, got %d", r.DuplicateKeys)
	}
}

func TestCensusUppercaseAndLowercaseKeys(t *testing.T) {
	entries := []Entry{
		{Key: "FOO_BAR", Value: "v1"},
		{Key: "foo_bar", Value: "v2"},
		{Key: "MixedCase", Value: "v3"},
		{Key: "ANOTHER", Value: "v4"},
	}
	r := TakeCensus(entries)
	if r.UppercaseKeys != 2 {
		t.Errorf("expected UppercaseKeys=2, got %d", r.UppercaseKeys)
	}
	if r.LowercaseKeys != 1 {
		t.Errorf("expected LowercaseKeys=1, got %d", r.LowercaseKeys)
	}
}

func TestCensusAverageValueLen(t *testing.T) {
	entries := []Entry{
		{Key: "A", Value: "ab"},   // len 2
		{Key: "B", Value: "abcd"}, // len 4
	}
	r := TakeCensus(entries)
	expected := 3.0
	if r.AverageValueLen != expected {
		t.Errorf("expected AverageValueLen=%f, got %f", expected, r.AverageValueLen)
	}
}

func TestCensusNoDuplicatesWhenAllUnique(t *testing.T) {
	entries := []Entry{
		{Key: "X", Value: "1"},
		{Key: "Y", Value: "2"},
	}
	r := TakeCensus(entries)
	if r.DuplicateKeys != 0 {
		t.Errorf("expected DuplicateKeys=0, got %d", r.DuplicateKeys)
	}
}
