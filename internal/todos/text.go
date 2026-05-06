package todos

const commandHelp = `Manage project todos

Usage:
  squid todo <subcommand> [arguments]

Subcommands:
  add,    a  <text>   Create a new todo
  list,   l           List all todos
  show,   s  <id>     Show details of a todo
  toggle, t           Toggle a todo (prompts if no id given)
  del,    d  <id>     Delete a todo
  clear               Remove all completed todos

Examples:
  squid todo add "fix the auth bug"
  squid todo list
  squid todo toggle
  squid todo del 2
  squid todo clear
`
