package notes

import (
	"bufio"
	"fmt"
	"os"
)

const notesFile = "notes.txt"

func save(note Note) error {
	file, err := openFile(notesFile)
	if err != nil {
		fmt.Println("error opening file", err)
		return err
	}
	writer := bufio.NewWriter(file)
	defer func() {
		writer.Flush()
		file.Close()
	}()
	fmt.Fprintf(writer, "\n %v \n\n", note.Title)
	fmt.Fprintf(writer, "\n %v \n----------", note.Body)
	return nil
}

func openFile(fileName string) (*os.File, error) {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return file, nil
}

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

//func readFile(file *os.File) {
//	// Create a new reader
//	reader := bufio.NewReader(file)
//	counter := 0
//	for {
//		// Read the string until we met with a \n
//		line, err := reader.ReadString('\n')
//		counter += 1
//		fmt.Println(counter)
//		fmt.Print(line)
//		if err != nil {
//			break
//		}
//	}
//}
