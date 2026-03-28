package notes

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/AgustinPagotto/squid/internal/cli"
	validator "github.com/AgustinPagotto/squid/internal/validator"
)

type NoteType int

const (
	NoteTypeNote = iota
	NoteTypeChecklist
)

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
	case "add":
		err := handleAdd()
		if err != nil {
			fmt.Println(err)
		}
	case "edit":
		if !validator.HasArgAmount(args, 2) {
			fmt.Println("please provide a note id")
			return
		}
		id, err := parseID(args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		err = handleEdit(id)
		if err != nil {
			fmt.Println(err)
		}
	case "del":
		if !validator.HasArgAmount(args, 2) {
			fmt.Println("please provide a note id to delete")
			return
		}
		id, err := parseID(args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		err = handleDelete(id)
		if err != nil {
			fmt.Println(err)
		}
	default:
		fmt.Println("unknown subcommand:", args[0])
	}
}

func handleAdd() error {
	editor := getEditor()
	tmpFile, err := os.CreateTemp("", "squid-note-*.txt")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(cli.AddNoteTemplate)
	if err != nil {
		return err
	}
	tmpFile.Close()

	cmd := exec.Command(editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return err
	}
	title, body, err := parseNote(string(content))
	if err != nil {
		return err
	}
	fmt.Println("Title: ", title, "\nBody: ", body)
	note := Note{Title: title, Body: body, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	err = save(note)
	if err != nil {
		return err
	}
	return nil
}

func handleEdit(id int) error {
	note, err := findNote(id)
	if err != nil {
		return err
	}
	editor := getEditor()
	tmpFile, err := os.CreateTemp("", "squid-note-*.txt")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())
	writeContent := fmt.Sprintf("%s\n\n%s\n\n%s",
		note.Title,
		note.Body,
		cli.EditNoteTemplate,
	)

	_, err = tmpFile.WriteString(writeContent)
	if err != nil {
		return err
	}
	tmpFile.Close()

	cmd := exec.Command(editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return err
	}
	title, body, err := parseNote(string(content))
	if err != nil {
		return err
	}
	note.Title = title
	note.Body = body
	err = saveEdit(*note)
	if err != nil {
		return err
	}
	return nil
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
	if err := delete(id); err != nil {
		return err
	}
	fmt.Println("note deleted")
	return nil
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

func parseID(s string) (int, error) {
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid note id")
	}
	return id, nil
}
