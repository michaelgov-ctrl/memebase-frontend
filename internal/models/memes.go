package models

import (
	"bytes"
	"encoding/json"
	"errors"
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

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

type MemeListResponse struct {
	Memes    []Meme
	Metadata Metadata
}

type MemeModel struct {
	Client        *http.Client
	Host          string
	MemesEndpoint string
}

func (m *MemeModel) PostMeme(meme *Meme) (string, error) {
	url, err := url.JoinPath(m.Host, m.MemesEndpoint)
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
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusSeeOther {
		return "", ErrFailedCreation
	}

	loc := resp.Header.Get("Location")

	return loc, nil
}

func (m *MemeModel) GetById(id string) (*Meme, error) {
	url, err := url.JoinPath(m.Host, m.MemesEndpoint, id)
	if err != nil {
		return nil, err
	}

	resp, err := m.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, ErrNoMeme
		}
		return nil, errors.New(resp.Status)
	}

	var meme = &Meme{}
	err = json.NewDecoder(resp.Body).Decode(meme)
	if err != nil {
		return nil, err
	}

	return meme, nil
}

func (m *MemeModel) GetMemeList(queryString string) (*MemeListResponse, error) {
	url, err := url.JoinPath(m.Host, m.MemesEndpoint, queryString)
	if err != nil {
		return nil, err
	}

	resp, err := m.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrNoMeme
	}

	var mlr = &MemeListResponse{}
	err = json.NewDecoder(resp.Body).Decode(mlr)
	if err != nil {
		return nil, err
	}

	return mlr, err
}
