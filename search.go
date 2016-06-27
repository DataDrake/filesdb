package filesdb

func (t *Tree) FindAtCurrentLevel(name string) *FileRecord {
	for _, v := range t {
		if v.name == name {
			return v
		}
	}
	return nil
}