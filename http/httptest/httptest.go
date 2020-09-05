package httptest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/sasalatart/batcoms/domain"
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
// domain.Faction
func AssertJSONFaction(t *testing.T, res *http.Response, expectedFaction domain.Faction) {
	t.Helper()
	factionFromBody := new(domain.Faction)
	err := json.NewDecoder(res.Body).Decode(factionFromBody)
	require.NoError(t, err, "Decoding body into faction struct")
	assert.Equal(t, expectedFaction, *factionFromBody, "Comparing body with expected faction")
}

// AssertJSONCommander asserts that the given *http.Response contains the specified JSON-serialized
// domain.Commander
func AssertJSONCommander(t *testing.T, res *http.Response, expectedCommander domain.Commander) {
	t.Helper()
	commanderFromBody := new(domain.Commander)
	err := json.NewDecoder(res.Body).Decode(commanderFromBody)
	require.NoError(t, err, "Decoding body into commander struct")
	assert.Equal(t, expectedCommander, *commanderFromBody, "Comparing body with expected commander")
}

// AssertJSONBattle asserts that the given *http.Response contains the specified JSON-serialized
// domain.Battle
func AssertJSONBattle(t *testing.T, res *http.Response, expectedBattle domain.Battle) {
	t.Helper()
	battleFromBody := new(domain.Battle)
	err := json.NewDecoder(res.Body).Decode(battleFromBody)
	require.NoError(t, err, "Decoding body into battle struct")
	assert.Equal(t, expectedBattle, *battleFromBody, "Comparing body with expected battle")
}
