package wikiactors

// Repository is the interface through which WikiActors may be read and written
type Repository interface {
	Reader
	Writer
}

// Reader is the interface through which WikiActors may be read
type Reader interface {
	Find(kind Kind, id int) *Actor
	FindByURL(kind Kind, url string) *Actor
}

// Writer is the interface through which WikiActors may be written
type Writer interface {
	Save(a Actor) error
	Export(fileName string) error
}
