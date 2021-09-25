package service

import (
	"example.com/social/internal/domain"
)

type AuthService struct {
	Storage ProfileStorage
}

func NewAuthService(storage ProfileStorage) *AuthService {
	return &AuthService{
		storage,
	}
}

func (authService *AuthService) SignUp(profile *domain.Profile) error {
	err := authService.Storage.SaveProfile(profile)
	if err != nil {
		return err
	}

	return nil
}

func (authService *AuthService) SignIn(credentials *domain.Credentials) (bool, error) {
	profile, err := authService.Storage.GetProfileByUsername(credentials.Username)
	if err != nil {
		return false, err
	}

	if profile == nil {
		return false, domain.ProfileNotFound
	}

	isMatch, err := credentials.CheckPassword(profile)
	if err != nil {
		return false, err
	}

	return isMatch, nil
}