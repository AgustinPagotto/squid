package cli

import "fmt"

const AddNoteTemplate = `

# First line will be your note title
# The following lines will be for the actual note. Lines starting
# with '#' will be ignored, and an empty message aborts.
#
# Squid Notes
`

const EditNoteTemplate = `

# Edit your note as needed, remember first line is the title
# Lines starting with '#' will be ignored, and an empty message aborts.
#
# Squid Notes
`

func PrintHelp() {
	fmt.Print(`Squid — project-aware developer context CLI
Usage:
  squid <command> [arguments]

Available Commands:
  init       Initialize squid in the current directory
  notes      Manage project notes

Flags:
  -h         Show help for squid

Use "squid <command> -h" for more information about a command.
`)
}

func PrintInitHelp() {
	fmt.Print(`Initialize squid in the current directory

Usage:
  squid init

Creates a .squid/ directory in the current directory. Run this once
per project. All squid data (notes, todos, context) will be stored there.
`)
}

func PrintNotesHelp() {
	fmt.Print(`Manage project notes

Usage:
  squid notes <subcommand> [arguments]

Subcommands:
  add,  a           Create a new note in your editor
  list, l           List all notes
  show, s <id>      Show full content of a note
  edit, e <id>      Edit an existing note in your editor
  del,  d <id>      Delete a note

Examples:
  squid notes add
  squid notes list
  squid notes show 1
  squid notes edit 1
  squid notes del 1
`)
}

func PrintTodoHelp() {
	fmt.Print(`Manage project todos

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
`)
}
