package cli

import (
	"flag"
	"fmt"
	"strings"
)

func ChatCmd(args []string) {
	if len(args) < 1 {
		fmt.Println("Expected chat subcommand: create")
		return
	}

	switch args[0] {
	case "create":
		fs := flag.NewFlagSet("create", flag.ExitOnError)
		id := fs.String("id", "", "chat ID")
		members := fs.String("members", "", "comma-separated user IDs")
		fs.Parse(args[1:])

		if *id == "" || *members == "" {
			fmt.Println("id and members are required")
			return
		}

		memberList := strings.Split(*members, ",")
		postJSON("/chats", map[string]interface{}{"id": *id, "members": memberList})

	default:
		fmt.Println("Unknown chat subcommand:", args[0])
	}
}
