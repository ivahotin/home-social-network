package domain

import (
	"errors"
	"time"
)

type GenderType byte

const (
	Male GenderType = 0
	Female      	= 1
)

type Profile struct {
	Id          int64
	Username 	string
	Password 	string
	Firstname 	string
	Lastname 	string
	Birthdate 	time.Time
	Gender 		GenderType
	Interests 	string
	City  		string
}

var ProfileNotFound = errors.New("profile not found")
var SuchUsernameExists = errors.New("username already exists")

type ProfilesSearchResult struct {
	Profiles   []*Profile
	PrevCursor int64
	NextCursor int64
}

func NewProfilesSearchResult(profiles []*Profile, cursor int64, myId int64) *ProfilesSearchResult {
	filteredProfiles := make([]*Profile, 0, len(profiles))
	for _, profile := range profiles {
		if profile.Id != myId {
			filteredProfiles = append(filteredProfiles, profile)
		}
	}

	var nextCursor int64 = 0
	if len(filteredProfiles) > 0 {
		nextCursor = filteredProfiles[len(profiles) - 1].Id
	}

	return &ProfilesSearchResult{
		Profiles: filteredProfiles,
		PrevCursor: cursor,
		NextCursor: nextCursor,
	}
}