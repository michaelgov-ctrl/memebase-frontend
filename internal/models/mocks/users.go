package mocks

import "github.com/michaelgov-ctrl/memebase-front/internal/models"

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email, password string) (string, error) {
	if email == "alice@example.com" && password == "pa$$word" {
		return "66c78b9d663ce3e2a6bf35c4", nil
	}

	return "", models.ErrInvalidCredentials
}

func (m *UserModel) Exists(id string) (bool, error) {
	switch id {
	case "66c78b9d663ce3e2a6bf35c4":
		return true, nil
	default:
		return false, nil
	}
}
