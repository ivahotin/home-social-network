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

func (authService *AuthService) SignIn(credentials *domain.Credentials) (*domain.SignInResult, error) {
	profile, err := authService.Storage.GetProfileByUsername(credentials.Username)
	if err != nil {
		return nil, err
	}

	if profile == nil {
		return nil, domain.ProfileNotFound
	}

	isMatch, err := credentials.CheckPassword(profile)
	if err != nil {
		return nil, err
	}

	return &domain.SignInResult{IsMatch: isMatch, Id: profile.Id}, nil
}