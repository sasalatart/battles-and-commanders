package locations

// Location represents a place and coordinates where a battle took place
type Location struct {
	Place     string `validate:"required" json:"place"`
	Latitude  string `validate:"required_with=longitude" json:"latitude"`
	Longitude string `validate:"required_with=latitude" json:"longitude"`
}
