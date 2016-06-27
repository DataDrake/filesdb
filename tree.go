package filesdb

import (
	"bytes"
	"errors"
	"strings"
)

type Tree []*FileRecord

func NewTree() *Tree {
	return &make(Tree, 0)
}

func (t *Tree) Add(s string, delim rune) {
	i := strings.IndexRune(s, delim)
	// terminate the tree
	if len(s) == 0 {
		return
	}
	var car string
	var cdr string
	if i == -1 {
		car = s
		cdr = ""
	} else {
		car = s[0:i]
		cdr = s[i+1 : len(s)]
		// skip rune at the beginning of a file path
		if len(car) == 0 {
			i = strings.IndexRune(cdr, delim)
			car = cdr[0:i]
			cdr = cdr[i+1 : len(cdr)]
		}
	}

	// get matching file record
	r := t.FindAtCurrentLevel(car)
	if r == nil {
		// create new if not found
		r = &FileRecord{car, make(Tree, 0)}
		t = append(t, r)
	}

	//recurse to the next level
	r.children.Add(cdr, delim)
	return
}

func (t *Tree) ToCBOR() []bytes {
	cbor := make([]bytes, 0)
	cbor = append(cbor, INF_ARRAY)
	for _, v := range t {
		cbor = append(cbor, v.ToCBOR())
	}
	cbor = append(cbor, BREAK)
	return cbor
}

func (t *Tree) FromCBOR(cbor *bytes.Buffer) (err error) {
	b, err := cbor.ReadByte()
	if err != nil {
		return
	}
	if b != INF_ARRAY {
		// rewind
		cbor.UnreadByte()
		err = errors.New("Not an array for filerecords")
		return
	}
	// read records
	done := false
	for !done {
		r := &FileRecord{}
		done, err = r.FromCBOR(cbor)
		if err != nil {
			return
		}
		t = append(t, r)
	}

	// try to read terminator
	b, err = cbor.ReadByte()
	if err != nil {
		return
	}
	if b != BREAK {
		err = errors.New("Record listing was not correctly terminated")
	}
	return
}
