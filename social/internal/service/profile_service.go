package service

import "example.com/social/internal/domain"

type ProfileService struct {
	ProfileStorage ProfileStorage
}

func NewProfileService(storage ProfileStorage) *ProfileService {
	return &ProfileService{
		storage,
	}
}

func (profileService *ProfileService) GetProfileByUsername(username string) (*domain.Profile, error) {
	return profileService.ProfileStorage.GetProfileByUsername(username)
}