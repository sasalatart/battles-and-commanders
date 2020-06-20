package domain

// ParticipantKind represents the kind of a participant in a battle
type ParticipantKind int

const (
	// FactionKind represents a faction participant
	FactionKind ParticipantKind = iota
	// CommanderKind represents a commander participant
	CommanderKind
)

// Participant stores the details of an entity that participated in a battle. These can either be a
// faction or a commander
type Participant struct {
	Kind        ParticipantKind
	ID          int
	URL         string
	Flag        string
	Name        string
	Description string
	Extract     string
}
