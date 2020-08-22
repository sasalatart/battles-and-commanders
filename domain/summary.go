package domain

// Summary stores information retrieved from querying Wikipedia's summary API
type Summary struct {
	PageID       int    `json:"pageid"`
	Type         string `json:"type"`
	Title        string `json:"title"`
	DisplayTitle string `json:"displaytitle"`
	Description  string `json:"description"`
	Extract      string `json:"extract"`
}