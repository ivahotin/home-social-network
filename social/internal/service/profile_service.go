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

func (profileService *ProfileService) GetProfilesBySearchTerm(firstname, lastname string, cursor int64, limit int, myId int64) (*domain.ProfilesSearchResult, error) {
	profiles, err := profileService.ProfileStorage.GetProfilesBySearchTerm(firstname, lastname, cursor, limit)
	if err != nil {
		return nil, err
	}

	profilesSearchResult := domain.NewProfilesSearchResult(profiles, cursor, myId)
	return profilesSearchResult, nil
}

func (profileService *ProfileService) GetProfilesByUserIds(userIds []int64) ([]*domain.Profile, error) {
	return profileService.ProfileStorage.GetProfilesByIds(userIds)
}

func (profileService *ProfileService) GetProfileByUserId(userId int64) (*domain.Profile, error) {
	return profileService.ProfileStorage.GetProfileByUserId(userId)
}