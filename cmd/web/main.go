package main

import (
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/rynhndrcksn/portfolio/internal/vcs"
)

// config contains all the project configuration.
type config struct {
	port int
	env  string
}

// application contains the stuff used across the project.
type application struct {
	config         config
	debug          bool
	logger         *slog.Logger
	wg             sync.WaitGroup
	templateCache  map[string]*template.Template
	sessionManager *scs.SessionManager
}

func main() {
	var conf config
	flag.IntVar(&conf.port, "port", 4000, "Web server port")
	flag.StringVar(&conf.env, "env", "development", "Environment (development|staging|production)")
	debug := flag.Bool("debug", false, "Enable debug mode")
	displayVersion := flag.Bool("version", false, "Display version and exit")
	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version: %s\n", vcs.Version())
		os.Exit(0)
	}

	// Initialize new structured logger that writes to stdout.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Initialize a new template cache.
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Initialize a new session manager with an in memory session store.
	// Documentation can be found here: https://pkg.go.dev/github.com/alexedwards/scs/v2
	sessionManager := scs.New()
	sessionManager.Cookie.Secure = true
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.HashTokenInStore = true
	sessionManager.Store = memstore.New()

	// Initialize a new application struct.
	app := &application{
		config:         conf,
		debug:          *debug,
		logger:         logger,
		templateCache:  templateCache,
		sessionManager: sessionManager,
	}

	// Launch the site.
	err = app.serve()
	if err != nil {
		app.logger.Error("Error shutting down server", slog.String("error", err.Error()))
		os.Exit(1)
	}
	os.Exit(0)
}
