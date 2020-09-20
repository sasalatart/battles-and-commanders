package summaries

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/pkg/strclean"
)

// Fetch queries a Wikipedia URL page for its corresponding summary API endpoint, and returns a
// Summary if successful, or an error if not
func Fetch(url string) (Summary, error) {
	summary := Summary{}

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
	summary.Title = strclean.Apply(summary.Title)
	summary.Description = strclean.Apply(summary.Description)
	summary.Extract = strclean.Apply(summary.Extract)
	return summary, nil
}
