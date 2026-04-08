package notes

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/AgustinPagotto/squid/internal/cli"
	"github.com/AgustinPagotto/squid/internal/validator"
)

type Note struct {
	ID        int
	Title     string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func Handle(args []string, ns NoteStorageInterface) {
	if len(args) == 0 {
		fmt.Println("expected subcommand (add, list, etc.)")
		return
	}
	switch args[0] {
	case "-h":
		cli.PrintNotesHelp()
	case "add", "a":
		err := handleAdd(ns)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("note added successfully")
	case "edit", "e":
		id, err := validator.ValidateAndParseID(args, 1)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = handleEdit(ns, id)
		if err != nil {
			fmt.Println(err)
		}
	case "del", "d":
		id, err := validator.ValidateAndParseID(args, 1)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = handleDelete(ns, id)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("note deleted")
	case "list", "l":
		err := handleList(ns)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "show", "s":
		id, err := validator.ValidateAndParseID(args, 1)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = handleShow(ns, id)
		if err != nil {
			fmt.Println(err)
			return
		}
	default:
		fmt.Println("unknown subcommand:", args[0])
	}
}

func handleAdd(ns NoteStorageInterface) error {
	content, err := openInEditor(cli.AddNoteTemplate)
	if err != nil {
		return err
	}
	title, body, err := parseNote(content)
	if err != nil {
		return err
	}
	note := Note{Title: title, Body: body, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	return ns.addNote(note)
}

func handleEdit(ns NoteStorageInterface, id int) error {
	note, err := ns.findNote(id)
	if err != nil {
		return err
	}
	initial := fmt.Sprintf("%s\n\n%s\n\n%s", note.Title, note.Body, cli.EditNoteTemplate)
	content, err := openInEditor(initial)
	if err != nil {
		return err
	}
	title, body, err := parseNote(content)
	if err != nil {
		return err
	}
	note.Title = title
	note.Body = body
	note.UpdatedAt = time.Now()
	return ns.editNote(*note)
}

func handleDelete(ns NoteStorageInterface, id int) error {
	note, err := ns.findNote(id)
	if err != nil {
		return err
	}
	fmt.Printf("Are you sure you want to delete note: %q? [Y/n]: ", note.Title)
	var input string
	fmt.Scanln(&input)
	input = strings.TrimSpace(strings.ToLower(input))
	if input != "" && input != "y" && input != "yes" {
		return fmt.Errorf("delete aborted")
	}
	return ns.delNote(id)
}

func handleList(ns NoteStorageInterface) error {
	notes, err := ns.loadNotes()
	if err != nil {
		return err
	}
	if len(notes) == 0 {
		fmt.Println("no notes yet — run 'squid notes add' to create one")
		return nil
	}

	for _, n := range notes {
		title := truncate(n.Title, 50)
		preview := truncate(strings.ReplaceAll(n.Body, "\n", " "), 80)
		fmt.Printf("%2d  %s\n", n.ID, title)
		fmt.Printf("    %s\n\n", preview)
	}

	fmt.Println()
	return nil
}

func handleShow(ns NoteStorageInterface, id int) error {
	note, err := ns.findNote(id)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n\n", note.Title)
	fmt.Println(note.Body)
	fmt.Println()
	fmt.Printf("Created: %s\n", note.CreatedAt.Format("2006-01-02"))
	fmt.Printf("Updated: %s\n", note.UpdatedAt.Format("2006-01-02"))

	return nil
}

func openInEditor(initial string) (string, error) {
	tmpFile, err := os.CreateTemp("", "squid-note-*.txt")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())

	if _, err = tmpFile.WriteString(initial); err != nil {
		return "", err
	}
	tmpFile.Close()

	cmd := exec.Command(getEditor(), tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return "", err
	}

	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func getEditor() string {
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	return "nano"
}

func parseNote(input string) (string, string, error) {
	lines := strings.Split(input, "\n")
	var result []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" && len(result) == 0 {
			continue
		}
		if line == "" {
			result = append(result, "")
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		result = append(result, line)
	}
	if len(result) < 2 {
		return "", "", fmt.Errorf("no sufficient content was written")
	}
	for len(result) > 0 && result[len(result)-1] == "" {
		result = result[:len(result)-1]
	}
	return strings.TrimSpace(result[0]), strings.TrimSpace(strings.Join(result[1:], "\n")), nil
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
