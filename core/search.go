package filesdb

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"io"
	"bytes"
)

func (t *Tree) FindAtCurrentLevel(name string) *FileRecord {
	for _, v := range *t {
		if v.name == name {
			return v
		}
	}
	return nil
}

func SearchRecurse(query string, path string , db io.ByteReader) {
	//try to read start of filearray
	b, err := db.ReadByte()
	if err != nil {
		panic(err.Error())
	}
	if b != INF_ARRAY {
		panic("Missing children in file record")
	}
	for {
		// try to read start of record
		b, err = db.ReadByte()
		if err != nil {
			return
		}
		if b == BREAK {
			break
		}
		if b != FILE_RECORD {
			panic("Not a filerecord")
		}
		//try to read start of name
		b, err = db.ReadByte()
		if err != nil {
			panic(err.Error())
		}
		if b != INF_ARRAY {
			panic("Missing name in file record")
		}
		// try to read name
		b, err =db.ReadByte()
		name := bytes.NewBuffer(make([]byte,0))
		for err == nil && b != BREAK {
			name.WriteByte(b)
			b, err = db.ReadByte()
		}
		if err != nil {
			panic(err.Error())
		}
		if strings.Contains(name.String(), query) || strings.Contains(path, query){
			fmt.Println(path + name.String())
		}
		// read the children
		SearchRecurse(query, path+name.String()+string(filepath.Separator),db)
	}
}

func Search(name string, db io.ByteReader) {
	if strings.ContainsRune(name, '/') {
		fmt.Println("ERROR: search string should be a filename, not a path")
		os.Exit(1)
	}
	SearchRecurse(name, "/", db)
}

func (t *Tree) Fill(path string) {
	i, err := os.Stat(path)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if !i.IsDir() {
		fmt.Println("Path should be a directory, not a file")
		return
	}
	fs, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, f := range fs {
		r := &FileRecord{f.Name(), NewTree()}
		if f.IsDir() {
			r.children.Fill(filepath.Join(path, f.Name()))
		}
		*t = append(*t, r)
	}
}
