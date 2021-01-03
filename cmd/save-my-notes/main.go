package main

import (
	"github.com/semihsemih/save-my-notes/internal/controllers"
	"github.com/semihsemih/save-my-notes/internal/db"
	"github.com/semihsemih/save-my-notes/internal/router"
	"github.com/subosito/gotenv"
	"log"
	"net/http"
)

func init() {
	gotenv.Load()
}

func main() {
	database := db.ConnectDB()
	controller := controllers.Controller{DB: database}
	r := router.Init(&controller)
	log.Println("Listen on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", r))
}
