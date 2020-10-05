package schema

import "gorm.io/datatypes"

// Battle is used to store data that defines a specific battle. This struct defines the SQL schema
type Battle struct {
	Base
	WikiID                  int    `gorm:"not null;uniqueIndex"`
	URL                     string `gorm:"not null;uniqueIndex"`
	Name                    string `gorm:"not null;uniqueIndex"`
	PartOf                  string
	Summary                 string  `gorm:"not null"`
	StartDate               string  `gorm:"not null;index"`
	StartDateNum            float64 `gorm:"not null;index"`
	EndDate                 string  `gorm:"not null;index"`
	EndDateNum              float64 `gorm:"not null;index"`
	Place                   string  `gorm:"not null"`
	Latitude                string
	Longitude               string
	Result                  string `gorm:"not null"`
	TerritorialChanges      string
	Strength                datatypes.JSON
	Casualties              datatypes.JSON
	BattleCommanders        []BattleCommander
	BattleFactions          []BattleFaction
	BattleCommanderFactions []BattleCommanderFaction
}

// SideKind is the type used to represent all of the sides that participated in a battle. These can
// be SideA, SideB, and SideAB. This separation is completely unrelated to who emerged victorious
// after the battle, or who "was in the right side of history"
type SideKind int64

const (
	// SideA represents one of the two possible opposing sides of a battle
	SideA SideKind = iota
	// SideB represents one of the two possible opposing sides of a battle
	SideB
	// SideAB represent a fallback for data that may not be assigned to a specific side of a battle
	SideAB
)
