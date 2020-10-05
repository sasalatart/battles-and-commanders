package factions

import uuid "github.com/satori/go.uuid"

// Faction is an organization to which the commanders and other units involved in a battle belong.
// These may be countries, kingdoms, empires, or other similar entities depending on the historical
// context of the period of time
type Faction struct {
	ID      uuid.UUID
	WikiID  int
	URL     string
	Name    string
	Summary string
}
