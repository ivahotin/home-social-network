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
	chatDb := db2.NewChatDbPool()
	profileStorage := storage.NewMySqlProfileStorage(db)
	authService := service.NewAuthService(profileStorage)
	profileService := service.NewProfileService(profileStorage)
	followerStorage := storage.NewMysqlFollowerStorage(db)
	followerService := service.NewFriendshipService(followerStorage)
	chatStorage := storage.NewCockroachChatStorage(chatDb)
	chatService := service.NewChatService(chatStorage)
	application := app.NewApplication(
		authService,
		authService,
		profileService,
		profileService,
		followerService,
		followerService,
		profileService,
		followerService,
		profileService,
		chatService)

	go application.Run()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	application.Stop()
}
