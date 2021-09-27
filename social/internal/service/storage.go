package service

import "example.com/social/internal/domain"

type ProfileStorage interface {
	SaveProfile(profile *domain.Profile) error
	GetProfileByUsername(username string) (*domain.Profile, error)
	GetProfilesBySearchTerm(term string, offset int64, limit int) ([]*domain.Profile, error)
}
