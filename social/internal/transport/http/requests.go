package http

type SignInRequest struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type SignUpRequest struct {
	Username 	string 	`form:"username" json:"username" binding:"required"`
	Password 	string 	`form:"password" json:"password" binding:"required"`
	Firstname 	string 	`form:"firstname" json:"firstname" binding:"required"`
	Lastname 	string 	`form:"lastname" json:"lastname" binding:"required"`
	Age 		int 	`form:"age" json:"age" binding:"required"`
	Gender 		string 	`form:"gender" json:"gender" binding:"required"`
	Interests 	string 	`form:"interests" json:"interests" binding:"required"`
	City 		string 	`form:"city" json:"city" binding:"required"`
}

type SearchUsersRequest struct {
	SearchTerm	string  `form:"term" binding:"required"`
	Cursor      int64	`form:"cursor,default=0"`
	Limit       int     `form:"limit" binding:"required,min=5,max=100"`
}

type FollowRequest struct {
	FollowedId    int64   `uri:"followed" binding:"required"`
}

type UnfollowRequest struct {
	FollowedId    int64   `uri:"followed" binding:"required"`
}