package main

import (
	"cligram/internal/cli"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		printUsage()
		return
	}

	commands := map[string]func([]string){
		"user":        cli.UserCmd,
		"chat":        cli.ChatCmd,
		"msg":         cli.MsgCmd,
		"interactive": interactiveCmd,
	}

	cmd := args[0]
	if handler, ok := commands[cmd]; ok {
		handler(args[1:])
	} else {
		printUsage()
	}
}

func interactiveCmd(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: cligram interactive <user_id> <server_addr>")
		return
	}
	userID := args[0]
	serverAddr := args[1]
	cli.InteractiveChat(userID, serverAddr)
}

func printUsage() {
	fmt.Println("Usage: cligram <user|chat|msg|interactive> ...")
}
