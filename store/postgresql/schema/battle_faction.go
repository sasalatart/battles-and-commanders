package schema

import uuid "github.com/satori/go.uuid"

// BattleFaction is used to store which factions were part of a battle, and in which side they
// fought. This struct defines the SQL schema
type BattleFaction struct {
	BattleID  uuid.UUID `gorm:"type:uuid;primary_key;not null"`
	FactionID uuid.UUID `gorm:"type:uuid;primary_key;not null"`
	Side      SideKind  `gorm:"type:integer;index;not null"`
	Battle    *Battle
	Faction   *Faction
}
