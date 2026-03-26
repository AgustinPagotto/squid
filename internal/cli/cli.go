package cli

import "fmt"

func PrintHelp() {
	fmt.Print(`Squid — project-aware developer context CLI
Usage:
  squid <command> [arguments]

Available Commands:
  notes      Manage project notes and checklists
  context    Manage project context and settings

Flags:
  -h, Show help for squid

Examples:
  squid notes add "fix auth bug"
  squid notes list
  squid context set my-api

Use "squid <command> -h" for more information about a command.
`)
}
