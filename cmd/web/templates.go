package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/rynhndrcksn/go-starter-site/internal/env"
	"github.com/rynhndrcksn/go-starter-site/ui"
)

var (
	errHashAssetPathIsEmpty       = errors.New("cannot pass an empty path in")
	errHashAssetCantReadFile      = errors.New("error trying to read original file with fs.ReadFile()")
	errPropsKeyValueCountMismatch = errors.New("mismatched amount of key/value pairs")
	errPropsKeyValueCountIsZero   = errors.New("length of 'pairs' must be greater than 0")
)

// functions contains a template.FuncMap that maps the above functions to functions that can then be called inside the templates.
var functions = template.FuncMap{
	"hashAssetPath": hashAssetPath,
	"humanDate":     humanDate,
	"props":         props,
}

// hashAssetPath takes an asset path, computes the has for the asset, appends it to the asset name, and returns it.
// This acts as a way to enable caching the assets long term, but if they change, the cache can be bypassed.
// The path passed in must start with a "/", and is expected to be used in a similar setup to this project.
// This works by appending a ?v=<hash> to the end of the file, so it isn't incredibly robust for advanced needs.
// This is heavily inspired by https://github.com/c9845/hashfs
func hashAssetPath(originalPath string) (string, error) {
	if len(strings.TrimSpace(originalPath)) == 0 {
		return "", errHashAssetPathIsEmpty
	}

	// Strip off the first "/", so the following fs.ReadFile line succeeds.
	contents, err := fs.ReadFile(ui.Files, originalPath[1:])
	if err != nil {
		return originalPath, errHashAssetCantReadFile
	}

	// Get the file hash and encode it to a string.
	hash := md5.Sum(contents)
	encodedHash := hex.EncodeToString(hash[:])

	// Add hash to filename.
	hashedFilename := fmt.Sprintf("%s?v=%s", originalPath, encodedHash)

	return hashedFilename, nil
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
	if len(pairs) == 0 {
		return nil, errPropsKeyValueCountIsZero
	}

	if len(pairs)%2 != 0 {
		return nil, errPropsKeyValueCountMismatch
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
	CanonicalUrl string
	CurrentYear  int
	Description  string
	Flash        string
	ImageUrl     string
	PageType     string
	SiteName     string
	Title        string
}

// newTemplateData initializes a new templateData struct and returns it.
func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CanonicalUrl: getCanonicalURL(r),
		CurrentYear:  time.Now().Year(),
		Flash:        app.sessionManager.PopString(r.Context(), "flash"),
		SiteName:     env.GetStringOrDefault("SITE_NAME", "Site"),
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

// getCanonicalURL creates the canonical URL and returns it.
func getCanonicalURL(r *http.Request) string {
	// Get the full URL from the request
	scheme := "https"

	// Remove www if present and clean the host
	host := r.Host
	host = strings.TrimPrefix(host, "www.")

	// Build the canonical URL without query parameters
	canonicalURL := scheme + "://" + host + r.URL.Path

	return canonicalURL
}
