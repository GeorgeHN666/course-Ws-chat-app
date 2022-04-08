package main

import (
	"log"
	"net/http"

	"github.com/GeorgeHN666/course-ws-chat-app/internal/handlers"
)

func main() {

	mux := routes()

	log.Println("starting channel listener")

	go handlers.ListenToWsChannel()

	log.Println("Starting Server on PORT::8080")

	_ = http.ListenAndServe(":8080", mux)

}
