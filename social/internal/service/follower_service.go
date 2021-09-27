package service

import (
	"example.com/social/internal/domain"
)

type FollowerService struct {
	FollowersStorage FollowersStorage
}

func NewFriendshipService(storage FollowersStorage) *FollowerService {
	return &FollowerService{
		storage,
	}
}

func (followerService *FollowerService) Follow(followerId, userId int64) error {
	follower := domain.Follower{
		FollowerId:   	followerId,
		UserId: 		userId,
		IsActive: 		true,
	}

	err := followerService.FollowersStorage.AddFollower(&follower)
	if err != nil {
		return nil
	}

	return nil
}

func (followerService *FollowerService) Unfollow(followerId, userId int64) error {
	follower := domain.Follower{
		FollowerId:   	followerId,
		UserId: 		userId,
		IsActive: 		false,
	}

	err := followerService.FollowersStorage.RemoveFollower(&follower)
	if err != nil {
		return nil
	}

	return nil
}

func (followerService *FollowerService) GetFollowingByUserId(userId int64) ([]int64, error) {
	return followerService.FollowersStorage.GetFollowingByUserId(userId)
}