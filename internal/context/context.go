package context

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/AgustinPagotto/squid/internal/config"
	"github.com/AgustinPagotto/squid/internal/validator"
	"golang.org/x/term"
)

type Alias struct {
	ID           int
	AliasCommand string
}

// AliasRegex matches name="value" or name='value' with a non-empty value.
const AliasRegex = `^[a-zA-Z_][a-zA-Z0-9_-]*=("[^"]+"|'[^']+')$`

const (
	MaxAliasNameLen    = 50
	MaxAliasCommandLen = 250
)

func Handle(args []string, cs ContextStorageInterface) {
	if len(args) == 0 {
		fmt.Println("expected subcommand (add, list, etc.)")
		return
	}
	switch args[0] {
	case "-h":
		fmt.Print(commandHelp)
	case "activate", "ac":
		err := handleActivate(cs)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		fmt.Fprintln(os.Stdout, "activation of the context succeded")
	case "add", "a", "edit":
		err := handleAdd(cs)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "list", "l":
		err := handleList(cs)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "del", "d":
		id, err := validator.ValidateAndParseID(args, 1)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = handleDelete(cs, id)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("alias deleted")
	default:
		fmt.Println("unknown subcommand:", args[0])
	}
}

func HandleContextCreation(entries []string, cs ContextStorageInterface) error {
	var aliases []Alias
	for i, aliasString := range entries {
		aliases = append(aliases, Alias{ID: i, AliasCommand: aliasString})
	}
	err := cs.persistContext(aliases)
	if err != nil {
		return err
	}
	return nil
}

func handleDelete(cs ContextStorageInterface, id int) error {
	alias, err := cs.findAlias(id)
	if err != nil {
		return err
	}
	fmt.Printf("Are you sure you want to delete alias: %q? [Y/n]: ", alias.AliasCommand)
	var input string
	fmt.Scanln(&input)
	input = strings.TrimSpace(strings.ToLower(input))
	if input != "" && input != "y" && input != "yes" {
		return fmt.Errorf("delete aborted")
	}
	return cs.delAlias(id)
}

func handleActivate(cs ContextStorageInterface) error {
	if term.IsTerminal(int(os.Stdout.Fd())) {
		return fmt.Errorf("hint: run as: eval \"$(squid context activate)\"")
	}
	aliases, err := cs.loadAliases()
	if err != nil {
		return err
	}
	for _, alias := range aliases {
		fmt.Printf("alias %s\n", alias.AliasCommand)
	}
	return nil
}

func handleList(cs ContextStorageInterface) error {
	aliases, err := cs.loadAliases()
	if err != nil {
		return err
	}
	if len(aliases) == 0 {
		fmt.Println("no aliases defined yet — run 'squid context add' to create one")
		return nil
	}

	for _, a := range aliases {
		fmt.Printf("%2d  %s\n", a.ID, a.AliasCommand)
	}

	return nil
}

func handleAdd(cs ContextStorageInterface) error {
	existing, err := cs.loadAliases()
	if err != nil {
		return err
	}

	var sb strings.Builder
	sb.WriteString(addAliasesHeader)
	for _, a := range existing {
		sb.WriteString(a.AliasCommand + "\n")
	}

	content, err := config.OpenInEditor(sb.String())
	if err != nil {
		return err
	}
	parsedAliases, err := ParseAliases(content)
	if err != nil {
		return err
	}
	return cs.persistContext(assignIDs(existing, parsedAliases))
}

func assignIDs(existing []Alias, parsed []string) []Alias {
	existingIDs := make(map[string]int, len(existing))
	for _, a := range existing {
		existingIDs[a.AliasCommand] = a.ID
	}
	var aliases []Alias
	for _, aliasString := range parsed {
		if id, ok := existingIDs[aliasString]; ok {
			aliases = append(aliases, Alias{ID: id, AliasCommand: aliasString})
		} else {
			aliases = append(aliases, Alias{ID: nextID(aliases), AliasCommand: aliasString})
		}
	}
	return aliases
}

func ParseAliases(input string) ([]string, error) {
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
		name := strings.SplitN(line, "=", 2)[0]
		if len(name) > MaxAliasNameLen {
			continue
		}
		// value is everything after name= and inside the surrounding quotes
		value := line[len(name)+2 : len(line)-1]
		if len(value) > MaxAliasCommandLen {
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
