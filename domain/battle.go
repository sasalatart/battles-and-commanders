package domain

import uuid "github.com/satori/go.uuid"

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
	Location            Location
	Result              string
	TerritorialChanges  string
	Strength            SideNumbers
	Casualties          SideNumbers
	Factions            FactionsBySide
	Commanders          CommandersBySide
	CommandersByFaction CommandersByFaction
}

// FactionsBySide groups all of the factions that participated in a battle into the two opposing
// sides. These sides are A and B
type FactionsBySide struct {
	A []Faction
	B []Faction
}

// CommandersBySide groups all of the commanders that participated in a battle into the two opposing
// sides. These sides are A and B
type CommandersBySide struct {
	A []Commander
	B []Commander
}

// ParticipantsIDsBySide groups the IDs of the participants of a specific kind (either factions or
// commanders) that participated in a battle into the two opposing sides. These sides are A and B
type ParticipantsIDsBySide struct {
	A []uuid.UUID
	B []uuid.UUID
}

// CommandersByFaction is a map that indexes the IDs of commanders according to the faction to which
// they belong during a specific battle
type CommandersByFaction map[uuid.UUID][]uuid.UUID

// CreateBattleInput is a struct that contains all of the data required to create a battle. This
// includes annotations required by validations
type CreateBattleInput struct {
	WikiID              int    `validate:"required"`
	URL                 string `validate:"required,url"`
	Name                string `validate:"required"`
	PartOf              string
	Summary             string `validate:"required"`
	StartDate           string `validate:"required"`
	EndDate             string `validate:"required"`
	Location            Location
	Result              string `validate:"required"`
	TerritorialChanges  string
	Strength            SideNumbers
	Casualties          SideNumbers
	FactionsBySide      ParticipantsIDsBySide
	CommandersBySide    ParticipantsIDsBySide
	CommandersByFaction CommandersByFaction
}
