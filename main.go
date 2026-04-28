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

	store, err := db.InitDB()
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		_ = store.Close()
	}()

	app := api.New(store)
	app.Init()

	if err := server.Run(); err != nil {
		log.Println(err)
		return
	}
}
