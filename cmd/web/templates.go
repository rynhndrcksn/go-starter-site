package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"time"

	"github.com/rynhndrcksn/portfolio/ui"
)

// functions contains a template.FuncMap that maps the above functions to functions that can then be called inside the templates.
var functions = template.FuncMap{
	"humanDate": humanDate,
	"props":     props,
}

// humanDate returns a nicely formatted string representation of a time.Time object.
func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

// props takes any number of key/value pairs and passes them into a child template.
func props(pairs ...any) (map[string]any, error) {
	if len(pairs)%2 != 0 {
		return nil, errors.New("mismatched amount of key/value pairs")
	}

	m := make(map[string]any, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		key, ok := pairs[i].(string)

		if !ok {
			return nil, fmt.Errorf("cannot use type %T as map key", pairs[i])
		}
		m[key] = pairs[i+1]
	}
	return m, nil
}

// templateData holds dynamic data that can be passed to the HTML templates.
type templateData struct {
	CurrentYear int
	Flash       string
}

// newTemplateData initializes a new templateData struct and returns it.
func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear: time.Now().Year(),
		Flash:       app.sessionManager.PopString(r.Context(), "flash"),
	}
}

// newTemplateCache grabs all the templates in ui/html/, renders them, and adds them to a map.
// This way the template doesn't have to be rendered on every request.
func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	// Use fs.Glob() to get a slice of all file paths in the ui.Files embedded filesystem which match the pattern 'html/pages/*.tmpl'.
	// This gives us a slice of all the 'page' templates for the application.
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	// Loop through the page file paths one-by-one.
	for _, page := range pages {
		// Extract the file name (like 'home.tmpl') from the full filepath and assign it to the new variable.
		name := filepath.Base(page)

		// Create a slice containing the filepath patterns for the templates we want to parse.
		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			"html/components/*.tmpl",
			page,
		}

		// Use ParseFS() to parse the template files from the ui.Files embedded filesystem.
		var ts *template.Template
		ts, err = template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map.
		cache[name] = ts
	}

	return cache, nil
}
