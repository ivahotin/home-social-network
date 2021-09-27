package usecases

type FollowUseCase interface {
	Follow(followerId, userId int64) error
}

type UnFollowUseCase interface {
	Unfollow(followerId, userId int64) error
}

type GetFollowingByUserIdQuery interface {
	GetFollowingByUserId(userId int64) ([]int64, error)
}