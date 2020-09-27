package schema

import uuid "github.com/satori/go.uuid"

// BattleFaction is used to store which factions were part of a battle, and in which side they
// fought. This struct defines the SQL schema
type BattleFaction struct {
	BattleID  uuid.UUID `gorm:"type:uuid;primaryKey;not null"`
	FactionID uuid.UUID `gorm:"type:uuid;primaryKey;not null"`
	Side      SideKind  `gorm:"type:integer;primaryKey;not null"`
	Battle    Battle
	Faction   Faction
}
