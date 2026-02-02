package cli

import (
	"cligram/internal/app"
	"fmt"
)

func userCmd(service *app.ChatService, args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: cligram user create <id> <name>")
		return
	}
	switch args[0] {
	case "create":
		if len(args) != 3 {
			fmt.Println("Usage: cligram user create <id> <name>")
			return
		}
		id := args[1]
		name := args[2]
		if err := service.CreateUser(id, name); err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("User created:", id)
	default:
		fmt.Println("Unknown user command:", args[0])
	}
}
