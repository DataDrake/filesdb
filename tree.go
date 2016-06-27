package filesdb

import (
	"errors"
	"io"
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

func (t *Tree) ToCBOR(o io.Writer) (err error) {
	_, err = o.Write([]byte{INF_ARRAY})
	if err != nil {
		return
	}
	for _, v := range t {
		err = v.ToCBOR(o)
		if err != nil {
			return
		}
	}
	_, err = o.Write([]byte{BREAK})
	return
}

func (t *Tree) FromCBOR(i io.Reader) (err error) {
	b := make([]byte,1)
	_,err = i.Read(b)
	if err != nil {
		return
	}
	if b[0] != INF_ARRAY {
		err = errors.New("Not an array for filerecords")
		return
	}
	// read records
	done := false
	for !done {
		r := &FileRecord{}
		done, err = r.FromCBOR(i)
		if err != nil {
			return
		}
		t = append(t, r)
	}

	// try to read terminator
	_, err = i.Read(b)
	if err != nil {
		return
	}
	if b[0] != BREAK {
		err = errors.New("Record listing was not correctly terminated")
	}
	return
}
