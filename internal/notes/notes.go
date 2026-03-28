package notes

import (
	"fmt"
	"os"
	"os/exec"
	"time"
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
		handleAdd()
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
	cmd := exec.Command(editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("error opening editor:", err)
		return err
	}
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		fmt.Println("error reading file:", err)
		return err
	}
	if len(content) == 0 {
		fmt.Println("no content was written: ")
		return err
	}
	fmt.Println(string(content))
	return nil
}

func Add(content string) error {
	note := Note{Title: "example", Body: content}
	return save(note)
}
