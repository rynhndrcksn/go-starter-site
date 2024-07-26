package main

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

// logError is a generic helper for logging error messages.
func (app *application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		trace  = string(debug.Stack())
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), slog.String("method", method), slog.String("uri", uri), slog.String("trace", trace))
}
