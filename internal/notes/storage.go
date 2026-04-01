package notes

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
)

var notesFile = "notes.json"

func findNote(id int) (*Note, error) {
	notes, err := loadNotes()
	if err != nil {
		return nil, err
	}
	for i := range notes {
		if notes[i].ID == id {
			return &notes[i], nil
		}
	}
	return nil, fmt.Errorf("note with id %d not found", id)
}

func editNote(note Note) error {
	notes, err := loadNotes()
	if err != nil {
		return err
	}
	for i := range notes {
		if notes[i].ID == note.ID {
			notes[i] = note
			break
		}
	}
	return persistNotes(notes)
}

func addNote(note Note) error {
	notes, err := loadNotes()
	if err != nil {
		return err
	}
	note.ID = nextID(notes)
	notes = append(notes, note)
	return persistNotes(notes)
}

func delNote(id int) error {
	notes, err := loadNotes()
	if err != nil {
		return err
	}
	for i := range notes {
		if notes[i].ID == id {
			notes = slices.Delete(notes, i, i+1)
			break
		}
	}
	return persistNotes(notes)
}

func persistNotes(notes []Note) error {
	jsonNotes, err := json.MarshalIndent(notes, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(notesFile, jsonNotes, 0644)
}

func loadNotes() ([]Note, error) {
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
	if len(notes) == 0 {
		return 0
	}
	maxID := 0
	for _, note := range notes {
		if note.ID > maxID {
			maxID = note.ID
		}
	}
	return maxID + 1
}
