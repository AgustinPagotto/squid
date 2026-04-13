package context

import "fmt"

func printHelp() {
	fmt.Print(`Manage project context and aliases

Usage:
  squid context <subcommand>

Subcommands:
  activate, a    Load context aliases into the current shell session

Examples:
  eval "$(squid context activate)"
`)
}
