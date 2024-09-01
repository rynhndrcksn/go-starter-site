package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

// commonHeaders sets all the default headers we want on every request.
func (app *application) commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Note: This is split across multiple lines for readability.
		w.Header().Set("Content-Security-Policy",
			"default-src 'self';"+
				"frame-ancestors 'none';")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		next.ServeHTTP(w, r)
	})
}

// logRequests will log information for each request the site gets.
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.RequestURI
		)
		if !strings.Contains(uri, "static") {
			app.logger.Info("Received request", slog.String("ip", ip), slog.String("proto", proto), slog.String("method", method), slog.String("uri", uri))
		}
		next.ServeHTTP(w, r)
	})
}

// recoverPanic checks if there's been a panic and returns a 500 error rather than no response.
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event of a panic)
		defer func() {
			// Use the builtin recover function to check if there has been a panic or not.
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorHandler(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
