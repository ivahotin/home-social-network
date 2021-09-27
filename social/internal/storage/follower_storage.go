package storage

import (
	"database/sql"
	"example.com/social/internal/domain"
)

const (
	insertFollowersStmt = "insert into followers (follower_id, user_id, is_active) values (?, ?, ?) on duplicate key update is_active = true"
	getFollowingByUserIdStmt = "select user_id from followers where follower_id = ? and is_active"
	updateFollowerStmt = "update followers set is_active = false where follower_id = ? and user_id = ?"
)

type MysqlFollowerStorage struct {
	db *sql.DB
}

func NewMysqlFollowerStorage(db *sql.DB) *MysqlFollowerStorage {
	return &MysqlFollowerStorage{
		db: db,
	}
}

func (followerStorage *MysqlFollowerStorage) AddFollower(follower *domain.Follower) error {
	stmt, err := followerStorage.db.Prepare(insertFollowersStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(follower.FollowerId, follower.UserId, follower.IsActive)
	if err != nil {
		return err
	}

	return nil
}

func (followerStorage *MysqlFollowerStorage) RemoveFollower(follower *domain.Follower) error {
	stmt, err := followerStorage.db.Prepare(updateFollowerStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(follower.FollowerId, follower.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (followerStorage *MysqlFollowerStorage) GetFollowingByUserId(userId int64) ([]int64, error) {
	following := make([]int64, 0, 10)

	stmt, err := followerStorage.db.Prepare(getFollowingByUserIdStmt)
	if err != nil {
		return following, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(userId)
	if err != nil {
		return following, err
	}
	defer rows.Close()

	for rows.Next() {
		var friendId int64

		if err = rows.Scan(&friendId); err != nil {
			return nil, err
		}

		following = append(following, friendId)
	}

	return following, nil
}