package domain

// ScrapedSideParticipants groups participants' WikiIDs into each side of a battle. The struct may
// be used for either factions or commanders
type ScrapedSideParticipants struct {
	A []int `validate:"unique"`
	B []int `validate:"unique"`
}

// ScrapedCommandersByFaction is a map that groups all of the WikiIDs of commanders that
// participated in a specific battle into their corresponding faction WikiIDs
type ScrapedCommandersByFaction map[int][]int

// SBattle stores all the details regarding a specific battle as scraped from Wikipedia
type SBattle struct {
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
	Factions            ScrapedSideParticipants
	Commanders          ScrapedSideParticipants
	CommandersByFaction ScrapedCommandersByFaction
}
