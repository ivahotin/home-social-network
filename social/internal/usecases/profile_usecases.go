package usecases

import "example.com/social/internal/domain"

type GetProfileGetUsernameUseCase interface {
	GetProfileByUsername(username string) (*domain.Profile, error)
}
