package main

import (
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mongodbstore"
	"github.com/alexedwards/scs/v2"
	"github.com/michaelgov-ctrl/memebase-front/internal/models"
)

type config struct {
	addr      string
	staticDir string
}

type application struct {
	config         config
	logger         *slog.Logger
	models         models.Models
	templateCache  map[string]*template.Template
	sessionManager *scs.SessionManager
}

func main() {
	var cfg config

	flag.StringVar(&cfg.addr, "addr", ":4040", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")

	var dba dbAuth
	flag.StringVar(&dba.user, "db-user", "root", "MongoDB user")
	flag.StringVar(&dba.password, "db-password", "example", "MongoDB password")
	flag.StringVar(&dba.uri, "db-uri", "mongodb://localhost:27017", "MongoDB URI")

	flag.Parse()

	logOpts := &slog.HandlerOptions{AddSource: true}
	logger := slog.New(slog.NewTextHandler(os.Stdout, logOpts))

	db, err := openDB(dba)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer closeDB(db)
	logger.Info("database connection pool established")

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	sessionManager := scs.New()
	sessionManager.Store = mongodbstore.New(db.Database("memebase"))
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		config:         cfg,
		logger:         logger,
		models:         models.NewModels(db),
		templateCache:  templateCache,
		sessionManager: sessionManager,
	}

	app.logger.Info("starting server", "addr", app.config.addr)
	err = http.ListenAndServe(app.config.addr, app.routes())
	app.logger.Error(err.Error())
	os.Exit(1)
}
