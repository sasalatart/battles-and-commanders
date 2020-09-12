package summaries

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/services/parser"
)

// Fetch queries a Wikipedia URL page for its corresponding summary API endpoint, and returns a
// domain.Summary if successful, or an error if not
func Fetch(url string) (domain.Summary, error) {
	summary := domain.Summary{}

	if !strings.Contains(url, "/wiki/") {
		return summary, ErrNotWiki
	}

	summaryURL := strings.ReplaceAll(url, "/wiki/", "/api/rest_v1/page/summary/")
	resp, err := http.Get(summaryURL)
	if err != nil {
		return summary, errors.Wrapf(err, "GET %s failed", summaryURL)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return summary, errors.Wrapf(err, "Reading response from %s", summaryURL)
	}
	if err = json.Unmarshal(body, &summary); err != nil {
		return summary, errors.Wrapf(err, "Unmarshalling response from %s", summaryURL)
	}
	if strings.Contains(strings.ToLower(summary.Title), "not found") || strings.Contains(summary.Type, "not_found") {
		return summary, ErrNoSummary
	}
	summary.Title = parser.Clean(summary.Title)
	summary.Description = parser.Clean(summary.Description)
	summary.Extract = parser.Clean(summary.Extract)
	return summary, nil
}
