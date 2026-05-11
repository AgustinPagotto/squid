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

func FindShellConfigurationFile(configFileName string) (string, error) {

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	candidates := []string{filepath.Join(home, configFileName),
		filepath.Join(home, ".config", configFileName),
		filepath.Join(home, ".config", "zsh", configFileName),
	}

	for _, c := range candidates {
		info, err := os.Stat(c)
		if err == nil && !info.IsDir() {
			return filepath.Join(c), nil
		}

	}
	return "", ErrRootNotFound
}
