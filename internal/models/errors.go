package models

import "errors"

var (
	ErrNoMeme         = errors.New("models: no matching meme found")
	ErrDocNotFound    = errors.New("document not found")
	ErrEditConflict   = errors.New("edit conflict")
	ErrFailedCreation = errors.New("meme creation failed")
)
