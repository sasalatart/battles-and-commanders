package schema

// Commander is used to store data that defines a specific commander. This struct defines the SQL
// schema
type Commander struct {
	Base
	WikiID   int    `gorm:"not null;index"`
	URL      string `gorm:"unique_index;not null"`
	Name     string `gorm:"not null"`
	Summary  string
	Battles  []Battle  `gorm:"many2many:battle_commanders;"`
	Factions []Faction `gorm:"many2many:battle_commander_factions;"`
}
