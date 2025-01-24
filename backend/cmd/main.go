package main

import (
	"log"

	"github.com/jeagerism/medium-clone/backend/config"
	"github.com/jeagerism/medium-clone/backend/pkg/database"
	"github.com/jeagerism/medium-clone/backend/server"
)

func main() {
	conf := config.GetConfig()
	if conf == nil {
		log.Fatalf("Failed to load configuration")
	}

	db, err := database.NewPostgresDatabase(conf)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.GetDb().Close()

	log.Println("Database initialized successfully")

	server.NewGinServer(conf, db).Start()

	log.Println("Application started successfully")
}
