package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

var ErrRootNotFound = errors.New("no .squid directory found in current directory")
var ErrNotEnoughPermissions = errors.New("not enough permissions to read the current directory")

const (
	DirPerm  os.FileMode = 0755
	FilePerm os.FileMode = 0644
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
	dir, err := FindRoot()
	if err == nil && dir != "" {
		fmt.Println("squid already initialized in this directory, if you want to restart squid use squid ink command")
		return nil
	}
	if err != ErrRootNotFound {
		fmt.Println(err)
		return err
	}
	option := giveInitialOptions()
	switch option {
	case 1:
		option := predefinedInit()
		fmt.Println(option)
	case 2:
		fmt.Println("selected custom")
	default:
		return nil
	}
	return os.Mkdir(".squid", DirPerm)
}

func FindRoot() (string, error) {
	currDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not get current directory: %w", err)
	}
	for {
		info, err := os.Stat(filepath.Join(currDir, ".squid"))
		if err == nil && info.IsDir() {
			return filepath.Join(currDir, ".squid"), nil
		}
		parent := filepath.Dir(currDir)
		if currDir == parent {
			break
		}
		currDir = parent
	}
	return "", ErrRootNotFound
}

func FindShellConfigurationFile() (string, error) {
	currDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not get current directory: %w", err)
	}
	for {
		info, err := os.Stat(filepath.Join(currDir, ".zshrc"))
		if err == nil && !info.IsDir() {
			return filepath.Join(currDir, ".zshrc"), nil
		}
		parent := filepath.Dir(currDir)
		if currDir == parent {
			break
		}
		currDir = parent
	}
	return "", ErrRootNotFound
}

func printHelp() {
	fmt.Print(`Initialize squid in the current directory

Usage:
  squid init

Creates a .squid/ directory in the current directory. Run this once
per project. All squid data (notes, todos, context) will be stored there.
`)
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
			return int(input[0] - '0')
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
  │   7  ·  Vue.js                                 │
  │   8  ·  Rust                                   │
  │   9  ·  Django                                 │
  │  10  ·  Rails                                  │
  │  11  ·  Flutter                                │
  │  12  ·  Docker                                 │
  │                                                │
  │  13  ·  ← Back                                 │
  │  14  ·  Exit                                   │
  │                                                │
  └────────────────────────────────────────────────┘

`)

	for {
		fmt.Print("  Select [1-14]: ")
		var input string
		fmt.Scanln(&input)
		switch input {
		case "1", "2", "3", "4", "5", "6", "7",
			"8", "9", "10", "11", "12", "13", "14":
			n, _ := strconv.Atoi(input)
			return n
		default:
			fmt.Println("\n  Invalid option. Please enter a number between 1 and 14.\n")
		}
	}
}
