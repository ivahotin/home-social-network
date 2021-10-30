package http

import (
	"example.com/social/internal/domain"
	"time"
)

type User struct {
	Id        int64
	UserName  string
}

type Profile struct {
	Id        int64  	`json:"id"`
	Username  string 	`json:"username"`
	Firstname string 	`json:"firstname"`
	Lastname  string 	`json:"lastname"`
	Birthdate time.Time `json:"birthdate"`
	Gender    string 	`json:"gender"`
	City      string 	`json:"city"`
	Interests string 	`json:"interests"`
}

type GetProfilesBySearchTerm struct {
	Profiles  	[]*Profile	`json:"profiles"`
	PrevCursor 	int64      	`json:"prev_cursor"`
	NextCursor 	int64      	`json:"next_cursor"`
}

type MeProfileResponse struct {
	Profile   	*Profile    `json:"profile"`
	Following   []*Profile  `json:"following"`
}

type GetProfileByUserIdResponse struct {
	Profile     *Profile    `json:"profile"`
}

func ConvertDomainProfileToResponseProfile(domainProfile *domain.Profile) *Profile {
	profile 			:= 	new(Profile)
	profile.Id 			= 	domainProfile.Id
	profile.Username 	= 	domainProfile.Username
	profile.Firstname 	= 	domainProfile.Firstname
	profile.Lastname 	=	domainProfile.Lastname
	profile.Birthdate 	= 	domainProfile.Birthdate
	profile.City        = 	domainProfile.City
	switch domainProfile.Gender {
	case domain.Male: profile.Gender = "Male"
	case domain.Female: profile.Gender = "Female"
	}
	profile.Interests 	= domainProfile.Interests

	return profile
}