package wikibattles

import (
	"github.com/sasalatart/batcoms/domain/locations"
	"github.com/sasalatart/batcoms/domain/statistics"
)

// Battle stores all the details regarding a specific battle as scraped from Wikipedia
type Battle struct {
	ID                  int    `validate:"required,min=1"`
	URL                 string `validate:"required,url"`
	Name                string `validate:"required"`
	PartOf              string
	Description         string
	Extract             string `validate:"required"`
	Date                string `validate:"required"`
	Location            locations.Location
	Result              string `validate:"required"`
	TerritorialChanges  string
	Strength            statistics.SideNumbers
	Casualties          statistics.SideNumbers
	Factions            SideActors
	Commanders          SideActors
	CommandersByFaction CommandersByFaction
}

// SideActors groups actors' WikiIDs into each side of a battle. The struct may be used for either
// factions or commanders
type SideActors struct {
	A []int `validate:"unique"`
	B []int `validate:"unique"`
}

// CommandersByFaction is a map that groups all of the WikiIDs of commanders that participated
// in a specific battle into their corresponding faction WikiIDs
type CommandersByFaction map[int][]int

// BattleItem stores the name and URL of a battle after being scraped from Wikipedia's indexed
// list of battles
type BattleItem struct {
	Name string `validate:"required"`
	URL  string `validate:"required,url"`
}
