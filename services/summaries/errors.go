package summaries

import (
	"github.com/sasalatart/batcoms/domain"
)

// ErrNoSummary is used to communicate that no summary could be obtained from Wikipedia's API
const ErrNoSummary = domain.Error("No summary found")

// ErrNotWiki is used to communicate that the URL did not correspond to a Wiki page
const ErrNotWiki = domain.Error("Not a Wiki page")
