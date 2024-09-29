package main

import (
	"Goose47/storage/internal/api/controllers"
	"Goose47/storage/internal/api/services"
	"Goose47/storage/internal/config"
	database "Goose47/storage/internal/db"
	"Goose47/storage/internal/server"
	"Goose47/storage/internal/tasks"
	"Goose47/storage/internal/utils/repositories"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"os"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Mode)
	db, err := database.New(cfg.DB)
	if err != nil {
		log.Error(fmt.Sprintf("could not connect to database: %s", err.Error()))
		return
	}

	taskManager := tasks.New()
	taskManager.RunTasks(&db, cfg.FS)

	repo := repositories.NewStorageRepository(db.GetCollection())

	itemService := services.NewItemService(log, repo, repo, repo, cfg.FS.Base)
	permsService, err := services.NewPermsService(log, cfg.AuthAddress)
	if err != nil {
		log.Error(fmt.Sprintf("could not create perms service: %s", err.Error()))
		return
	}

	controller := controllers.NewStorageController(itemService)

	router := server.NewRouter(log, controller, cfg.Secret, permsService)

	err = server.Serve(cfg, router) // blocking
	if err != nil {
		log.Error(fmt.Sprintf("server stopped: %s", err.Error()))
	}

	taskManager.StopTasks()
	//todo add grpc auth
}

const (
	debugMode   = gin.DebugMode
	releaseMode = gin.ReleaseMode
)

func setupLogger(mode string) *slog.Logger {
	var log *slog.Logger

	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	switch mode {
	case debugMode:
		log = slog.New(slog.NewTextHandler(os.Stdout, opts))
	case releaseMode:
		log = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	}

	return log
}
