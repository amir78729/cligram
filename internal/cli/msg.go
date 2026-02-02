package cli

import (
	"cligram/internal/app"
	"fmt"
	"strings"
)

func msgCmd(service *app.ChatService, args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: cligram msg send <from> <chat> <text> | msg list <user> <chat>")
		return
	}

	switch args[0] {
	case "send":
		if len(args) < 4 {
			fmt.Println("Usage: cligram msg send <from> <chat> <text>")
			return
		}
		from := args[1]
		chatID := args[2]
		text := strings.Join(args[3:], " ")
		if err := service.SendMessage(from, chatID, text); err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Message sent")
	case "list":
		if len(args) != 3 {
			fmt.Println("Usage: cligram msg list <user> <chat>")
			return
		}
		user := args[1]
		chatID := args[2]
		msgs, err := service.ListMessages(user, chatID)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		for _, m := range msgs {
			fmt.Printf("[%s] %s: %s\n", m.CreatedAt.Format("15:04"), m.From, m.Text)
		}
	default:
		fmt.Println("Unknown msg command:", args[0])
	}
}
