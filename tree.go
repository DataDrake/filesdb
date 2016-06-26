package filesdb

import "strings"

type Branch struct {
	Prev int
	Next int
}

type Tree struct {
	data []*HashTable
}

func NewTree() *Tree {
	return &Tree{make([]*HashTable, 0)}
}

func (t *Tree) AddItem(s string, delim rune, level int, prev int) int {
	i := strings.IndexRune(s, delim)
	// terminate the tree
	if len(s) == 0 {
		return -1
	}
	//add level if there is not one this deep
	if len(t.data) == level {
		t.data = append(t.data, NewHashTable(-1, nil))
	}
	//get current path piece
	car := s[0:i]
	//try to get existing
	curr, e := t.data[level].Get(car)
	if curr == -1 {
		//make a new one if not found
		curr, e = t.data[level].Put(car, make([]Branch, 0))
	}
	//get the remainder of the path
	cdr := s[i+1 : len(s)]
	//recurse to the next level
	next := t.AddItem(cdr, delim, level+1, curr)
	//save this new branch
	e.value = append(e.value.([]*Branch), &Branch{prev, next})
	return curr
}
