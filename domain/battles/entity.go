package battles

import (
	"github.com/sasalatart/batcoms/domain/commanders"
	"github.com/sasalatart/batcoms/domain/factions"
	"github.com/sasalatart/batcoms/domain/locations"
	"github.com/sasalatart/batcoms/domain/statistics"
	uuid "github.com/satori/go.uuid"
)

// Battle represents an armed encounter between two or more commanders belonging to different
// factions that occurred at some point, somewhere in history
type Battle struct {
	ID                  uuid.UUID
	WikiID              int
	URL                 string
	Name                string
	PartOf              string
	Summary             string
	StartDate           string
	EndDate             string
	Location            locations.Location
	Result              string
	TerritorialChanges  string
	Strength            statistics.SideNumbers
	Casualties          statistics.SideNumbers
	Factions            FactionsBySide
	Commanders          CommandersBySide
	CommandersByFaction CommandersByFaction
}

// FactionsBySide groups all of the factions that participated in a battle into the two opposing
// sides. These sides are A and B
type FactionsBySide struct {
	A []factions.Faction
	B []factions.Faction
}

// CommandersBySide groups all of the commanders that participated in a battle into the two opposing
// sides. These sides are A and B
type CommandersBySide struct {
	A []commanders.Commander
	B []commanders.Commander
}

// CreationInput is a struct that contains all of the data required to create a battle. This
// includes annotations required by validations
type CreationInput struct {
	WikiID              int    `validate:"required"`
	URL                 string `validate:"required,url"`
	Name                string `validate:"required"`
	PartOf              string
	Summary             string `validate:"required"`
	StartDate           string `validate:"required"`
	EndDate             string `validate:"required"`
	Location            locations.Location
	Result              string `validate:"required"`
	TerritorialChanges  string
	Strength            statistics.SideNumbers
	Casualties          statistics.SideNumbers
	FactionsBySide      IDsBySide
	CommandersBySide    IDsBySide
	CommandersByFaction CommandersByFaction
}

// IDsBySide groups the IDs of the factions or commanders that participated in a battle into the two
// opposing sides. These sides are A and B
type IDsBySide struct {
	A []uuid.UUID
	B []uuid.UUID
}

// CommandersByFaction is a map that indexes the IDs of commanders according to the faction to which
// they belong during a specific battle
type CommandersByFaction map[uuid.UUID][]uuid.UUID
