package app

import (
	"context"
	"example.com/social/internal/usecases"
	"fmt"
	"net/http"
	"os"
	"time"

	transport "example.com/social/internal/transport/http"
)

type Application struct {
	srv *http.Server
}

func NewApplication(
	signUpUseCase usecases.SignUpUseCase,
	signInUseCase usecases.SignInUseCase,
	getProfileByUsername usecases.GetProfileGetUsernameUseCase,
	getProfilesBySearchTerm usecases.GetProfilesBySearchTerm,
	followUserCase usecases.FollowUseCase,
	getFollowingByUserIdQuery usecases.GetFollowingByUserIdQuery,
	getProfilesByUserIdsQuery usecases.GetProfilesByUserIdsQuery,
	unfollowUseCase usecases.UnFollowUseCase,
	getUserProfileByUserIdQuery usecases.GetProfileByUserIdQuery) *Application {
	return &Application{
		transport.NewServer(
			":" + os.Getenv("PORT"),
			transport.MakeEndpoints(
				signUpUseCase,
				signInUseCase,
				getProfileByUsername,
				getProfilesBySearchTerm,
				followUserCase,
				getFollowingByUserIdQuery,
				getProfilesByUserIdsQuery,
				unfollowUseCase,
				getUserProfileByUserIdQuery)),
	}
}

func (app *Application) Run() {
	go func() {
		if err := app.srv.ListenAndServe(); err != nil {
			fmt.Println("Http server error, ", err)
		}
	}()
}

func (app *Application) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	if err := app.srv.Shutdown(ctx); err != nil {
		fmt.Println("Http closing error", err)
	}
}
