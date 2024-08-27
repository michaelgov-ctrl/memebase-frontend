package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(app.config.staticDir))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /meme/view/{id}", dynamic.ThenFunc(app.memeView))
	mux.Handle("GET /meme/create", dynamic.ThenFunc(app.memeCreate))
	mux.Handle("POST /meme/create", dynamic.ThenFunc(app.memeCreatePost))
	mux.HandleFunc("/teapot", app.teapot)
	mux.HandleFunc("/coffee", app.teapot)

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
