package todos

const commandHelp = `Manage project todoss

Usage:
  squid todos <subcommand> [arguments]

Subcommands:
  add,    a  <text>   Create a new todos
  list,   l           List all todoss
  show,   s  <id>     Show details of a todos
  toggle, t           Toggle a todos (prompts if no id given)
  del,    d  <id>     Delete a todos
  clear               Remove all completed todos

Examples:
  squid todos add "fix the auth bug"
  squid todos list
  squid todos toggle
  squid todos del 2
  squid todos clear
`
