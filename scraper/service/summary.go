package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sasalatart/batcoms/scraper/domain"
)

// PageSummary queries a Wikipedia URL page for its corresponding summary API endpoint, and returns
// a domain.Summary if successful, or an error if not
func PageSummary(url string) (domain.Summary, error) {
	summary := domain.Summary{}
	summaryURL := strings.ReplaceAll(url, "/wiki/", "/api/rest_v1/page/summary/")

	resp, err := http.Get(summaryURL)
	if err != nil {
		return summary, fmt.Errorf("Failed fetching %s: %s", summaryURL, err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return summary, fmt.Errorf("Failed reading response for %s: %s", summaryURL, err.Error())
	}
	if err = json.Unmarshal(body, &summary); err != nil {
		return summary, fmt.Errorf("Failed unmarshaling response for %s: %s", summaryURL, err.Error())
	}
	return summary, nil
}
