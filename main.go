package main

import (
	"fmt"
	"os"

	"github.com/AgustinPagotto/squid/internal/cli"
	"github.com/AgustinPagotto/squid/internal/notes"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("please add a subcommand or enter -h flag to show info")
		return
	}
	switch os.Args[1] {
	case "-h":
		cli.PrintHelp()
		return
	case "notes":
		ns := &notes.NoteStorage{Path: "notes.json"}
		notes.Handle(os.Args[2:], ns)
		return
	case "context":
		fmt.Println("context")
		return
	default:
		fmt.Println("subcommand not valid")
		return
	}
}
