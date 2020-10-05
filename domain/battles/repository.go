package battles

import (
	"github.com/sasalatart/batcoms/domain/locations"
	"github.com/sasalatart/batcoms/domain/statistics"
	"github.com/sasalatart/batcoms/pkg/dates"
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
	FromDate    dates.Historic
	ToDate      dates.Historic
}

// CreationInput is a struct that contains all of the data required to create a battle. This
// includes annotations required by validations
type CreationInput struct {
	WikiID              int    `validate:"required"`
	URL                 string `validate:"required,url"`
	Name                string `validate:"required"`
	PartOf              string
	Summary             string         `validate:"required"`
	StartDate           dates.Historic `validate:"required"`
	EndDate             dates.Historic `validate:"required"`
	Location            locations.Location
	Result              string `validate:"required"`
	TerritorialChanges  string
	Strength            statistics.SideNumbers
	Casualties          statistics.SideNumbers
	FactionsBySide      IDsBySide
	CommandersBySide    IDsBySide
	CommandersByFaction CommandersByFaction
}
