package service

import "example.com/social/internal/domain"

type ProfileStorage interface {
	SaveProfile(profile *domain.Profile) error
	GetProfileByUsername(username string) (*domain.Profile, error)
}
