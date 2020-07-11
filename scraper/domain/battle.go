package domain

// BattleItem stores the name and URL of a battle after being scraped from Wikipedia's indexed list
// of battles
type BattleItem struct {
	Name string `validate:"required"`
	URL  string `validate:"required,url"`
}

// Location represents a place where a battle took place
type Location struct {
	Place     string `validate:"required"`
	Latitude  string `validate:"required_with=longitude"`
	Longitude string `validate:"required_with=latitude"`
}

// SideNumbers stores numerical "raw" (text) data and statistics grouped into each side of a battle
type SideNumbers struct {
	A  string
	B  string
	AB string
}

// SideParticipants stores participants IDs grouped into each side of a battle. The struct may be
// used for either factions or commanders
type SideParticipants struct {
	A []int `validate:"unique"`
	B []int `validate:"unique"`
}

// Battle stores all the details regarding a specific battle as scraped from Wikipedia
type Battle struct {
	ID                  int    `validate:"required,min=1"`
	URL                 string `validate:"required,url"`
	Name                string `validate:"required"`
	PartOf              string
	Description         string
	Extract             string `validate:"required"`
	Date                string `validate:"required"`
	Location            Location
	Result              string `validate:"required"`
	TerritorialChanges  string
	Strength            SideNumbers
	Casualties          SideNumbers
	Factions            SideParticipants
	Commanders          SideParticipants
	CommandersByFaction map[int][]int
}
