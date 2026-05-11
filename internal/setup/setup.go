package setup

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/AgustinPagotto/squid/internal/config"
)

func Init(args []string) error {
	var first string
	if len(args) > 0 {
		first = args[0]
	}
	if first == "-h" {
		printHelp()
		return nil
	}
	err := handleInit()
	if err != nil {
		return err
	}
	return nil
}

func handleInit() error {
	dir, err := config.FindRoot()
	if err == nil && dir != "" {
		fmt.Println("squid already initialized in this directory, if you want to restart squid use squid ink command")
		return nil
	}
	if err != config.ErrRootNotFound {
		fmt.Println(err)
		return err
	}
	option := giveInitialOptions()
	switch option {
	case 1:
		option := predefinedInit()
		err := os.Mkdir(".squid", config.DirPerm)
		if err != nil {
			return err
		}
		return initiatePredefinedSquidFolder(option)
	case 2:
		fmt.Println("selected custom")
	default:
		return nil
	}
	return nil
}

func giveInitialOptions() int {
	fmt.Print(`
  ┌────────────────────────────────────────────────┐
  │                                                │
  │                 squid  ·  init                 │
  │                                                │
  ├────────────────────────────────────────────────┤
  │                                                │
  │   How would you like to set up your context?   │
  │                                                │
  │   1  ·  Predefined  ─  Node, Go, React, ...    │
  │   2  ·  Custom      ─  build your own          │
  │                                                │
  │   3  ·  Exit                                   │
  │                                                │
  └────────────────────────────────────────────────┘

`)
	for {
		fmt.Print("  Select [1-3]: ")
		var input string
		fmt.Scanln(&input)
		switch input {
		case "1", "2", "3":
			n, _ := strconv.Atoi(input)
			return n
		default:
			fmt.Println("\n  Invalid option. Please enter a number between 1 and 3.")
		}
	}
}

func predefinedInit() int {
	fmt.Print(`
  ┌────────────────────────────────────────────────┐
  │                                                │
  │              squid  ·  predefined              │
  │                                                │
  ├────────────────────────────────────────────────┤
  │                                                │
  │   Select a predefined context:                 │
  │                                                │
  │   1  ·  Go                                     │
  │   2  ·  Node.js                                │
  │   3  ·  Python                                 │
  │   4  ·  React                                  │
  │   5  ·  React Native                           │
  │   6  ·  Next.js                                │
  │                                                │
  │   7  ·  ← Back                                 │
  │   8  ·  Exit                                   │
  │                                                │
  └────────────────────────────────────────────────┘

`)

	for {
		fmt.Print("  Select [1-14]: ")
		var input string
		fmt.Scanln(&input)
		switch input {
		case "1", "2", "3", "4", "5", "6", "7",
			"8":
			n, _ := strconv.Atoi(input)
			return n
		default:
			fmt.Println("\n  Invalid option. Please enter a number between 1 and 8.")
		}
	}
}

func printHelp() {
	fmt.Print(`Initialize squid in the current directory

Usage:
  squid init

Creates a .squid/ directory in the current directory. Run this once
per project. All squid data (notes, todos, context) will be stored there.
`)
}

func initiatePredefinedSquidFolder(initOption int) error {
	var languageNameFile string
	switch initOption {
	case 1:
		languageNameFile = "go.json"
	case 2:
		languageNameFile = "nodejs.json"
	case 3:
		languageNameFile = "python.json"
	case 4:
		languageNameFile = "react.json"
	case 5:
		languageNameFile = "react-native.json"
	case 6:
		languageNameFile = "nextjs.json"
	}
	data, err := predefinedFS.ReadFile("assets/predefined/" + languageNameFile)
	if err != nil {
		return err
	}
	dest := filepath.Join(".squid", "context.json")
	return os.WriteFile(dest, data, config.FilePerm)
}
