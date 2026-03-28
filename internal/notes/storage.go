package notes

import (
	"encoding/json"
	"fmt"
	"os"
)

const notesFile = "notes.json"

func findNote(id int) (*Note, error) {
	notes, err := load()
	if err != nil {
		return &Note{}, err
	}
	for i := range notes {
		if notes[i].ID == id {
			return &notes[i], nil
		}
	}
	return &Note{}, fmt.Errorf("note with id %d not found", id)
}
func saveEdit(note Note) error {
	notes, err := load()
	if err != nil {
		return err
	}
	for i := range notes {
		if notes[i].ID == note.ID {
			notes[i] = note
		}
	}
	jsonNotes, err := json.MarshalIndent(notes, "", " ")
	if err != nil {
		return err
	}
	fmt.Println(string(jsonNotes))
	err = os.WriteFile(notesFile, jsonNotes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func save(note Note) error {
	notes, err := load()
	if err != nil {
		return err
	}
	note.ID = nextID(notes)
	notes = append(notes, note)
	jsonNotes, err := json.MarshalIndent(notes, "", " ")
	if err != nil {
		return err
	}
	fmt.Println(string(jsonNotes))
	err = os.WriteFile(notesFile, jsonNotes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func load() ([]Note, error) {
	file, err := os.ReadFile(notesFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []Note{}, nil
		}
		return nil, err
	}
	var notes []Note
	err = json.Unmarshal(file, &notes)
	return notes, err
}

func nextID(notes []Note) int {
	maxID := 0
	for _, note := range notes {
		if note.ID > maxID {
			maxID = note.ID
		}
	}
	return maxID + 1
}

func openFile(fileName string) (*os.File, error) {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return file, nil
}
