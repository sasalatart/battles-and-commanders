package domain

// ParticipantKind represents the kind of a participant in a battle
type ParticipantKind int

const (
	// FactionKind represents a faction participant
	FactionKind ParticipantKind = iota + 1
	// CommanderKind represents a commander participant
	CommanderKind
)

// SParticipant stores the details of an entity that participated in a battle as scraped from
// Wikipedia. These can be factions or commanders
type SParticipant struct {
	Kind        ParticipantKind `validate:"required"`
	ID          int             `validate:"required,min=1"`
	URL         string          `validate:"required,url"`
	Flag        string
	Name        string `validate:"required"`
	Description string
	Extract     string
}