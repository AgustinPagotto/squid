package notes

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/AgustinPagotto/squid/internal/cli"
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
	default:
		fmt.Println("unknown subcommand:", args[0])
	}
}

func handleAdd() error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano"
	}
	tmpFile, err := os.CreateTemp("", "squid-note-*.txt")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(cli.NoteTemplate)
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
	err = add(title, body)
	if err != nil {
		return err
	}
	return nil
}

func add(title, body string) error {
	note := Note{Title: title, Body: body, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	return save(note)
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
