package cli

import (
	"cligram/internal/app"
	"fmt"
	"os"
)

type CLI struct {
	Service *app.ChatService
}

func NewCLI(service *app.ChatService) *CLI {
	return &CLI{Service: service}
}

func (c *CLI) Execute() {
	if len(os.Args) < 2 {
		c.usage()
		return
	}

	switch os.Args[1] {
	case "user":
		userCmd(c.Service, os.Args[2:])
	case "chat":
		chatCmd(c.Service, os.Args[2:])
	case "msg":
		msgCmd(c.Service, os.Args[2:])
	default:
		fmt.Println("Unknown command:", os.Args[1])
		c.usage()
	}
}

func (c *CLI) usage() {
	fmt.Println(`Usage:
  cligram user create <id> <name>
  cligram chat create <id> <member1,member2,...>
  cligram msg send <from> <chat> <text>
  cligram msg list <user> <chat>`)
	os.Exit(1)
}
