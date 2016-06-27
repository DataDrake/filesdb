package filesdb

import (
	"bytes"
	"errors"
)

type FileRecord struct {
	name     string
	children *Tree
}

func (r *FileRecord) ToCBOR() []byte {
	result := make([]byte, 0)
	//start record, store name, start child record
	result = append(result, []byte{FILE_RECORD, INF_ARRAY, r.name, BREAK, INF_ARRAY})
	result = append(result, r.children.ToCBOR())
	result = append(result, BREAK)
	return result
}

func (r *FileRecord) FromCBOR(cbor *bytes.Buffer) (done bool, err error) {
	// try to read start of record
	b, err := cbor.ReadByte()
	if err != nil {
		return
	}
	if b == BREAK {
		goto SUCCESS
	}
	if b != FILE_RECORD {
		err = errors.New("Not a filerecord")
		goto FAILURE
	}
	//try to read start of name
	b, err = cbor.ReadByte()
	if err != nil {
		return
	}
	if b != INF_ARRAY {
		err = errors.New("Missing name in file record")
		goto FAILURE
	}
	// try to read name
	b, err = cbor.ReadByte()
	for err == nil && b != BREAK {
		r.name = append(r.name, b)
		b, err = cbor.ReadByte()
	}
	if err != nil {
		return
	}
	// read the children
	done, err = r.children.FromCBOR(cbor)
	if err != nil {
		return
	}
	if !done {
		err = errors.New("Children not terminated correctly")
		return
	}
	goto SUCCESS
FAILURE:
	// rewind
	cbor.UnreadByte()
	return

SUCCESS:
	done = true
	return
}
