package domain

// Location represents a place and coordinates where a battle took place
type Location struct {
	Place     string `validate:"required"`
	Latitude  string `validate:"required_with=longitude"`
	Longitude string `validate:"required_with=latitude"`
}
