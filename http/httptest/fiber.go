package httptest

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

// AssertFiberGET asserts that the specified route, when handled by the given *fiber.App, renders
// the specified status and satisfies the given assertResponse function
func AssertFiberGET(t *testing.T, app *fiber.App, route string, status int, assertResponse func(*http.Response)) {
	t.Helper()
	req, err := http.NewRequest("GET", route, nil)
	require.NoError(t, err, route)
	res, err := app.Test(req, -1)
	require.NoError(t, err, route)
	defer res.Body.Close()
	require.Equalf(t, status, res.StatusCode, "HTTP status for %q", route)
	assertResponse(res)
}

// AssertFailedFiberGET asserts that the specified route, when handled by the given *fiber.App, renders
// the specified fiber.Error
func AssertFailedFiberGET(t *testing.T, app *fiber.App, route string, status int, message string) {
	t.Helper()
	AssertFiberGET(t, app, route, status, func(res *http.Response) {
		AssertErrorMessage(t, res, message)
	})
}
