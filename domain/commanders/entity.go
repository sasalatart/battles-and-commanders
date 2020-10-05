package commanders

import uuid "github.com/satori/go.uuid"

// Commander represents a military leader involved in one or more battles in history
type Commander struct {
	ID      uuid.UUID
	WikiID  int
	URL     string
	Name    string
	Summary string
}
