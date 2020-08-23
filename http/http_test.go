package http_test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gofiber/fiber"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type response struct {
	status       int
	errorMessage string
	body         interface{}
}

func mustGet(t *testing.T, app *fiber.App, route string) *http.Response {
	t.Helper()
	req, err := http.NewRequest("GET", route, nil)
	require.NoError(t, err, route)
	res, err := app.Test(req, -1)
	require.NoError(t, err, route)
	return res
}

func assertErrorMessage(t *testing.T, res *http.Response, expectedMessage string) {
	t.Helper()
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err, "Reading from body")
	assert.Equal(t, expectedMessage, string(body), "Comparing body with expected error message")
}
