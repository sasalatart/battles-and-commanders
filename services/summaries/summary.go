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

// Get queries a Wikipedia URL page for its corresponding summary API endpoint, and returns a
// domain.Summary if successful, or an error if not
func Get(url string) (domain.Summary, error) {
	summary := domain.Summary{}

	if !strings.Contains(url, "/wiki/") {
		return summary, errors.Errorf("URL is not a wiki page: %s", url)
	}

	summaryURL := strings.ReplaceAll(url, "/wiki/", "/api/rest_v1/page/summary/")
	resp, err := http.Get(summaryURL)
	if err != nil {
		return summary, errors.Wrapf(err, "Fetching %s", summaryURL)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return summary, errors.Wrapf(err, "Reading response for %s", summaryURL)
	}
	if err = json.Unmarshal(body, &summary); err != nil {
		return summary, errors.Wrapf(err, "Unmarshalling response for %s", summaryURL)
	}
	if strings.Contains(strings.ToLower(summary.Title), "not found") || strings.Contains(summary.Type, "not_found") {
		return summary, errors.Errorf("No summary found in %s", summaryURL)
	}
	for _, p := range []*string{&summary.DisplayTitle, &summary.Description, &summary.Extract} {
		*p = parser.Clean(*p)
	}
	return summary, nil
}
