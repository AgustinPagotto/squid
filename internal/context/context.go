package context

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

type Alias struct {
	ID           int
	AliasCommand string
}

func Handle(args []string, cs ContextStorageInterface) {
	if len(args) == 0 {
		fmt.Println("expected subcommand (add, list, etc.)")
		return
	}
	switch args[0] {
	case "-h":
		printHelp()
	case "activate", "a":
		err := handleActivate(cs)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		fmt.Fprintln(os.Stderr, "activation of the context succeded")
	case "list", "l":
		err := handleList(cs)
		if err != nil {
			fmt.Println(err)
			return
		}
	default:
		fmt.Println("unknown subcommand:", args[0])
	}
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
