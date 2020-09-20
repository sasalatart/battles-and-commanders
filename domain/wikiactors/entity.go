package wikiactors

// Actor stores the details of an entity that participated in a battle as scraped from
// Wikipedia. These can be factions or commanders
type Actor struct {
	Kind        Kind   `validate:"required"`
	ID          int    `validate:"required,min=1"`
	URL         string `validate:"required,url"`
	Flag        string
	Name        string `validate:"required"`
	Description string
	Extract     string `validate:"required"`
}

// Kind represents the kind of an actor in a battle
type Kind int

const (
	// FactionKind represents a faction actor
	FactionKind Kind = iota + 1
	// CommanderKind represents a commander actor
	CommanderKind
)
