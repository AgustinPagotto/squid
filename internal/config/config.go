package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var ErrRootNotFound = errors.New("no .squid directory found in current directory")
var ErrNotEnoughPermissions = errors.New("not enough permissions to read the current directory")

func Init() error {
	exists, err := checkInitFolderExists()
	if err != nil {
		return err
	}
	if exists {
		fmt.Println("squid already initialized in this directory")
		return nil
	}
	return os.Mkdir(".squid", 0755)
}

func CheckRoot() error {
	exists, err := checkInitFolderExists()
	if err != nil {
		return err
	}
	if !exists {
		return ErrRootNotFound
	}
	return nil
}

func checkInitFolderExists() (bool, error) {
	c, err := os.ReadDir(".")
	if err != nil {
		if os.IsPermission(err) {
			return false, ErrNotEnoughPermissions
		}
		return false, fmt.Errorf("could not read current directory: %w", err)
	}
	for _, dir := range c {
		if dir.Name() == ".squid" && dir.IsDir() {
			return true, nil
		}
	}
	return false, nil
}

func FindRoot() (string, error) {
	currDir, _ := os.Getwd()
	for {
		directories, err := os.ReadDir(currDir)
		if err != nil {
			return "", err
		}
		for _, dir := range directories {
			if dir.Name() == ".squid" && dir.IsDir() {
				return currDir, nil
			}
		}
		parent := filepath.Dir(currDir)
		if currDir == parent {
			break
		}
		currDir = parent
	}
	return "", ErrRootNotFound
}
