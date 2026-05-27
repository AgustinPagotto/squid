package setup

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/AgustinPagotto/squid/internal/config"
	"github.com/AgustinPagotto/squid/internal/context"
)

func Init(args []string) error {
	var firstArg string
	if len(args) > 0 {
		firstArg = args[0]
	}
	switch firstArg {
	case "-h":
		fmt.Print(commandHelp)
	case "shell":
		if shell, exited := handleShell(); !exited {
			fmt.Print(buildSuccessTemplate(shell))
		}
	case "":
		dir, err := config.FindRoot()
		if err != nil && err != config.ErrRootNotFound {
			return err
		}
		var overrideConfirmed bool
		if err == nil && dir != "" {
			overrideConfirmed = showOverrideWarning()
			if !overrideConfirmed {
				return nil
			}
		}
		selectedOption, action := giveInitialOptions()
		if action == ActionExit {
			return nil
		}
		switch selectedOption {
		case FlowPredefined:
			option, action := predefinedInit()
			if action == ActionExit {
				fmt.Println("  squid init exited")
				return nil
			}
			if overrideConfirmed {
				err := os.RemoveAll(".squid")
				if err != nil {
					return err
				}
			}
			err := os.Mkdir(".squid", config.DirPerm)
			if err != nil {
				return err
			}
			if err := initiatePredefinedSquidFolder(option); err != nil {
				return err
			}
			var shell Shell
			if shellNextDialog() {
				shell, _ = handleShell()
			}
			fmt.Print(buildSuccessTemplate(shell))
			return nil
		case FlowCustom:
			if overrideConfirmed {
				err := os.RemoveAll(".squid")
				if err != nil {
					return err
				}
			}
			err := os.Mkdir(".squid", config.DirPerm)
			if err != nil {
				return err
			}
			if err := handleCustom(); err != nil {
				os.RemoveAll(".squid")
				fmt.Println(err)
				return nil
			}
			var shell Shell
			if shellNextDialog() {
				shell, _ = handleShell()
			}
			fmt.Print(buildSuccessTemplate(shell))
			return nil
		default:
			return nil
		}
	default:
		fmt.Printf(unknownSubcommandTemplate, firstArg)
	}
	return nil
}

func buildSuccessTemplate(shell Shell) string {
	if shell == 0 {
		return successTemplate
	}
	return fmt.Sprintf(successWithShellTemplate, shell.sourceCmd())
}

func handleShell() (shell Shell, exited bool) {
	selectedShell, action := selectShell()
	if action == ActionExit {
		return 0, true
	}
	err := addEvalPermissions(selectedShell)
	if err != nil {
		fmt.Println(err)
		return 0, true
	}
	return selectedShell, false
}

func handleCustom() error {
	content, err := config.OpenInEditor(addAliasesTemplate)
	if err != nil {
		return err
	}
	parsedAliases, err := context.ParseAliases(content)
	if err != nil {
		return err
	}
	cs := &context.ContextStorage{Path: ".squid/context.json"}
	err = context.HandleContextCreation(parsedAliases, cs)
	if err != nil {
		return err
	}
	return nil
}

func giveInitialOptions() (Flow, Action) {
	fmt.Print(initialOptionTemplate)
	for {
		fmt.Print("  Select [1-3]: ")
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			return 0, ActionExit
		}
		switch input {
		case "1", "2":
			n, _ := strconv.Atoi(input)
			return Flow(n), ActionSelect
		case "3":
			return 0, ActionExit
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
		_, err := fmt.Scanln(&input)
		if err != nil {
			return 0, ActionExit
		}
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
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	path, err := config.FindShellConfigurationFile(selectedShell.configFile(), home)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, config.FilePerm)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("failed to close config file: %v", err)
		}
	}()

	existing, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	if strings.Contains(string(existing), selectedShell.evalString()) {
		return fmt.Errorf("Shell hook already present, no need to add it again")
	}

	_, err = f.WriteString(selectedShell.evalString())
	return err
}

func selectShell() (Shell, Action) {
	fmt.Print(selectShellTemplate)
	for {
		fmt.Print("  Select [1-4]: ")
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			return 0, ActionExit
		}
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
		_, err := fmt.Scanln(&input)
		if err != nil {
			return false
		}
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
		_, err := fmt.Scanln(&input)
		if err != nil {
			return false
		}
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
