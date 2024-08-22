package main

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/michaelgov-ctrl/memebase-front/internal/data"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

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
}

func (app *application) memeView(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if len(id) != 24 { // len of mongo returned id
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display a specific meme with ID %s...", id)
}

func (app *application) memeCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new meme..."))
}

func (app *application) memeCreatePost(w http.ResponseWriter, r *http.Request) {
	meme := data.Meme{
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
