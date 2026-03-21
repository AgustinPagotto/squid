package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	var note string
	var noteTitle string
	var help bool
	flag.StringVar(&note, "note", "", "Note to add to the project")
	flag.StringVar(&noteTitle, "notetitle", "", "Note title to add to the project")
	flag.StringVar(&note, "n", "", "Note to add to the project")
	flag.StringVar(&noteTitle, "nt", "", "Note title to add to the project")
	flag.BoolVar(&help, "help", false, "Ask for help")
	flag.BoolVar(&help, "h", false, "Ask for help")
	flag.Parse()
	if note == "" || noteTitle == "" {
		fmt.Print("you need to provide all the info")
	}
	if help {
		fmt.Print("Help")
	}
	// Create or trunkate file
	file, err := os.Create("file.md")
	if err != nil {
		fmt.Print("Couldn't create the initial file", err)
	}
	// Create buffer to read/write the file
	var a = make([]byte, 10)
	var b []byte
	b = []byte(note)
	// Write the phrase to the file
	file.Write(b)
	// Moves the cursor to the start of the file
	_, _ = file.Seek(0, 0)
	// Loop that reads the file 10 bytes by 10 bytes until EOF
	for {
		j, err := file.Read(a)
		if errors.Is(err, io.EOF) && j == 0 {
			break
		}
		fmt.Print(string(a[:j]))
	}
}
