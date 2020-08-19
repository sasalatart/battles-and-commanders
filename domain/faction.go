package domain

import uuid "github.com/satori/go.uuid"

// Faction a socio-political organization to which the participants of historical battles belong
type Faction struct {
	ID      uuid.UUID
	WikiID  int
	URL     string
	Name    string
	Summary string
}

// CreateFactionInput is a struct that contains all of the data required to create a faction. This
// includes annotations required by validations
type CreateFactionInput struct {
	WikiID  int    `validate:"required"`
	URL     string `validate:"required,url"`
	Name    string `validate:"required"`
	Summary string
}
