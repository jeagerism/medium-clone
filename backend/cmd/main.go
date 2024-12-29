package main

import (
	"github.com/jeagerism/medium-clone/backend/config"
	"github.com/jeagerism/medium-clone/backend/pkg/database"
	"github.com/jeagerism/medium-clone/backend/server"
)

func main() {
	cfg := config.GetConfig()

	db := database.NewPostgresDatabase(cfg)

	server.NewGInServer(cfg, db).Start()
}
