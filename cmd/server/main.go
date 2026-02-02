package main

import (
	"cligram/cmd/server/api"
	"cligram/internal/app"
	"cligram/internal/db"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	service *app.ChatService
}

func main() {
	client := db.Connect()

	users := db.NewUserRepo(client)
	chats := db.NewChatRepo(client)
	messages := db.NewMessageRepo(client)

	service := app.NewChatService(users, chats, messages)
	server := &api.Server{Service: service}

	r := mux.NewRouter()

	// User endpoints
	r.HandleFunc("/users", server.CreateUserHandler).Methods("POST")

	// Chat endpoints
	r.HandleFunc("/chats", server.CreateChatHandler).Methods("POST")

	// Message endpoints
	r.HandleFunc("/messages", server.SendMessageHandler).Methods("POST")
	r.HandleFunc("/messages", server.ListMessagesHandler).Methods("GET")

	httpHandler := api.LoggingMiddleware(r)

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", httpHandler))
}
