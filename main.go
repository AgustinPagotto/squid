package main

import (
	"flag"
	"fmt"
)

type NoteType int

const (
	NoteTypeNote = iota
	NoteTypeChecklist
)

func main() {
	var note string
	var noteType string
	var help bool
	var typeOfNote NoteType
	flag.StringVar(&note, "note", "", "Note to add to the project")
	flag.StringVar(&note, "n", "", "Note to add to the project")
	flag.StringVar(&noteType, "notetype", "", "Note type, you can do n for note or c for checklist")
	flag.StringVar(&noteType, "nt", "", "Note type, you can do n for note or c for checklist")
	flag.BoolVar(&help, "help", false, "Ask for help")
	flag.BoolVar(&help, "h", false, "Ask for help")
	flag.Parse()
	if help {
		fmt.Print("Help")
		return
	}
	if note == "" || noteType == "" {
		fmt.Print("you need to provide all the info")
		return
	}
	var fileName string
	switch noteType {
	case "n":
		typeOfNote = NoteTypeNote
		fileName = "notes.txt"
	case "c":
		typeOfNote = NoteTypeChecklist
		fileName = "checklist.txt"
	default:
		fmt.Print("You need to provide a compatible note title (c or n)")
		return
	}
	file, err := openFile(fileName)
	if err != nil {
		fmt.Println("There was an error trying to open the file: ", err)
		return
	}
	defer file.Close()
	writeFile(file, typeOfNote, note)
	file.Seek(0, 0)
	countNotes(file)
}
