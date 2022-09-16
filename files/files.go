package files

type Entry struct {
	Name string
	File *File
	Dir  *Dir
}

type File []byte

type Dir []Entry

type WalkFunc func(path string, d Entry, err error) error

func WalkDir(root Dir, call WalkFunc) error {
	XXX
}

func WriteDir(root Dir, path string) error {
	XXX
}
