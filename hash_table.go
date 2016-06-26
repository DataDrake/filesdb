package filesdb

import (
	"crypto/md5"
	"encoding/binary"
	"hash"
)

var DEFAULT_HASH_SIZE = 16

type HashTable struct {
	size int
	data []*HashEntry
	hash hash.Hash
}

func NewHashTable(size int, hash hash.Hash) *HashTable {
	if size <= 0 {
		size = DEFAULT_HASH_SIZE
	}
	if hash == nil {
		hash = md5.New()
	}
	t := &HashTable{size, make([]*HashEntry, 0), hash}
	for i := 0; i < size; i++ {
		t.data = append(t.data, NewHashEntry())
	}
	return t
}

func get_index(key string, size int, hash hash.Hash) int {
	hash.Reset()
	// hash
	hash.Write([]byte(key))
	// get the hash value
	bindex := hash.Sum(make([]byte, 0))
	// convert upper 4 bytes of hash to int and module by size
	return int(binary.BigEndian.Uint32(bindex[0:3])) % size
}

func (t *HashTable) Get(key string) (int, *HashEntry) {
	index := get_index(key, t.hash, t.size)
	e := t.data[index]
	// if no match, it does not exist
	if e.key != key {
		return -1, nil
	}
	return index, e.value
}

func rehash(t *HashTable) {
	// double the size of the existing hash
	t2 := NewHashTable(2*t.size, t.hash)
	//copy over the values
	for _, v := range t.data {
		index := get_index(v.key, t2.hash, t2.size)
		t2.data[index] = v
	}
	// change pointer
	t = t2
}

func (t *HashTable) Put(key string, value interface{}) (int, *HashEntry) {
	index := get_index(key, t.hash, t.size)
	e := t.data[index]
	//rehash if collision
	if (e.flags&USED) > 0 && e.key != key {
		rehash(t)
	}
	//mark as used
	t.data[index].flags |= USED
	//store the value
	t.data[index].value = value
	return index, t.data[index]
}

func (t *HashTable) Remove(key string) interface{} {
	index := get_index(key, t.hash, t.size)
	e := t.data[index]
	// fail early if unused
	if (e.flags & USED) == 0 {
		return nil
	}
	return e.Reset()
}
