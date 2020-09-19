package factions

import uuid "github.com/satori/go.uuid"

// Repository is the interface through which factions may be read and written
type Repository interface {
	Reader
	Writer
}

// Reader is the interface through which factions may be read
type Reader interface {
	FindOne(query interface{}, args ...interface{}) (Faction, error)
}

// Writer is the interface through which factions may be written
type Writer interface {
	CreateOne(data CreationInput) (uuid.UUID, error)
}
