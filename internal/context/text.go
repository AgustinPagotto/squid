package context

const commandHelp = `Manage project context and aliases

Usage:
  squid context <subcommand>

Subcommands:
  activate, ac   Load context aliases into the current shell session
  add, a         Add new aliases via your editor
  list, l        List all aliases

Examples:
  eval "$(squid context activate)"
  squid context add
  squid context list
`

const addAliasesHeader = `# Add new aliases below the existing ones, one per line.
# Format:  <name>="<command>"
#
# Examples:
#   start="npm run dev"
#   test="go test ./..."
#
# Lines starting with '#' are ignored.
# Save and close the editor when done.
#
`
