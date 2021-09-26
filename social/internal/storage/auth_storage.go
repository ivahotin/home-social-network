package storage

import (
	"database/sql"
	"example.com/social/internal/domain"
)

const (
	insertStmt = "insert into profiles (username, password, firstname, lastname, age, gender, interests, city) values (?, ?, ?, ?, ?, ?, ?, ?) on duplicate key update username = username"
	getProfileByUsernameStmt = "select username, password, firstname, lastname, age, gender, interests, city from profiles where username = ?"
)

type MySqlAuthStorage struct {
	db *sql.DB
}

func NewMySqlAuthStorage(db *sql.DB) *MySqlAuthStorage {
	return &MySqlAuthStorage{
		db: db,
	}
}

func (authStorage *MySqlAuthStorage) SaveProfile(profile *domain.Profile) error {
	stmt, err := authStorage.db.Prepare(insertStmt)
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

func (authStorage *MySqlAuthStorage) GetProfileByUsername(username string) (*domain.Profile, error) {
	stmt, err := authStorage.db.Prepare(getProfileByUsernameStmt)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var profile domain.Profile
	err = stmt.QueryRow(username).Scan(
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
			return nil, nil
		}
		return nil, err
	}

	return &profile, nil
}