package main

import (
	"chat/server"
	"log"
	"net/http"
)

func main() {
	v := server.NewServer()
	log.Println("Сервак жгёт ✅ ")
	log.Fatal(http.ListenAndServe(":8081", v))
}
