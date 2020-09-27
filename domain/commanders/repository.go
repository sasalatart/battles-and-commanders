package commanders

import uuid "github.com/satori/go.uuid"

// Repository is the interface through which commanders may be read and written
type Repository interface {
	Reader
	Writer
}

// Reader is the interface through which commanders may be read
type Reader interface {
	FindOne(query FindOneQuery) (Commander, error)
	FindMany(query FindManyQuery, page int) ([]Commander, int, error)
}

// Writer is the interface through which commanders may be written
type Writer interface {
	CreateOne(data CreationInput) (uuid.UUID, error)
}

// FindOneQuery is used to refine the filters when finding one commander
type FindOneQuery struct {
	ID   uuid.UUID
	Name string
	URL  string
}

// FindManyQuery is used to refine the filters when finding many commanders
type FindManyQuery struct {
	FactionID uuid.UUID
	Name      string
	Summary   string
}
