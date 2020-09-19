package summaries

import "github.com/sasalatart/batcoms/domain"

// Summary stores information retrieved from querying Wikipedia's summary API
type Summary struct {
	PageID      int    `json:"pageid"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Extract     string `json:"extract"`
}

// ErrNoSummary is used to communicate that no summary could be obtained from Wikipedia's API
const ErrNoSummary = domain.Error("No summary found")

// ErrNotWiki is used to communicate that the URL did not correspond to a Wiki page
const ErrNotWiki = domain.Error("Not a Wiki page")
