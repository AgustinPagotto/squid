package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var ErrRootNotFound = errors.New("no .squid directory found in current directory")
var ErrConfigFileNotFound = errors.New("no config file found on the system")
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

func FindShellConfigurationFile(configFileName, home string) (string, error) {
	candidates := []string{
		filepath.Join(home, configFileName),
		filepath.Join(home, ".config", configFileName),
		filepath.Join(home, ".config", "zsh", configFileName),
	}

	for _, c := range candidates {
		info, err := os.Stat(c)
		if err == nil && !info.IsDir() {
			return filepath.Join(c), nil
		}

	}
	return "", ErrConfigFileNotFound
}

func OpenInEditor(helpText string) (string, error) {
	tmpFile, err := os.CreateTemp("", "squid-temp-*.txt")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())

	if _, err = tmpFile.WriteString(helpText); err != nil {
		return "", err
	}
	tmpFile.Close()

	cmd := exec.Command(getEditor(), tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return "", err
	}

	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func getEditor() string {
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	return "nano"
}
