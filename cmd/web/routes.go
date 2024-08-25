package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(app.config.staticDir))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /meme/view/{id}", app.memeView)
	mux.HandleFunc("GET /meme/create", app.memeCreate)
	mux.HandleFunc("POST /meme/create", app.memeCreatePost)
	mux.HandleFunc("/teapot", app.teapot)
	mux.HandleFunc("/coffee", app.teapot)

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
