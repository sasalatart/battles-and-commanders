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

// CreationInput is a struct that contains all of the data required to create a commander. This
// includes annotations required by validations
type CreationInput struct {
	WikiID  int    `validate:"required"`
	URL     string `validate:"required,url"`
	Name    string `validate:"required"`
	Summary string
}
