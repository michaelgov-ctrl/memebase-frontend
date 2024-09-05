package models

import (
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type MemeModelInterface interface {
	PostMeme(meme *Meme) (string, error)
	GetById(id string) (*Meme, error)
	GetMemeList(queryString string) (*MemeListResponse, error)
}

type UserModelInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (string, error)
	Exists(id string) (bool, error)
}

type Models struct {
	Memes MemeModelInterface
	Users UserModelInterface
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
