package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/michaelgov-ctrl/memebase-front/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	mlr, err := app.models.Memes.GetMemeList("")
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	for _, meme := range mlr.Memes {
		fmt.Fprintf(w, "%+v\n", meme)
	}

	fmt.Fprintf(w, "%+v\n", mlr.Metadata)

	/*
		//base must be the first file in the slice
		files := []string{
			"./ui/html/base.tmpl.html",
			"./ui/html/partials/nav.tmpl.html",
			"./ui/html/pages/home.tmpl.html",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		err = ts.ExecuteTemplate(w, "base", nil)
		if err != nil {
			app.serverError(w, r, err)
		}
	*/
}

func (app *application) memeView(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if len(id) != 24 { // len of mongo returned id
		http.NotFound(w, r)
		return
	}

	// id = "66c78b9d663ce3e2a6bf35c4"
	meme, err := app.models.Memes.GetById(id)
	if err != nil {
		if errors.Is(err, models.ErrNoMeme) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	fmt.Fprintf(w, "id: %s\nmeme: %v", id, meme)
}

func (app *application) memeCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new meme..."))
}

func (app *application) memeCreatePost(w http.ResponseWriter, r *http.Request) {
	meme := models.Meme{
		Artist: "dddd",
		Title:  "qqqq",
		B64:    "yeep",
	}

	location, err := app.models.Memes.PostMeme(&meme)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	parts := strings.Split(location, "/")
	id := parts[len(parts)-1]

	http.Redirect(w, r, fmt.Sprintf("/meme/view/%s", id), http.StatusSeeOther)
}

func (app *application) teapot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTeapot)

	w.Write([]byte("I'm a teapot..."))
}
