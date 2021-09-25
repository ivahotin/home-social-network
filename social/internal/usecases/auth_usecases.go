package usecases

import "example.com/social/internal/domain"

type SignUpUseCase interface {
	SignUp(profile *domain.Profile) error
}

type SignInUseCase interface {
	SignIn(credentials *domain.Credentials) (bool, error)
}