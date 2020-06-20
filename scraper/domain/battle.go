package domain

// BattleItem stores the name and URL of a battle after being scraped from Wikipedia's indexed list
// of battles
type BattleItem struct {
	Name string
	URL  string
}

// Battle stores all the details regarding a specific battle as scraped from Wikipedia
type Battle struct {
	ID          int
	URL         string
	Name        string
	PartOf      string
	Description string
	Extract     string
	Date        string
	Location    struct {
		Place     string
		Latitude  string
		Longitude string
	}
	Result             string
	TerritorialChanges string
	Strength           struct {
		A string
		B string
	}
	Casualties struct {
		A string
		B string
	}
	Factions struct {
		A []int
		B []int
	}
	Commanders struct {
		A []int
		B []int
	}
	CommandersByFaction map[int][]int
}
