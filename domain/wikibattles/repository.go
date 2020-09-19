package wikibattles

// Repository is the interface through which WikiBattles may be read and written
type Repository interface {
	Reader
	Writer
}

// Reader is the interface through which WikiBattles may be read
type Reader interface {
	Find(id int) *Battle
}

// Writer is the interface through which WikiBattles may be written
type Writer interface {
	Save(b Battle) error
	Export(fileName string) error
}
