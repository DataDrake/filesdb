package filesdb

import (
	"errors"
	"io"
)

type FileRecord struct {
	name     string
	children Tree
}

func (r *FileRecord) ToCBOR(o io.Writer) (err error) {
	//start record, store name, start child record
	_, err = o.Write([]byte{FILE_RECORD, INF_ARRAY})
	if err != nil {
		return
	}
	_, err = o.Write([]byte(r.name))
	if err != nil {
		return
	}
	_, err = o.Write([]byte{BREAK})
	if err != nil {
		return
	}
	err = r.children.ToCBOR(o)
	return
}

func (r *FileRecord) FromCBOR(i io.Reader) (err error) {
	b := make([]byte, 1)
	// try to read start of record
	_, err = i.Read(b)
	if err != nil {
		return
	}
	if b[0] == BREAK {
		err = errors.New("FileRecord is over")
		return
	}
	if b[0] != FILE_RECORD {
		err = errors.New("Not a filerecord")
		return
	}
	//try to read start of name
	_, err = i.Read(b)
	if err != nil {
		return
	}
	if b[0] != INF_ARRAY {
		err = errors.New("Missing name in file record")
		return
	}
	// try to read name
	_, err = i.Read(b)
	for err == nil && b[0] != BREAK {
		r.name += string(b)
		_, err = i.Read(b)
	}
	if err != nil {
		return
	}
	// read the children
	r.children, err = ReadTreeFromCBOR(i)
	if err != nil {
		return
	}
	return
}
