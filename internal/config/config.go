package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
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
	switch first {
	case "-h":
		printHelp()
	case "-c", "-context":
		dir, err := FindRoot()
		if err == nil && dir != "" {
			fmt.Println("squid already initialized in this directory")
			return nil
		}
		return os.Mkdir(".squid", DirPerm)
	default:
		dir, err := FindRoot()
		if err == nil && dir != "" {
			fmt.Println("squid already initialized in this directory")
			return nil
		}
		return os.Mkdir(".squid", DirPerm)
	}
	return nil
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
