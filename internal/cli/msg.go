package cli

import (
	"bytes"
	"cligram/internal/domain"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func MsgCmd(args []string) {
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

		reqBody, _ := json.Marshal(map[string]string{
			"from":    from,
			"chat_id": chatID,
			"text":    text,
		})

		resp, err := http.Post(serverURL+"/messages", "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			fmt.Println("Request error:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			fmt.Println("Error:", resp.Status)
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

		url := fmt.Sprintf("%s/messages?user_id=%s&chat_id=%s", serverURL, user, chatID)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Request error:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("Error:", resp.Status)
			return
		}

		var msgs []domain.Message

		if err := json.NewDecoder(resp.Body).Decode(&msgs); err != nil {
			fmt.Println("Decode error:", err)
			return
		}
		for _, m := range msgs {
			fmt.Printf("[%s] %s: %s\n", m.CreatedAt.Format("15:04:05"), m.From, m.Text)
		}

	default:
		fmt.Println("Unknown msg command:", args[0])
	}
}
