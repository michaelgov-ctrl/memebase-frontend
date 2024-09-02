package models

import (
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type Models struct {
	Memes *MemeModel
	Users *UserModel
}

func NewModels(mongoClient *mongo.Client) Models {
	return Models{
		Memes: &MemeModel{
			Client:        &http.Client{Timeout: 5 * time.Second},
			Host:          "http://localhost:4000",
			MemesEndpoint: "/v1/memes",
		},
		Users: &UserModel{DB: mongoClient.Database("memebase")},
	}
}
