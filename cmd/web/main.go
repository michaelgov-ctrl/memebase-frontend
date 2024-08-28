package main

import (
	"crypto/tls"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mongodbstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
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
	formDecoder    *form.Decoder
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

	mongoClient, err := openMongoConnection(dba)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer closeMongoConnection(mongoClient)
	logger.Info("database connection pool established")

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	sessionManager := scs.New()
	sessionManager.Store = mongodbstore.New(mongoClient.Database("memebase_session_manager"))
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app := &application{
		config:         cfg,
		logger:         logger,
		models:         models.NewModels(mongoClient),
		templateCache:  templateCache,
		sessionManager: sessionManager,
		formDecoder:    form.NewDecoder(),
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         cfg.addr,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.logger.Info("starting server", "addr", app.config.addr)

	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	app.logger.Error(err.Error())
	os.Exit(1)
}
