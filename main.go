package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

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
			fmt.Print(`Initialize squid in the current directory

Usage:
  squid init

Creates a .squid/ directory in the current directory. Run this once
per project. All squid data (notes, todos, context) will be stored there.
`)
			return
		}
		if err := config.Init(); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Init run sucessfuly, .squid folder created")
	case "-h":
		fmt.Print(`Squid — project-aware developer context CLI

Usage:
  squid <command> [arguments]

Available Commands:
  init       Initialize squid in the current directory
  notes      Manage project notes
  todo       Manage project todos
  context    Manage project context and aliases

Flags:
  -h         Show help for squid

Use "squid <command> -h" for more information about a command.
`)
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
