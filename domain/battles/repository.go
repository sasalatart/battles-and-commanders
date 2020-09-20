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
	FindOne(query Battle) (Battle, error)
}

// Writer is the interface through which battles may be written
type Writer interface {
	CreateOne(data CreationInput) (uuid.UUID, error)
}
