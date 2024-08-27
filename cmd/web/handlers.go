package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/michaelgov-ctrl/memebase-front/internal/models"
	"github.com/michaelgov-ctrl/memebase-front/internal/validator"
)

type memeCreateForm struct {
	Title       string
	Artist      string
	ContentType string
	validator.Validator
}

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
	data.Form = memeCreateForm{}

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

	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	form := memeCreateForm{
		Title:       r.PostForm.Get("title"),
		Artist:      r.PostForm.Get("artist"),
		ContentType: http.DetectContentType(buf),
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Artist), "artist", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Artist, 100), "artist", "This field cannot be more than 100 characters long")
	form.CheckField(validator.StringMatch(form.ContentType, "image"), "image", "Provided file must be an image")

	// TODO: test the file sizes
	fmt.Printf("File Size: %+v\n", handler.Size)

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	meme := models.Meme{
		Title:  form.Title,
		Artist: form.Artist,
		Image: models.Image{
			Type:  form.ContentType,
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
