package domain

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Username 	string
	Password 	string
	RawPassword string
}

func NewCredentials(username string, password string) (*Credentials, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, err
	}

	return &Credentials{
		Username: username,
		Password: string(hash),
		RawPassword: password,
	}, nil
}

func (credentials *Credentials) CheckPassword(profile *Profile) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(profile.Password), []byte(credentials.RawPassword))
	switch {
	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
		return false, nil
	case err != nil:
		return false, err
	}

	return true, nil
}