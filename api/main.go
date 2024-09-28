package main

import (
	"Goose47/storage/config"
	database "Goose47/storage/db"
	"Goose47/storage/server"
	"Goose47/storage/tasks"
)

func main() {
	cfg := config.MustLoad()
	db, err := database.New(cfg.DB)
	if err != nil {

	}
	taskManager := tasks.New()
	taskManager.RunTasks(db, cfg.FS)
	server.Init()
}
