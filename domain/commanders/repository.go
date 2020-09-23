package commanders

import uuid "github.com/satori/go.uuid"

// Repository is the interface through which commanders may be read and written
type Repository interface {
	Reader
	Writer
}

// Reader is the interface through which commanders may be read
type Reader interface {
	FindOne(query Commander) (Commander, error)
	FindMany(query Query, page uint) ([]Commander, uint, error)
}

// Writer is the interface through which commanders may be written
type Writer interface {
	CreateOne(data CreationInput) (uuid.UUID, error)
}

// Query is used to refine the filters when finding many commanders
type Query struct {
	FactionID uuid.UUID
	Name      string
	Summary   string
}