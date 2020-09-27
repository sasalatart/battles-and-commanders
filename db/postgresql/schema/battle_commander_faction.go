package schema

import uuid "github.com/satori/go.uuid"

// BattleCommanderFaction is used to store under which faction each commander fought in each battle.
// Not all commanders were always part of the same faction, and due to the way data was obtained,
// in some cases it is not clear to which faction a commander belonged in a specific battle. This
// may result in some commanders being present in battles, but not being assigned a faction. This
// struct defines the SQL schema
type BattleCommanderFaction struct {
	BattleID    uuid.UUID `gorm:"type:uuid;primaryKey;not null"`
	CommanderID uuid.UUID `gorm:"type:uuid;primaryKey;not null"`
	FactionID   uuid.UUID `gorm:"type:uuid;primaryKey;not null"`
	Battle      Battle
	Commander   Commander
	Faction     Faction
}
