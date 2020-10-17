package battles

import (
	"github.com/sasalatart/batcoms/domain/commanders"
	"github.com/sasalatart/batcoms/domain/factions"
	"github.com/sasalatart/batcoms/domain/locations"
	"github.com/sasalatart/batcoms/domain/statistics"
	"github.com/sasalatart/batcoms/pkg/dates"
	uuid "github.com/satori/go.uuid"
)

// Battle represents an armed encounter between two or more commanders belonging to different
// factions that occurred at some point, somewhere in history
type Battle struct {
	ID                  uuid.UUID              `json:"id"`
	WikiID              int                    `json:"wikiID"`
	URL                 string                 `json:"url"`
	Name                string                 `json:"name"`
	PartOf              string                 `json:"partOf"`
	Summary             string                 `json:"summary"`
	StartDate           dates.Historic         `json:"startDate"`
	EndDate             dates.Historic         `json:"endDate"`
	Location            locations.Location     `json:"location"`
	Result              string                 `json:"result"`
	TerritorialChanges  string                 `json:"territorialChanges"`
	Strength            statistics.SideNumbers `json:"strength"`
	Casualties          statistics.SideNumbers `json:"casualties"`
	Factions            FactionsBySide         `json:"factions"`
	Commanders          CommandersBySide       `json:"commanders"`
	CommandersByFaction CommandersByFaction    `json:"commandersByFaction"`
}

// FactionsBySide groups all of the factions that participated in a battle into the two opposing
// sides. These sides are A and B
type FactionsBySide struct {
	A []factions.Faction `json:"a"`
	B []factions.Faction `json:"b"`
}

// CommandersBySide groups all of the commanders that participated in a battle into the two opposing
// sides. These sides are A and B
type CommandersBySide struct {
	A []commanders.Commander `json:"a"`
	B []commanders.Commander `json:"b"`
}

// IDsBySide groups the IDs of the factions or commanders that participated in a battle into the two
// opposing sides. These sides are A and B
type IDsBySide struct {
	A []uuid.UUID `json:"a"`
	B []uuid.UUID `json:"b"`
}

// CommandersByFaction is a map that indexes the IDs of commanders according to the faction to which
// they belong during a specific battle
type CommandersByFaction map[uuid.UUID][]uuid.UUID
