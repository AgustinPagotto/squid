package notes

import "fmt"

const addNoteTemplate = `

# First line will be your note title
# The following lines will be for the actual note. Lines starting
# with '#' will be ignored, and an empty message aborts.
#
# Squid Notes
`

const editNoteTemplate = `

# Edit your note as needed, remember first line is the title
# Lines starting with '#' will be ignored, and an empty message aborts.
#
# Squid Notes
`

func printHelp() {
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
