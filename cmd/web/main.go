package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rynhndrcksn/go-starter-site/internal/data"
	"github.com/rynhndrcksn/go-starter-site/internal/vcs"
)

// config contains all the project configuration.
type config struct {
	port int
	env  string
	dsn  string
}

// application contains the stuff used across the project.
type application struct {
	config         config
	debug          bool
	logger         *slog.Logger
	wg             sync.WaitGroup
	templateCache  map[string]*template.Template
	sessionManager *scs.SessionManager
	models         data.Models
}

func main() {
	var conf config
	flag.IntVar(&conf.port, "port", 4000, "Web server port")
	flag.StringVar(&conf.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&conf.dsn, "dsn", os.Getenv("DB_CONN"), "Database DSN")
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

	// Initialize a new DB connection
	db, err := openDB(conf)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("database connection pool established")

	// Initialize a new session manager with an in memory session store.
	// Documentation can be found here: https://pkg.go.dev/github.com/alexedwards/scs/v2
	sessionManager := scs.New()
	sessionManager.Cookie.Secure = true
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.HashTokenInStore = true
	sessionManager.Store = pgxstore.New(db)

	// Initialize a new application struct.
	app := &application{
		config:         conf,
		debug:          *debug,
		logger:         logger,
		templateCache:  templateCache,
		models:         data.NewModels(db),
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

// openDB returns a pgxpool.Pool.
func openDB(cfg config) (*pgxpool.Pool, error) {
	// Use pgxpool.New() to create an empty connection pool, using the DSN from the config struct.
	// Note: pgxpool parses configurations from the connection string.
	// i.e. postgres://.../myDatabase?sslmode=verify-ca&pool_max_conns=10
	// More information can be found here: https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool@v5.6.0#Config
	db, err := pgxpool.New(context.Background(), cfg.dsn)
	if err != nil {
		return nil, err
	}

	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use Ping() to establish a new connection to the database.
	// If the connection couldn't be established successfully within the 5-second deadline, then this will return an error.
	err = db.Ping(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	// Return the pgxpool.Pool connection pool.
	return db, nil
}
