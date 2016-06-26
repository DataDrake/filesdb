package filesdb

import (
	"testing"
)

func TestNewHashEntry(t *testing.T) {
	e := NewHashEntry()
	if e.flags != EMPTY {
		t.Log("Entry was not marked empty")
	}
	if e.key != "" {
		t.Logf("Key should be empty, found: %s", e.key)
	}
	if e.value != nil {
		t.Log("Value should be empty")
	}
}

func TestHashEntry_Reset(t *testing.T) {
	e := NewHashEntry()
	e.flags = USED
	e.key = "bob"
	e.value = "test123"
	if e.flags != EMPTY {
		t.Log("Entry was not marked empty")
	}
	if e.key != "bob" {
		t.Logf("Key should be %s, found: %s", "bob", e.key)
	}
	switch e.value.(type) {
	case string:
		if e.value.(string) != "test123" {
			t.Logf("Value should be '%s', found: %s", "test123", e.value.(string))
		}
	default:
		t.Log("Type should be string")
	}

	e.Reset()
	if e.flags != EMPTY {
		t.Log("Entry was not marked empty")
	}
	if e.key != "" {
		t.Logf("Key should be empty, found: %s", e.key)
	}
	if e.value != nil {
		t.Log("Value should be empty")
	}
}
