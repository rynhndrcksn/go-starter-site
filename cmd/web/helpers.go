package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
)

// render renders the specified template if it exists.
func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	// Retrieve the appropriate template set from the cache based on the page name.
	// If no entry exists in the cache with the provided name, then create a new error and call serverErrorHandler() and return.
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverErrorHandler(w, r, err)
		return
	}

	// Write the template to the buffer, instead of straight to the http.ResponseWriter.
	// If there's an error, call serverErrorHandler() and then return.
	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverErrorHandler(w, r, err)
		return
	}

	// If the template is written to the buffer without any errors, we are safe to go ahead and write the HTTP status code to http.ResponseWriter.
	// Write the contents of the buffer to the http.ResponseWriter.
	w.WriteHeader(status)
	_, err = buf.WriteTo(w)
	if err != nil {
		app.serverErrorHandler(w, r, err)
		return
	}
}

// background accepts an arbitrary function as a parameter.
func (app *application) background(r *http.Request, fn func()) {
	app.wg.Add(1)
	// Launch a background goroutine.
	go func() {
		// Use "defer" to decrement the WaitGroup counter right before the goroutine returns.
		defer app.wg.Done()

		// Recover any panics that arise.
		defer func() {
			if err := recover(); err != nil {
				var (
					method = r.Method
					trace  = string(debug.Stack())
					uri    = r.URL.RequestURI()
				)
				app.logger.Error(fmt.Sprintf("%v", err), slog.String("method", method), slog.String("uri", uri), slog.String("trace", trace))
			}
		}()

		// Execute our arbitrary function that was passed as the parameter.
		fn()
	}()
}
