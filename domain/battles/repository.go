package battles

import (
	uuid "github.com/satori/go.uuid"
)

// Repository is the interface through which battles may be read and written
type Repository interface {
	Reader
	Writer
}

// Reader is the interface through which battles may be read
type Reader interface {
	FindOne(query FindOneQuery) (Battle, error)
	FindMany(query FindManyQuery, page int) ([]Battle, int, error)
}

// Writer is the interface through which battles may be written
type Writer interface {
	CreateOne(data CreationInput) (uuid.UUID, error)
}

// FindOneQuery is used to refine the filters when finding one battle
type FindOneQuery struct {
	ID   uuid.UUID
	Name string
	URL  string
}

// FindManyQuery is used to refine the filters when finding many battles
type FindManyQuery struct {
	FactionID   uuid.UUID
	CommanderID uuid.UUID
	Name        string
	Summary     string
	Place       string
	Result      string
}
