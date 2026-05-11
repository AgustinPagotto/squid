package setup

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/AgustinPagotto/squid/internal/config"
)

func Init(args []string) error {
	var first string
	if len(args) > 0 {
		first = args[0]
	}
	dir, err := config.FindRoot()
	if err == nil && dir != "" {
		keepGoing := showOverrideWarning()
		if !keepGoing {
			return nil
		}
	}
	if err != nil && err != config.ErrRootNotFound {
		return err
	}
	switch first {
	case "-h":
		fmt.Print(commandHelp)
		return nil
	case "shell":
		selectedShell := selectShell()
		var shellFile, shellEvalConfig string
		if selectedShell == 2 {
			shellFile = ".zshrc"
			shellEvalConfig = zshEvalConfig
		}
		addEvalPermissions(shellFile, shellEvalConfig)
		return nil
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
	fmt.Print(initialOptionTemplate)
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
	fmt.Print(predefinedInitTemplate)
	for {
		fmt.Print("  Select [1-7]: ")
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

func addEvalPermissions(selectedShell string, shellEvalConfig string) error {
	path, err := config.FindShellConfigurationFile(selectedShell)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, config.FilePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(shellEvalConfig)
	return nil
}

func selectShell() int {
	fmt.Print(selectShellTemplate)
	for {
		fmt.Print("  Select [1-4]: ")
		var input string
		fmt.Scanln(&input)
		switch input {
		case "1", "2", "3", "4":
			n, _ := strconv.Atoi(input)
			return n
		default:
			fmt.Println("\n  Invalid option. Please enter a number between 1 and 4.")
		}
	}
}
func showOverrideWarning() bool {
	fmt.Print("Are you sure you want to override your .squid? [Y/n]: ")
	var input string
	fmt.Scanln(&input)
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "" || input == "y" || input == "yes" {
		return true
	}
	return false
}
