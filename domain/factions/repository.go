package factions

import uuid "github.com/satori/go.uuid"

// Repository is the interface through which factions may be read and written
type Repository interface {
	Reader
	Writer
}

// Reader is the interface through which factions may be read
type Reader interface {
	FindOne(query FindOneQuery) (Faction, error)
	FindMany(query FindManyQuery, page int) ([]Faction, int, error)
}

// Writer is the interface through which factions may be written
type Writer interface {
	CreateOne(data CreationInput) (uuid.UUID, error)
}

// FindOneQuery is used to refine the filters when finding one faction
type FindOneQuery struct {
	ID   uuid.UUID
	Name string
	URL  string
}

// FindManyQuery is used to refine the filters when finding many factions
type FindManyQuery struct {
	CommanderID uuid.UUID
	Name        string
	Summary     string
}

// CreationInput is a struct that contains all of the data required to create a faction. This
// includes annotations required by validations
type CreationInput struct {
	WikiID  int    `validate:"required"`
	URL     string `validate:"required,url"`
	Name    string `validate:"required"`
	Summary string
}
