package main

import (
	"net/http"
	"runtime/debug"
)

// homeHandler displays the home page.
func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "home.tmpl", data)
}

// notFoundHandler displays a 404 page.
func (app *application) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusNotFound, "404.tmpl", data)
}

// serverErrorHandler displays an error page due to a server error.
func (app *application) serverErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	if app.debug {
		http.Error(w, string(debug.Stack()), http.StatusInternalServerError)
		return
	}
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusInternalServerError, "500.tmpl", data)
}
