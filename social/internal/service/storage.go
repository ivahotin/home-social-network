package service

import "example.com/social/internal/domain"

type ProfileStorage interface {
	SaveProfile(profile *domain.Profile) error
	GetProfileByUsername(username string) (*domain.Profile, error)
	GetProfilesBySearchTerm(term string, offset int64, limit int) ([]*domain.Profile, error)
	GetProfilesByIds(userIds []int64) ([]*domain.Profile, error)
}

type FollowersStorage interface {
	AddFollower(follower *domain.Follower) error
	RemoveFollower(follower *domain.Follower) error
	GetFollowingByUserId(userId int64) ([]int64, error)
}