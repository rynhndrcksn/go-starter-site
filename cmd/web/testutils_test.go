package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
)

// newTestApplication creates a new application struct containing mocked dependencies.
func newTestApplication(t *testing.T) *application {
	// Initialize a new template cache.
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	// Initialize a new session manager instance.
	// Note: these settings should match the ones used in production.
	// The only exception is to use an in-memory store.
	sessionManager := scs.New()
	sessionManager.Cookie.Secure = true
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.HashTokenInStore = true
	sessionManager.Store = memstore.New()

	return &application{
		logger:         slog.New(slog.NewTextHandler(io.Discard, nil)),
		templateCache:  templateCache,
		sessionManager: sessionManager,
	}
}

// testServer embeds a *httptest.Server instance.
type testServer struct {
	*httptest.Server
}

// newTestServer initializes and returns a new testServer instance.
func newTestServer(t *testing.T, h http.Handler) *testServer {
	// Initialize a new test server.
	ts := httptest.NewTLSServer(h)

	// Initialize a new cookie jar.
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add the cookie jar to the test server client.
	// Any response cookies will now be stored and sent with other requests when using this client.
	ts.Client().Jar = jar

	// Disable redirect-following for the test server client by setting a custom CheckRedirect function.
	// This function will be called whenever a 3xx response is received by the client, and by always returning a
	// http.ErrUseLastResponse error it forces the client to immediately return the received response.
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

// get will make a GET request to a given URL path using the test server client
// and returns the response status code, response headers, and response body.
func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(rs.Body)

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}
