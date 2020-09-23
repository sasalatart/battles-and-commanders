package factions

import uuid "github.com/satori/go.uuid"

// Repository is the interface through which factions may be read and written
type Repository interface {
	Reader
	Writer
}

// Reader is the interface through which factions may be read
type Reader interface {
	FindOne(query Faction) (Faction, error)
	FindMany(query Query, page uint) ([]Faction, uint, error)
}

// Writer is the interface through which factions may be written
type Writer interface {
	CreateOne(data CreationInput) (uuid.UUID, error)
}

// Query is used to refine the filters when finding many factions
type Query struct {
	CommanderID uuid.UUID
	Name        string
	Summary     string
}
