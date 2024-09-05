package mocks

import (
	"os"

	"github.com/michaelgov-ctrl/memebase-front/internal/models"
)

var TestImage = models.Image{
	Type:  "jpg",
	Bytes: mockMemeImage("../../internal/models/mocks/test.jpg"),
}

func mockMemeImage(filePath string) []byte {
	body, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	return body
}

var mockMeme = models.Meme{
	ID:     "66c78b9d663ce3e2a6bf35c4",
	Artist: "me",
	Title:  "hilarious meme",
	Image:  TestImage,
}

type MemeModel struct{}

func (m *MemeModel) PostMeme(meme *models.Meme) (string, error) {
	return "69c78b9d663ce3e2a6bf35c4", nil
}

func (m *MemeModel) GetById(id string) (*models.Meme, error) {
	switch id {
	case "66c78b9d663ce3e2a6bf35c4":
		return &mockMeme, nil
	default:
		return &models.Meme{}, models.ErrDocNotFound
	}
}

func (m *MemeModel) GetMemeList(queryString string) (*models.MemeListResponse, error) {
	mlr := &models.MemeListResponse{
		Memes:    []models.Meme{mockMeme},
		Metadata: models.Metadata{},
	}

	return mlr, nil
}
