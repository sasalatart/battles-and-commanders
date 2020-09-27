package httptest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/sasalatart/batcoms/domain/battles"
	"github.com/sasalatart/batcoms/domain/commanders"
	"github.com/sasalatart/batcoms/domain/factions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertFailedGET asserts that issuing a GET request to the specified route renders the given
// status and error message
func AssertFailedGET(t *testing.T, route string, expectedStatus int, expectedMessage string) {
	t.Helper()
	res, err := http.Get(route)
	require.NoErrorf(t, err, "Requesting %s", route)
	defer res.Body.Close()
	assert.Equal(t, expectedStatus, res.StatusCode, "Comparing status with expected value")
	AssertErrorMessage(t, res, expectedMessage)
}

// AssertErrorMessage asserts that the given *http.Response contains the specified error message
func AssertErrorMessage(t *testing.T, res *http.Response, expectedMessage string) {
	t.Helper()
	contents, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err, "Reading from response body")
	assert.Equal(t, expectedMessage, string(contents), "Comparing body with expected error message")
}

// AssertJSONFaction asserts that the given *http.Response contains the specified JSON-serialized
// factions.Faction
func AssertJSONFaction(t *testing.T, res *http.Response, expectedFaction factions.Faction) {
	t.Helper()
	factionFromBody := new(factions.Faction)
	err := json.NewDecoder(res.Body).Decode(factionFromBody)
	require.NoError(t, err, "Decoding body into faction struct")
	assert.Equal(t, expectedFaction, *factionFromBody, "Comparing body with expected faction")
}

// AssertJSONFactions is like AssertJSONFaction, but for a slice of factions.Faction
func AssertJSONFactions(t *testing.T, res *http.Response, expectedFactions []factions.Faction) {
	t.Helper()
	factionsFromBody := new([]factions.Faction)
	err := json.NewDecoder(res.Body).Decode(factionsFromBody)
	require.NoError(t, err, "Decoding body into factions slice")
	assert.Equal(t, expectedFactions, *factionsFromBody, "Comparing body with expected factions")
}

// AssertJSONCommander asserts that the given *http.Response contains the specified JSON-serialized
// commanders.Commander
func AssertJSONCommander(t *testing.T, res *http.Response, expectedCommander commanders.Commander) {
	t.Helper()
	commanderFromBody := new(commanders.Commander)
	err := json.NewDecoder(res.Body).Decode(commanderFromBody)
	require.NoError(t, err, "Decoding body into commander struct")
	assert.Equal(t, expectedCommander, *commanderFromBody, "Comparing body with expected commander")
}

// AssertJSONCommanders is like AssertJSONCommander, but for a slice of commanders.Commander
func AssertJSONCommanders(t *testing.T, res *http.Response, expectedCommanders []commanders.Commander) {
	t.Helper()
	commandersFromBody := new([]commanders.Commander)
	err := json.NewDecoder(res.Body).Decode(commandersFromBody)
	require.NoError(t, err, "Decoding body into commanders slice")
	assert.Equal(t, expectedCommanders, *commandersFromBody, "Comparing body with expected commanders")
}

// AssertJSONBattle asserts that the given *http.Response contains the specified JSON-serialized
// battles.Battle
func AssertJSONBattle(t *testing.T, res *http.Response, expectedBattle battles.Battle) {
	t.Helper()
	battleFromBody := new(battles.Battle)
	err := json.NewDecoder(res.Body).Decode(battleFromBody)
	require.NoError(t, err, "Decoding body into battle struct")
	assert.Equal(t, expectedBattle, *battleFromBody, "Comparing body with expected battle")
}

// AssertJSONBattles is like AssertJSONBattle, but for a slice of battles.Battle
func AssertJSONBattles(t *testing.T, res *http.Response, expectedBattles []battles.Battle) {
	t.Helper()
	battlesFromBody := new([]battles.Battle)
	err := json.NewDecoder(res.Body).Decode(battlesFromBody)
	require.NoError(t, err, "Decoding body into battles slice")
	assert.Equal(t, expectedBattles, *battlesFromBody, "Comparing body with expected battles")
}

// AssertHeaderPages asserts that the given *http.Response has the expected "x-pages" header value
func AssertHeaderPages(t *testing.T, res *http.Response, expectedPages int) {
	t.Helper()
	expected := fmt.Sprint(expectedPages)
	got := res.Header.Get("x-pages")
	assert.Equal(t, expected, got, "Comparing with the expected 'x-pages' header")
}
