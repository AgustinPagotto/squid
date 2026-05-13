package setup

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/AgustinPagotto/squid/internal/config"
)

func Init(args []string) error {
	var firstArg string
	if len(args) > 0 {
		firstArg = args[0]
	}
	switch firstArg {
	case "-h":
		fmt.Print(commandHelp)
		return nil
	case "shell":
		if exited := handleShell(); !exited {
			fmt.Print(successTemplate)
		}
		return nil
	case "":
		dir, err := config.FindRoot()
		var overrideConfirmed bool
		if err == nil && dir != "" {
			overrideConfirmed = showOverrideWarning()
			if !overrideConfirmed {
				return nil
			}
		}
		if err != nil && err != config.ErrRootNotFound {
			return err
		}
		option := giveInitialOptions()
		switch option {
		case 1:
			option, action := predefinedInit()
			if action == ActionExit {
				fmt.Println("  squid init exited")
				return nil
			}
			if overrideConfirmed {
				os.RemoveAll(".squid")
			}
			err := os.Mkdir(".squid", config.DirPerm)
			if err != nil {
				return err
			}
			if err := initiatePredefinedSquidFolder(option); err != nil {
				return err
			}
			continueWithShell := shellNextDialog()
			if continueWithShell {
				handleShell()
			}
			fmt.Print(successTemplate)
			return nil
		case 2:
			fmt.Println("selected custom")
		default:
			return nil
		}
	default:
		fmt.Println("Not valid init subcommand")
	}
	return nil
}

func handleShell() (exited bool) {
	selectedShell, action := selectShell()
	if action == ActionExit {
		return true
	}
	addEvalPermissions(selectedShell)
	return false
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

func predefinedInit() (Language, Action) {
	fmt.Print(predefinedInitTemplate)
	for {
		fmt.Print("  Select [1-7]: ")
		var input string
		fmt.Scanln(&input)
		switch input {
		case "1", "2", "3", "4", "5", "6":
			n, _ := strconv.Atoi(input)
			return Language(n), ActionSelect
		case "7":
			return 0, ActionExit
		default:
			fmt.Println("\n  Invalid option. Please enter a number between 1 and 7.")
		}
	}
}

func initiatePredefinedSquidFolder(initOption Language) error {
	data, err := predefinedFS.ReadFile("assets/predefined/" + initOption.filename())
	if err != nil {
		return err
	}
	dest := filepath.Join(".squid", "context.json")
	return os.WriteFile(dest, data, config.FilePerm)
}

func addEvalPermissions(selectedShell Shell) error {
	path, err := config.FindShellConfigurationFile(selectedShell.configFile())
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, config.FilePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	existing, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	if strings.Contains(string(existing), selectedShell.evalString()) {
		fmt.Println("  Shell hook already present, skipping.")
		return nil
	}

	_, err = f.WriteString(selectedShell.evalString())
	return err
}

func selectShell() (Shell, Action) {
	fmt.Print(selectShellTemplate)
	for {
		fmt.Print("  Select [1-4]: ")
		var input string
		fmt.Scanln(&input)
		switch input {
		case "1", "2", "3":
			n, _ := strconv.Atoi(input)
			return Shell(n), ActionSelect
		case "4":
			return 0, ActionExit
		default:
			fmt.Println("\n  Invalid option. Please enter a number between 1 and 4.")
		}
	}
}
func showOverrideWarning() bool {
	fmt.Print(overrideWarningTemplate)
	for {
		fmt.Print("  Select [1-2]: ")
		var input string
		fmt.Scanln(&input)
		switch strings.TrimSpace(input) {
		case "1":
			return true
		case "2":
			return false
		default:
			fmt.Println("\n  Invalid option. Please enter 1 or 2.")
		}
	}
}

func shellNextDialog() bool {
	fmt.Print(shellNextDialogTemplate)
	for {
		fmt.Print("  Select [1-2]: ")
		var input string
		fmt.Scanln(&input)
		switch strings.TrimSpace(input) {
		case "1":
			return true
		case "2":
			return false
		default:
			fmt.Println("\n  Invalid option. Please enter 1 or 2.")
		}
	}
}
