package notes

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"

	"github.com/AgustinPagotto/squid/internal/config"
)

type NoteStorageInterface interface {
	findNote(id int) (*Note, error)
	editNote(note Note) error
	addNote(note Note) error
	delNote(id int) error
	loadNotes() ([]Note, error)
}

type NoteStorage struct {
	Path string
}

func (ns *NoteStorage) findNote(id int) (*Note, error) {
	notes, err := ns.loadNotes()
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

func (ns *NoteStorage) editNote(note Note) error {
	notes, err := ns.loadNotes()
	if err != nil {
		return err
	}
	for i := range notes {
		if notes[i].ID == note.ID {
			notes[i] = note
			return ns.persistNotes(notes)
		}
	}
	return fmt.Errorf("note with id %d not found", note.ID)
}

func (ns *NoteStorage) addNote(note Note) error {
	notes, err := ns.loadNotes()
	if err != nil {
		return err
	}
	note.ID = nextID(notes)
	notes = append(notes, note)
	return ns.persistNotes(notes)
}

func (ns *NoteStorage) delNote(id int) error {
	notes, err := ns.loadNotes()
	if err != nil {
		return err
	}
	for i := range notes {
		if notes[i].ID == id {
			return ns.persistNotes(slices.Delete(notes, i, i+1))
		}
	}
	return fmt.Errorf("note with id %d not found", id)
}

func (ns *NoteStorage) persistNotes(notes []Note) error {
	jsonNotes, err := json.MarshalIndent(notes, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(ns.Path, jsonNotes, config.FilePerm)
}

func (ns *NoteStorage) loadNotes() ([]Note, error) {
	file, err := os.ReadFile(ns.Path)
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
