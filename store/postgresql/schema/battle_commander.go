package schema

import uuid "github.com/satori/go.uuid"

// BattleCommander is used to store which commanders were part of a battle, and in which side they
// fought. This struct defines the SQL schema
type BattleCommander struct {
	BattleID    uuid.UUID `gorm:"type:uuid;not null;primary_key"`
	CommanderID uuid.UUID `gorm:"type:uuid;not null;primary_key"`
	Side        SideKind  `gorm:"type:integer;not null;index"`
	Battle      *Battle
	Commander   *Commander
}
