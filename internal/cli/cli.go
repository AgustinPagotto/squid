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
