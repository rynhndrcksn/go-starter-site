package main

import (
	"net/http"
	"testing"

	"github.com/rynhndrcksn/go-starter-site/internal/assert"
)

func TestHomeHandler(t *testing.T) {
	// Make a new test application.
	app := newTestApplication(t)

	// Make a new test server to make calls with.
	ts := newTestServer(t, app.sessionManager.LoadAndSave(app.routes()))
	defer ts.Close()

	// Make a request to the handler being tested.
	code, _, body := ts.get(t, "/")

	// Assert we're getting a 200 code response.
	assert.Equal(t, code, http.StatusOK)

	// Assert that the body contains the text from the <title> tag.
	assert.StringContains(t, body, "<title>Home - Site</title>")
}

func TestNotFoundHandler(t *testing.T) {
	// Make a new test application.
	app := newTestApplication(t)

	// Make a new test server to make calls with.
	ts := newTestServer(t, app.sessionManager.LoadAndSave(app.routes()))
	defer ts.Close()

	// Make a request to the handler being tested.
	code, _, body := ts.get(t, "/not-found")

	// Assert we're getting a 200 code response.
	assert.Equal(t, code, http.StatusNotFound)

	// Assert that the body contains the text from the <title> tag.
	assert.StringContains(t, body, "<title>Not Found - Site</title>")
}

func TestServerErrorHandler(t *testing.T) {
	// Make a new test application.
	app := newTestApplication(t)

	// Make a new test server to make calls with.
	ts := newTestServer(t, app.sessionManager.LoadAndSave(app.routeThatPanics()))
	defer ts.Close()

	// Make a request to the handler being tested.
	code, _, body := ts.get(t, "/server-error")

	// Assert we're getting a 200 code response.
	assert.Equal(t, code, http.StatusInternalServerError)

	// Assert that the body contains the text from the <title> tag.
	assert.StringContains(t, body, "<title>Server Error - Site</title>")
}
