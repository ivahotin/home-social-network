package storage

import (
	"database/sql"
	"example.com/social/internal/domain"
	"strings"
)

const (
	insertStmt = "insert into profiles (username, password, firstname, lastname, age, gender, interests, city) values (?, ?, ?, ?, ?, ?, ?, ?) on duplicate key update username = username"
	getProfileByUsernameStmt = "select id, username, password, firstname, lastname, age, gender, interests, city from profiles where username = ?"
	getProfilesBySearchTerm = "select id, username, password, firstname, lastname, age, gender, interests, city from profiles where (firstname like ? or lastname like ?) and id > ? order by id asc limit ?"
	getProfilesByUserIds = "select id, username, password, firstname, lastname, age, gender, interests, city from profiles where id in "
	getProfileByUserId = "select id, username, password, firstname, lastname, age, gender, interests, city from profiles where id = ?"
)

type MySqlProfileStorage struct {
	db *sql.DB
}

func NewMySqlProfileStorage(db *sql.DB) *MySqlProfileStorage {
	return &MySqlProfileStorage{
		db: db,
	}
}

func (profileStorage *MySqlProfileStorage) GetProfileByUsername(username string) (*domain.Profile, error) {
	stmt, err := profileStorage.db.Prepare(getProfileByUsernameStmt)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var profile domain.Profile
	err = stmt.QueryRow(username).Scan(
		&profile.Id,
		&profile.Username,
		&profile.Password,
		&profile.Firstname,
		&profile.Lastname,
		&profile.Age,
		&profile.Gender,
		&profile.Interests,
		&profile.City)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &profile, nil
}

func (profileStorage *MySqlProfileStorage) SaveProfile(profile *domain.Profile) error {
	stmt, err := profileStorage.db.Prepare(insertStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(
		profile.Username,
		profile.Password,
		profile.Firstname,
		profile.Lastname,
		profile.Age,
		profile.Gender,
		profile.Interests,
		profile.City)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected > 0 {
		return nil
	} else {
		return domain.SuchUsernameExists
	}
}

func (profileStorage *MySqlProfileStorage) GetProfilesBySearchTerm(
	term string,
	offset int64,
	limit int) ([]*domain.Profile, error) {

	profiles := make([]*domain.Profile, 0, limit)
	stmt, err := profileStorage.db.Prepare(getProfilesBySearchTerm)
	if err != nil {
		return profiles, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(
		term + "%",
		term + "%",
		offset,
		limit)
	if err != nil {
		return profiles, err
	}
	defer rows.Close()

	for rows.Next() {
		profile := new(domain.Profile)

		if err = rows.Scan(
			&profile.Id,
			&profile.Username,
			&profile.Password,
			&profile.Firstname,
			&profile.Lastname,
			&profile.Age,
			&profile.Gender,
			&profile.Interests,
			&profile.City); err != nil {

			return nil, err
		}

		profiles = append(profiles, profile)
	}

	return profiles, nil
}

func (profileStorage *MySqlProfileStorage) GetProfilesByIds(userIds []int64) ([]*domain.Profile, error) {
	profiles := make([]*domain.Profile, 0)
	if len(userIds) == 0 {
		return profiles, nil
	}
	query := getProfilesByUserIds + `(?` + strings.Repeat(",?", len(userIds)-1) + `)`

	stmt, err := profileStorage.db.Prepare(query)
	if err != nil {
		return profiles, err
	}
	defer stmt.Close()

	var args []interface{}
	for _, userId := range userIds {
		args = append(args, userId)
	}
	rows, err := stmt.Query(args...)
	if err != nil {
		return profiles, err
	}

	for rows.Next() {
		profile := new(domain.Profile)

		if err = rows.Scan(
			&profile.Id,
			&profile.Username,
			&profile.Password,
			&profile.Firstname,
			&profile.Lastname,
			&profile.Age,
			&profile.Gender,
			&profile.Interests,
			&profile.City); err != nil {

			return nil, err
		}

		profiles = append(profiles, profile)
	}

	return profiles, nil
}

func (profileStorage *MySqlProfileStorage) GetProfileByUserId(userId int64) (*domain.Profile, error) {
	stmt, err := profileStorage.db.Prepare(getProfileByUserId)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var profile domain.Profile
	err = stmt.QueryRow(userId).Scan(
		&profile.Id,
		&profile.Username,
		&profile.Password,
		&profile.Firstname,
		&profile.Lastname,
		&profile.Age,
		&profile.Gender,
		&profile.Interests,
		&profile.City)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &profile, nil
}