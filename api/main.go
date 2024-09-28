package main

import (
	"Goose47/storage/config"
	"Goose47/storage/db"
	"Goose47/storage/server"
	"Goose47/storage/tasks"
)

func main() {
	cfg := config.MustLoad()
	db, err := db.New(cfg.DB)
	if err != nil {

	}
	tasks.Init()
	server.Init()
}
