package main

import (
	"cligram/internal/cli"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	switch os.Args[1] {
	case "user":
		cli.UserCmd(os.Args[2:])
	case "chat":
		cli.ChatCmd(os.Args[2:])
	case "msg":
		cli.MsgCmd(os.Args[2:])
	default:
		printUsage()
	}
}

func printUsage() {
	println("Usage: cligram <user|chat|msg> ...")
}
