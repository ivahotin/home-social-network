package usecases

import "example.com/social/internal/domain"

type GetProfileGetUsernameUseCase interface {
	GetProfileByUsername(username string) (*domain.Profile, error)
}

type GetProfilesBySearchTerm interface {
	GetProfilesBySearchTerm(term string, cursor int64, limit int, myId int64) (*domain.ProfilesSearchResult, error)
}
