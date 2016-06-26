package filesdb

const (
	EMPTY = uint8(0)
	USED  = uint8(1)
)

type HashEntry struct {
	flags uint8
	key   string
	value interface{}
}

func NewHashEntry() *HashEntry {
	return &HashEntry{0, "", nil}
}

func (e *HashEntry) Reset() interface{} {
	// save value
	v := e.value
	// reset flags
	e.flags = EMPTY
	// reset key
	e.key = ""
	// reset value
	e.value = nil
	return v
}
