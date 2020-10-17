package statistics

// SideNumbers stores numerical "raw" (text) data and statistics grouped into each side of a battle
type SideNumbers struct {
	A  string `json:"a"`
	B  string `json:"b"`
	AB string `json:"ab"`
}
