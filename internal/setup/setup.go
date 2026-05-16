package setup

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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
				os.RemoveAll(".squid")
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
				os.RemoveAll(".squid")
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
	content, err := openInEditor(addAliasesTemplate)
	if err != nil {
		return err
	}
	parsedAliases, err := parseAliases(content)
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

func parseAliases(input string) ([]string, error) {
	lines := strings.Split(input, "\n")
	r := regexp.MustCompile(AliasRegex)
	var result []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			continue
		}
		if !r.MatchString(line) {
			continue
		}
		result = append(result, line)
	}
	for len(result) > 0 && result[len(result)-1] == "" {
		result = result[:len(result)-1]
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no valid aliases found — aborting")
	}
	return result, nil
}

func openInEditor(initial string) (string, error) {
	tmpFile, err := os.CreateTemp("", "squid-context-*.txt")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())

	if _, err = tmpFile.WriteString(initial); err != nil {
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

func giveInitialOptions() (Flow, Action) {
	fmt.Print(initialOptionTemplate)
	for {
		fmt.Print("  Select [1-3]: ")
		var input string
		fmt.Scanln(&input)
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
		return fmt.Errorf("Shell hook already present, no need to add it again.")
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
