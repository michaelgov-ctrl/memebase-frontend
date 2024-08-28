package models

import "errors"

var (
	ErrNoMeme             = errors.New("models: no matching meme found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
	ErrDocNotFound        = errors.New("document not found")
	ErrEditConflict       = errors.New("edit conflict")
	ErrFailedCreation     = errors.New("meme creation failed")
)
