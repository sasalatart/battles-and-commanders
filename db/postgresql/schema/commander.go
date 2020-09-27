package schema

// Commander is used to store data that defines a specific commander. This struct defines the SQL
// schema
type Commander struct {
	Base
	WikiID  int    `gorm:"not null;uniqueIndex"`
	URL     string `gorm:"not null;uniqueIndex"`
	Name    string `gorm:"not null;index"`
	Summary string `gorm:"not null"`
}
