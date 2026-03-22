package main

import (
	"flag"
	"fmt"
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
	file, err := openFile("file.txt")
	if err != nil {
		fmt.Println("There was an error trying to open the file: ", err)
		return
	}
	defer file.Close()
	readFile(file)
	writeFile(file)
	file.Seek(0, 0)
	readFile(file)
}
