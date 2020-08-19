package schema

import uuid "github.com/satori/go.uuid"

// BattleCommanderFaction is used to store under which faction each commander fought in each battle.
// Not all commanders were always part of the same faction, and due to the way data was obtained,
// in some cases it is not clear to which faction a commander belonged in a specific battle. This
// struct defines the SQL schema
type BattleCommanderFaction struct {
	BattleID    uuid.UUID `gorm:"type:uuid;primary_key;not null"`
	CommanderID uuid.UUID `gorm:"type:uuid;primary_key;not null"`
	FactionID   uuid.UUID `gorm:"type:uuid;primary_key;not null"`
}
