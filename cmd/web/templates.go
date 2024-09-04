package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/michaelgov-ctrl/memebase-front/internal/models"
	"github.com/michaelgov-ctrl/memebase-front/ui"
)

type templateData struct {
	CurrentYear     int
	Flash           string
	Form            any
	Meme            models.Meme
	Memes           []models.Meme
	Metadata        models.Metadata
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	var cache = make(map[string]*template.Template)

	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl.html",
			"html/partials/*.html",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
