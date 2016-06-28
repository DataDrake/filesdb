package filesdb

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func (t *Tree) FindAtCurrentLevel(name string) *FileRecord {
	for _, v := range *t {
		if v.name == name {
			return v
		}
	}
	return nil
}

func (t *Tree) SearchRecurse(name string, path string) {
	for _, v := range *t {
		if strings.Contains(v.name, name) {
			fmt.Println(path + v.name)
		}
		v.children.SearchRecurse(name, path+v.name+string(filepath.Separator))
	}
}

func (t *Tree) Search(name string) {
	if strings.ContainsRune(name, '/') {
		fmt.Println("ERROR: search string should be a filename, not a path")
		os.Exit(1)
	}
	t.SearchRecurse(name, "/")
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
