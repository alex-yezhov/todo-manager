package main

import (
	"log"

	"github.com/joho/godotenv"

	"scheduler/pkg/api"
	"scheduler/pkg/db"
	"scheduler/pkg/server"
)

func main() {
	_ = godotenv.Load()

	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = database.Close()
	}()

	api.Init()

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
