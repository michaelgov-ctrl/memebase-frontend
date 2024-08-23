package models

import (
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type Models struct {
	Auth  AuthModel
	Memes MemeModel
}

func NewModels(db *mongo.Client) Models {
	return Models{
		Auth: AuthModel{DB: db},
		Memes: MemeModel{
			Client:        &http.Client{Timeout: 5 * time.Second},
			Host:          "http://localhost:4000",
			MemesEndpoint: "/v1/memes",
		},
	}
}
