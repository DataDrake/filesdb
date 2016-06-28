package filesdb

import (
	"bytes"
	"testing"
)

var BYTES1 = []byte{FILE_RECORD, INF_ARRAY, 0x62, 0x6f, 0x62, BREAK, INF_ARRAY, BREAK}
var BYTES2 = []byte{FILE_RECORD, INF_ARRAY, 0x62, 0x6f, 0x62, BREAK,
	INF_ARRAY, FILE_RECORD, INF_ARRAY, 0x73, 0x61, 0x6c, 0x6c, 0x79, BREAK, INF_ARRAY, BREAK, BREAK}

func TestFileRecord_ToCBOR(t *testing.T) {
	r := &FileRecord{"bob", make(Tree, 0)}
	b := bytes.NewBuffer(make([]byte, 0))
	err := r.ToCBOR(b)
	if err != nil {
		t.Errorf("Failed to serialize record, reason: %s", err.Error())
	}
	if !bytes.Equal(BYTES1, b.Bytes()) {
		t.Errorf("Bytes did not match, found: %x", b.Bytes())
	}

	b = bytes.NewBuffer(make([]byte, 0))
	r2 := &FileRecord{"sally", make(Tree, 0)}
	r.children = append(r.children, r2)
	err = r.ToCBOR(b)
	if err != nil {
		t.Errorf("Failed to serialize record, reason: %s", err.Error())
	}
	if !bytes.Equal(BYTES2, b.Bytes()) {
		t.Errorf("Bytes did not match, found: %x", b.Bytes())
	}
}

func TestFileRecord_FromCBOR(t *testing.T) {
	b := bytes.NewBuffer(BYTES1)
	r := &FileRecord{}
	err := r.FromCBOR(b)
	if err != nil {
		t.Errorf("Failed to deserialize record, reason: %s", err.Error())
	}
	if r.name != "bob" {
		t.Errorf("Name should be 'bob', found: %s", r.name)
	}
	if len(r.children) != 0 {
		t.Error("Should not have any children")
	}
}

func TestFileRecord_FromCBOR2(t *testing.T) {
	b := bytes.NewBuffer(BYTES2)
	r := &FileRecord{}
	err := r.FromCBOR(b)
	if err != nil {
		t.Errorf("Failed to deserialize record, reason: %s", err.Error())
	}
	if r.name != "bob" {
		t.Errorf("Name should be 'bob', found: %s", r.name)
	}
	if len(r.children) != 1 {
		t.Error("Should have one child")
	}
	if r.children[0].name != "sally" {
		t.Errorf("Name should be 'sally', found: %s", r.children[0].name)
	}
	if len(r.children[0].children) != 0 {
		t.Error("Should have no children")
	}
}
