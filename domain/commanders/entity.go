package commanders

import uuid "github.com/satori/go.uuid"

// Commander represents a military leader involved in one or more battles in history
type Commander struct {
	ID      uuid.UUID `json:"id"`
	WikiID  int       `json:"wikiID"`
	URL     string    `json:"url"`
	Name    string    `json:"name"`
	Summary string    `json:"summary"`
}
