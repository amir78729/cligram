package cli

import (
	"cligram/internal/app"
	"fmt"
	"strings"
)

func chatCmd(service *app.ChatService, args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: cligram chat create <id> <member1,member2,...>")
		return
	}
	switch args[0] {
	case "create":
		if len(args) != 3 {
			fmt.Println("Usage: cligram chat create <id> <member1,member2,...>")
			return
		}
		id := args[1]
		members := strings.Split(args[2], ",")
		if err := service.CreateChat(id, members); err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Chat created:", id)
	default:
		fmt.Println("Unknown chat command:", args[0])
	}
}
