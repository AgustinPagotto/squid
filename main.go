package main

import (
	"fmt"
	"os"

	"github.com/AgustinPagotto/squid/internal/cli"
)

type NoteType int

const (
	NoteTypeNote = iota
	NoteTypeChecklist
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("please add a subcommand or enter -h flag to show info")
		return
	}
	switch os.Args[1] {
	case "-h":
		cli.PrintHelp()
		return
	case "notes":
		fmt.Println("notes")
		return
	case "context":
		fmt.Println("context")
		return
	default:
		fmt.Println("subcommand not valid")
		return
	}
	// fmt.Println(os.Args[1])
	//
	//	if help {
	//		fmt.Print("Help")
	//		return
	//	}
	//
	//	if note == "" || noteType == "" {
	//		fmt.Print("you need to provide all the info")
	//		return
	//	}
	//
	// var fileName string
	// switch noteType {
	// case "n":
	//
	//	typeOfNote = NoteTypeNote
	//	fileName = "notes.txt"
	//
	// case "c":
	//
	//	typeOfNote = NoteTypeChecklist
	//	fileName = "checklist.txt"
	//
	// default:
	//
	//		fmt.Print("You need to provide a compatible note title (c or n)")
	//		return
	//	}
	//
	// file, err := openFile(fileName)
	//
	//	if err != nil {
	//		fmt.Println("There was an error trying to open the file: ", err)
	//		return
	//	}
	//
	// defer file.Close()
	// writeFile(file, typeOfNote, note)
	// file.Seek(0, 0)
	// countNotes(file)
}
