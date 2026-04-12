package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AgustinPagotto/squid/internal/cli"
	"github.com/AgustinPagotto/squid/internal/config"
	"github.com/AgustinPagotto/squid/internal/context"
	"github.com/AgustinPagotto/squid/internal/notes"
	"github.com/AgustinPagotto/squid/internal/todos"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("please add a subcommand or enter -h flag to show info")
		return
	}
	switch os.Args[1] {
	case "init":
		if len(os.Args) > 2 && os.Args[2] == "-h" {
			cli.PrintInitHelp()
			return
		}
		if err := config.Init(); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Init run sucessfuly, .squid folder created")
	case "-h":
		cli.PrintHelp()
		return
	case "notes":
		dir, err := config.FindRoot()
		if err != nil {
			if errors.Is(err, config.ErrRootNotFound) {
				fmt.Println("squid is not initialized, run 'squid init'")
			} else {
				fmt.Println(err)
			}
			return
		}
		ns := &notes.NoteStorage{Path: filepath.Join(dir, "notes.json")}
		notes.Handle(os.Args[2:], ns)
		return
	case "todos":
		dir, err := config.FindRoot()
		if err != nil {
			if errors.Is(err, config.ErrRootNotFound) {
				fmt.Println("squid is not initialized, run 'squid init'")
			} else {
				fmt.Println(err)
			}
			return
		}
		ts := &todos.TodoStorage{Path: filepath.Join(dir, "todos.json")}
		todos.Handle(os.Args[2:], ts)
		return
	case "context":
		dir, err := config.FindRoot()
		if err != nil {
			if errors.Is(err, config.ErrRootNotFound) {
				fmt.Println("squid is not initialized, run 'squid init'")
			} else {
				fmt.Println(err)
			}
			return
		}
		cs := &context.ContextStorage{Path: filepath.Join(dir, "context.json")}
		context.Handle(os.Args[2:], cs)
		return
	default:
		fmt.Println("subcommand not valid")
		return
	}
}
