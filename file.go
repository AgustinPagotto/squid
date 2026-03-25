package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

//func createAndWrite(note string) {
//	// Create or trunkate file
//	file, err := os.Create("file.md")
//	if err != nil {
//		fmt.Print("Couldn't create the initial file", err)
//	}
//	// Create buffer to read/write the file
//	var a = make([]byte, 10)
//	var b []byte
//	b = []byte(note)
//	// Write the phrase to the file
//	file.Write(b)
//	// Moves the cursor to the start of the file
//	_, _ = file.Seek(0, 0)
//	// Loop that reads the file 10 bytes by 10 bytes until EOF
//	for {
//		j, err := file.Read(a)
//		if errors.Is(err, io.EOF) && j == 0 {
//			break
//		}
//		fmt.Print(string(a[:j]))
//	}
//}

func readFile(file *os.File) {
	// Create a new reader
	reader := bufio.NewReader(file)
	counter := 0
	for {
		// Read the string until we met with a \n
		line, err := reader.ReadString('\n')
		counter += 1
		fmt.Println(counter)
		fmt.Print(line)
		if err != nil {
			break
		}
	}
}

func writeFile(file *os.File, typeOfNote NoteType, note string) {
	writer := bufio.NewWriter(file)
	defer writer.Flush()
	fmt.Fprintf(writer, "\n %v \n----------", note)
}

func openFile(fileName string) (*os.File, error) {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func countNotes(file *os.File) {
	reader := bufio.NewReader(file)
	counter := 0
	var notes map[int]string = make(map[int]string)
	for {
		line, err := reader.ReadString('\n')
		if !strings.Contains(line, "----------") {
			if notes[counter] != "" {
				notes[counter] = fmt.Sprintf("%s\n%s", notes[counter], line)
			} else {
				notes[counter] = line
			}
		} else {
			counter += 1
		}
		if err != nil {
			break
		}
	}
	fmt.Println(notes)
}
