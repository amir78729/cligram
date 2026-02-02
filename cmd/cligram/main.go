package main

import (
	"cligram/internal/app"
	"cligram/internal/cli"
	"cligram/internal/db"
)

func main() {
	client := db.Connect()

	users := db.NewUserRepo(client)
	chats := db.NewChatRepo(client)
	messages := db.NewMessageRepo(client)

	service := app.NewChatService(users, chats, messages)

	c := cli.NewCLI(service)
	c.Execute()
}
