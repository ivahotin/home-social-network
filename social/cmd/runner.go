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
	authStorage := storage.NewMySqlAuthStorage(db)
	authService := service.NewAuthService(authStorage)
	profileService := service.NewProfileService(authStorage)
	application := app.NewApplication(authService, authService, profileService)

	go application.Run()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	application.Stop()
}
