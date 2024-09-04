package main

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/michaelgov-ctrl/memebase-front/ui"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)
	protected := dynamic.Append(app.requireAuthentication)

	// Create and read meme endpoints
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /meme/view/{id}", dynamic.ThenFunc(app.memeView))
	mux.Handle("GET /meme/create", protected.ThenFunc(app.memeCreate))
	mux.Handle("POST /meme/create", protected.ThenFunc(app.memeCreatePost))

	// Authorization endpoints
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	// Joke endpoints
	mux.HandleFunc("/teapot", app.teapot)
	mux.HandleFunc("/coffee", app.teapot)

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
