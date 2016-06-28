package filesdb

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/scanner"
)

func SearchRecurse(query string, path string, db io.Reader) {
	b := make([]byte,1)
	//try to read start of filearray
	_, err := db.Read(b)
	if err != nil {
		panic(err.Error())
	}
	if b[0] != INF_ARRAY {
		panic("Missing children in file record")
	}
	for {
		// try to read start of record
		_, err = db.Read(b)
		if err != nil {
			panic(err.Error())
		}
		if b[0] == BREAK {
			break
		}
		if b[0] != FILE_RECORD {
			panic("Not a filerecord")
		}
		//try to read start of name
		_, err = db.Read(b)
		if err != nil {
			panic(err.Error())
		}
		if b[0] != INF_ARRAY {
			panic("Missing name in file record")
		}
		// try to read name
		_, err = db.Read(b)
		name := bytes.NewBuffer(make([]byte, 0))
		for err == nil && b[0] != BREAK{
			name.WriteByte(b[0])
			_, err = db.Read(b)
		}
		if err != nil {
			fmt.Println(filepath.Join(path,name.String()))
			panic(err.Error())
		}
		if strings.Contains(name.String(), query) || strings.Contains(path, query) {
			fmt.Println(path + name.String())
		}
		// read the children
		SearchRecurse(query, path+name.String()+string(filepath.Separator), db)
	}
}

func Search(name string, db io.Reader) {
	if strings.ContainsRune(name, '/') {
		fmt.Println("ERROR: search string should be a filename, not a path")
		os.Exit(1)
	}
	SearchRecurse(name, "/", db)
}

func Fill(path string, db io.Writer) {
	fs, err := ioutil.ReadDir(path)
	if err != nil {
		_, e := db.Write([]byte{INF_ARRAY,BREAK})
		if e != nil {
			panic(e.Error())
		}
		fmt.Fprintf(os.Stderr,"%s\n",err.Error())
		return
	}
	_, err = db.Write([]byte{INF_ARRAY})
	if err != nil {
		panic(err.Error())
	}
	for _, f := range fs {
		name := f.Name()
		r := scanner.Scanner{}
		r.Init(strings.NewReader(name))
		for r.IsValid(){
			curr := r.Next()
			if curr == scanner.EOF {
				panic(err.Error())
			}
		}
		_, err = db.Write([]byte{FILE_RECORD,INF_ARRAY})
		if err != nil {
			panic(err.Error())
		}
		//write name
		_, err = db.Write([]byte(name))
		if err != nil {
			panic(err.Error())
		}
		_, err = db.Write([]byte{BREAK})
		if err != nil {
			panic(err.Error())
		}
		if f.IsDir() {
			Fill(filepath.Join(path, name), db)
		} else {
			_, err = db.Write([]byte{INF_ARRAY,BREAK})
			if err != nil {
				panic(err.Error())
			}
		}
	}
	_, err = db.Write([]byte{BREAK})
	if err != nil {
		panic(err.Error())
	}
}
