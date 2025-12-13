package main

import (
	"chat/db"
	"chat/server"
	"log"
	"net/http"
)

func main() {

	db.InitDB()

	m := server.NewServer()

	log.Println("🔥 Сервак зажог!")
	log.Fatal(http.ListenAndServe(":8081", m))
}
