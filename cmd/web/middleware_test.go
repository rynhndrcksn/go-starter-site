package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rynhndrcksn/go-starter-site/internal/assert"
)

func TestCommonHeaders(t *testing.T) {
	// Initialize a new httptest.ResponseRecorder and fake http.Request.
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock HTTP handler that we can pass to our commonHeaders middleware, which writes a 200 status code and an "OK" response body.
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Initialize a new test application.
	app := newTestApplication(t)

	// Pass the mock HTTP handler to our commonHeaders middleware.
	app.commonHeaders(next).ServeHTTP(rr, r)

	// Call the Result() method on the http.ResponseRecorder to get the results of the test.
	rs := rr.Result()

	// Check that the middleware has correctly set the Content-Security-Policy header on the response.
	expectedValue := "default-src 'self';frame-ancestors 'none';"
	assert.Equal(t, rs.Header.Get("Content-Security-Policy"), expectedValue)

	// Check that the middleware has correctly set the Referrer-Policy header on the response.
	expectedValue = "origin-when-cross-origin"
	assert.Equal(t, rs.Header.Get("Referrer-Policy"), expectedValue)

	expectedValue = "max-age=63072000; includeSubDomains; preload"
	assert.Equal(t, rs.Header.Get("Strict-Transport-Security"), expectedValue)

	// Check that the middleware has correctly set the X-Content-Type-Options header on the response.
	expectedValue = "nosniff"
	assert.Equal(t, rs.Header.Get("X-Content-Type-Options"), expectedValue)

	// Check that the middleware has correctly set the X-Frame-Options header on the response.
	expectedValue = "deny"
	assert.Equal(t, rs.Header.Get("X-Frame-Options"), expectedValue)

	// Check that the middleware has correctly called the next handler in line and the response status code and body are as expected.
	assert.Equal(t, rs.StatusCode, http.StatusOK)

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

	assert.Equal(t, string(body), "OK")
}
