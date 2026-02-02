package cli

import (
	"flag"
	"fmt"
)

func UserCmd(args []string) {
	if len(args) < 1 {
		fmt.Println("Expected user subcommand: create")
		return
	}

	switch args[0] {
	case "create":
		fs := flag.NewFlagSet("create", flag.ExitOnError)
		id := fs.String("id", "", "user ID")
		name := fs.String("name", "", "user name")
		fs.Parse(args[1:])

		if *id == "" || *name == "" {
			fmt.Println("id and name are required")
			return
		}

		postJSON("/users", map[string]string{"id": *id, "name": *name})

	default:
		fmt.Println("Unknown user subcommand:", args[0])
	}
}
