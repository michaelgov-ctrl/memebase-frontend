package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Meme struct {
	ID      string     `json:"id,omitempty"` //"go.mongodb.org/mongo-driver/bson/primitive" ID     primitive.ObjectID
	Created *time.Time `json:"created,omitempty"`
	Artist  string     `json:"artist"`
	Title   string     `json:"title"`
	B64     string     `json:"b64"`
	Version int32      `json:"version,omitempty"`
}

type MemeModel struct {
	Client       *http.Client
	Host         string
	GetEndpoint  string
	PostEndpoint string
}

func (m *MemeModel) PostMeme(meme *Meme) (string, error) {
	url, err := url.JoinPath(m.Host, m.PostEndpoint)
	if err != nil {
		return "", err
	}

	contentType := "application/json"

	b, err := json.Marshal(meme)
	if err != nil {
		return "", err
	}

	bodyReader := bytes.NewReader(b)

	resp, err := m.Client.Post(url, contentType, bodyReader)
	if err != nil {
		return "", err
	}

	return resp.Header.Get("Location"), nil
}

func (m *MemeModel) GetById(id string) (*Meme, error) {
	id = "66c78b9d663ce3e2a6bf35c4"
	url, err := url.JoinPath(m.Host, m.GetEndpoint, id)
	if err != nil {
		return nil, err
	}

	resp, err := m.Client.Get(url)
	if err != nil {
		return nil, err
	}

	fmt.Println(resp)
	return &Meme{}, nil
}
