package domain

import "errors"

type GenderType byte

const (
	Male GenderType = 0
	Female      	= 1
)

type Profile struct {
	Username 	string
	Password 	string
	Firstname 	string
	Lastname 	string
	Age 		int
	Gender 		GenderType
	Interests 	string
	City  		string
}

var ProfileNotFound = errors.New("profile not found")