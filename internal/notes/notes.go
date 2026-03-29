package notes

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/AgustinPagotto/squid/internal/cli"
	"github.com/AgustinPagotto/squid/internal/validator"
)

//type NoteType int

//const (
//	NoteTypeNote = iota
//	NoteTypeChecklist
//)

type Note struct {
	ID        int
	Title     string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func Handle(args []string) {
	if len(args) == 0 {
		fmt.Println("expected subcommand (add, list, etc.)")
		return
	}
	switch args[0] {
	case "add", "a":
		err := handleAdd()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("note added successfully")
	case "edit", "e":
		id, err := parseID(args, 1)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = handleEdit(id)
		if err != nil {
			fmt.Println(err)
		}
	case "del", "d":
		id, err := parseID(args, 1)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = handleDelete(id)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("note deleted")
	case "list", "l":
		err := handleList()
		if err != nil {
			fmt.Println(err)
			return
		}
	case "show", "s":
		id, err := parseID(args, 1)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = handleShow(id)
		if err != nil {
			fmt.Println(err)
			return
		}
	default:
		fmt.Println("unknown subcommand:", args[0])
	}
}

func handleAdd() error {
	content, err := openInEditor(cli.AddNoteTemplate)
	if err != nil {
		return err
	}
	title, body, err := parseNote(content)
	if err != nil {
		return err
	}
	note := Note{Title: title, Body: body, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	return addNote(note)
}

func handleEdit(id int) error {
	note, err := findNote(id)
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
	return editNote(*note)
}

func handleDelete(id int) error {
	note, err := findNote(id)
	if err != nil {
		return err
	}
	fmt.Printf("Are you sure you want to delete note: %q? [Y/n]: ", note.Title)
	var input string
	fmt.Scanln(&input)
	input = strings.TrimSpace(strings.ToLower(input))
	if input != "" && input != "y" && input != "yes" {
		fmt.Println("delete aborted")
		return nil
	}
	if err := delNote(id); err != nil {
		return err
	}
	return nil
}

func handleList() error {
	notes, err := loadNotes()
	if err != nil {
		return err
	}
	if len(notes) == 0 {
		return fmt.Errorf("No notes found")
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

func handleShow(id int) error {
	note, err := findNote(id)
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
	return strings.TrimSpace(result[0]), strings.TrimSpace(strings.Join(result[1:], "\n")), nil
}

func parseID(args []string, loc int) (int, error) {
	if !validator.HasArgAmount(args, loc+1) {
		fmt.Println("please provide a note id")
		return 0, fmt.Errorf("invalid amount of arguments")
	}
	id, err := strconv.Atoi(args[loc])
	if err != nil {
		return 0, fmt.Errorf("invalid note id")
	}
	return id, nil
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
