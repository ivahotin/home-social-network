package http

type User struct {
	Id        int64
	UserName  string
}

type Profile struct {
	Id        int64  `json:"id"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	City      string `json:"city"`
	Interests string `json:"interests"`
}

type GetProfilesBySearchTerm struct {
	Profiles  	[]*Profile	`json:"profiles"`
	PrevCursor 	int64      	`json:"prev_cursor"`
	NextCursor 	int64      	`json:"next_cursor"`
}