package main

import (
	"path/filepath"
	"text/template"

	"github.com/michaelgov-ctrl/memebase-front/internal/models"
)

type templateData struct {
	Meme     models.Meme
	Memes    []models.Meme
	Metadata models.Metadata
}

func newTemplateCache() (map[string]*template.Template, error) {
	var cache = make(map[string]*template.Template)

	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
