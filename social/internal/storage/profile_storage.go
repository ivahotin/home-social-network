package storage

import (
	"database/sql"
	"example.com/social/internal/domain"
	"sort"
	"strings"
	"time"
)

const (
	insertStmt = "insert into profiles (username, password, firstname, lastname, birthdate, gender, interests, city) values (?, ?, ?, ?, ?, ?, ?, ?) on duplicate key update username = username"
	getProfileByUsernameStmt = "select id, username, password, firstname, lastname, birthdate, gender, interests, city from profiles where username = ?"
	getProfilesBySearchTerm = "select id, username, password, firstname, lastname, birthdate, gender, interests, city from profiles where (firstname like ? and lastname like ?) and id > ? limit ?"
	getProfilesByUserIds = "select id, username, password, firstname, lastname, birthdate, gender, interests, city from profiles where id in "
	getProfileByUserId = "select id, username, password, firstname, lastname, birthdate, gender, interests, city from profiles where id = ?"
)

type MySqlProfileStorage struct {
	db *sql.DB
}

func NewMySqlProfileStorage(db *sql.DB) *MySqlProfileStorage {
	return &MySqlProfileStorage{
		db: db,
	}
}

type ById []*domain.Profile

func (a ById) Len() int { return len(a) }
func (a ById) Less(i, j int) bool { return a[i].Id < a[j].Id }
func (a ById) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (profileStorage *MySqlProfileStorage) GetProfileByUsername(username string) (*domain.Profile, error) {
	stmt, err := profileStorage.db.Prepare(getProfileByUsernameStmt)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var profile domain.Profile
	var birthdate string
	err = stmt.QueryRow(username).Scan(
		&profile.Id,
		&profile.Username,
		&profile.Password,
		&profile.Firstname,
		&profile.Lastname,
		&birthdate,
		&profile.Gender,
		&profile.Interests,
		&profile.City)

	profile.Birthdate, err = time.Parse("2006-01-02", birthdate)
	if err != nil {
		return nil, err
	}

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
		profile.Birthdate,
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
	firstname, lastname string,
	offset int64,
	limit int) ([]*domain.Profile, error) {

	profiles := make([]*domain.Profile, 0, limit)
	var stmt *sql.Stmt
	var err error
	stmt, err = profileStorage.db.Prepare(getProfilesBySearchTerm)
	if err != nil {
		return profiles, err
	}
	defer stmt.Close()

	var rows *sql.Rows
	rows, err = stmt.Query(
		firstname + "%",
		lastname + "%",
		offset,
		limit)
	if err != nil {
		return profiles, err
	}
	defer rows.Close()

	for rows.Next() {
		profile := new(domain.Profile)

		var birthdate string
		if err = rows.Scan(
			&profile.Id,
			&profile.Username,
			&profile.Password,
			&profile.Firstname,
			&profile.Lastname,
			&birthdate,
			&profile.Gender,
			&profile.Interests,
			&profile.City); err != nil {

			return nil, err
		}
		profile.Birthdate, err = time.Parse("2006-01-02", birthdate)
		if err != nil {
			return nil, err
		}

		profiles = append(profiles, profile)
	}

	sort.Sort(ById(profiles))

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

		var birthdate string
		if err = rows.Scan(
			&profile.Id,
			&profile.Username,
			&profile.Password,
			&profile.Firstname,
			&profile.Lastname,
			&birthdate,
			&profile.Gender,
			&profile.Interests,
			&profile.City); err != nil {

			return nil, err
		}
		profile.Birthdate, err = time.Parse("2006-01-02", birthdate)
		if err != nil {
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
	var birthdate string
	err = stmt.QueryRow(userId).Scan(
		&profile.Id,
		&profile.Username,
		&profile.Password,
		&profile.Firstname,
		&profile.Lastname,
		&birthdate,
		&profile.Gender,
		&profile.Interests,
		&profile.City)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	profile.Birthdate, err = time.Parse("2006-01-02", birthdate)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}