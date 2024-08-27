package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/michaelgov-ctrl/memebase-front/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	mlr, err := app.models.Memes.GetMemeList("")
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Memes = mlr.Memes
	data.Metadata = mlr.Metadata

	app.render(w, r, http.StatusOK, "home.tmpl.html", data)
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

	data := app.newTemplateData(r)
	data.Meme = *meme

	app.render(w, r, http.StatusOK, "view.tmpl.html", data)
}

func (app *application) memeCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	app.render(w, r, http.StatusOK, "create.tmpl.html", data)
}

func (app *application) memeCreatePost(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576)

	err := r.ParseMultipartForm(1_280)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image_file")
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	defer file.Close()

	contentType := handler.Header["Content-Type"]
	if len(contentType) != 1 || !strings.Contains(contentType[0], "image") {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// TODO: test the file sizes
	fmt.Printf("File Size: %+v\n", handler.Size)

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	meme := models.Meme{
		Title:  r.PostForm.Get("title"),
		Artist: r.PostForm.Get("artist"),
		Image: models.Image{
			Type:  contentType[0],
			Bytes: fileBytes,
		},
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
