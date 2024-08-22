package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rynhndrcksn/go-starter-site/ui"
)

// routes handles assigning all the routes for the site and what HTTP methods are used for them.
func (app *application) routes() http.Handler {
	// Initialize new httprouter instance.
	router := httprouter.New()

	// Use the http.FileServerFS() function to create an HTTP handler which serves the embedded files in ui.Files.
	// It's important to note that the static files are contained in the "static" folder of the ui.Files embedded filesystem.
	// So, for example, our CSS stylesheet is located at "static/css/main.css".
	router.Handler(http.MethodGet, "/static/*path", http.FileServerFS(ui.Files))

	// Tell httprouter to use our custom notFoundResponse handler.
	router.NotFound = http.HandlerFunc(app.notFoundHandler)

	// Register routes.
	router.HandlerFunc(http.MethodGet, "/", app.homeHandler)
	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	return app.recoverPanic(app.logRequest(app.commonHeaders(router)))
}
