package main

import (
	"Goose47/storage/config"
	"Goose47/storage/db"
	"Goose47/storage/server"
	"Goose47/storage/tasks"
)

func main() {
	config.Init()
	db.Init()
	tasks.Init()
	server.Init()
}
