package schema

import (
	uuid "github.com/satori/go.uuid"
)

// Base contains common columns for all tables.
type Base struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
}
