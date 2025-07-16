package main

import (
	"expvar"
	"net/http"

	"github.com/rynhndrcksn/go-starter-site/ui"
)

// routes handles assigning all the routes for the site and what HTTP methods are used for them.
func (app *application) routes() http.Handler {
	// Initialize a new http.ServeMux instance.
	mux := http.NewServeMux()

	// Use the http.FileServerFS() function to create an HTTP handler which serves the embedded files in ui.Files.
	// It's important to note that the static files are contained in the "static" folder of the ui.Files embedded filesystem.
	// So, for example, our CSS stylesheet is located at "static/css/main.css".
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	// Register routes.
	mux.HandleFunc("GET /", app.notFoundHandler)
	mux.HandleFunc("GET /{$}", app.homeHandler)
	mux.HandleFunc("GET /about", app.aboutHandler)
	mux.Handle("GET /debug/vars", expvar.Handler())

	return app.recoverPanic(app.logRequest(app.commonHeaders(mux)))
}
