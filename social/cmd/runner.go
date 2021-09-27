package main

import (
	db2 "example.com/social/internal/infra/db"
	"example.com/social/internal/service"
	"example.com/social/internal/storage"
	"os"
	"os/signal"
	"syscall"

	"example.com/social/internal/app"
)

func main() {
	db := db2.NewDbPool()
	profileStorage := storage.NewMySqlProfileStorage(db)
	authService := service.NewAuthService(profileStorage)
	profileService := service.NewProfileService(profileStorage)
	followerStorage := storage.NewMysqlFollowerStorage(db)
	followerService := service.NewFriendshipService(followerStorage)
	application := app.NewApplication(
		authService,
		authService,
		profileService,
		profileService,
		followerService,
		followerService,
		profileService,
		followerService,
		profileService)

	go application.Run()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	application.Stop()
}
