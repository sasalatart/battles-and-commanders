package schema

// Faction is used to store data that defines a specific faction. This struct defines the SQL schema
type Faction struct {
	Base
	WikiID     int    `gorm:"not null;index"`
	URL        string `gorm:"unique_index;not null"`
	Name       string `gorm:"not null"`
	Summary    string
	Battles    []Battle    `gorm:"many2many:battle_commanders;"`
	Commanders []Commander `gorm:"many2many:battle_commander_factions;"`
}
