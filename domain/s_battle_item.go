package domain

// SBattleItem stores the name and URL of a battle after being scraped from Wikipedia's indexed list
// of battles
type SBattleItem struct {
	Name string `validate:"required"`
	URL  string `validate:"required,url"`
}
