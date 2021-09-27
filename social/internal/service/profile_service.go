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

func (profileService *ProfileService) GetProfilesBySearchTerm(term string, cursor int64, limit int, myId int64) (*domain.ProfilesSearchResult, error) {
	profiles, err := profileService.ProfileStorage.GetProfilesBySearchTerm(term, cursor, limit)
	if err != nil {
		return nil, err
	}

	profilesSearchResult := domain.NewProfilesSearchResult(profiles, cursor, myId)
	return profilesSearchResult, nil
}