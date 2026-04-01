package notes

import (
	"testing"
	"time"
)

func TestNextID(t *testing.T) {
	tests := []struct {
		name  string
		notes []Note
		want  int
	}{
		{"continues after max id", []Note{{ID: 1}, {ID: 3}, {ID: 2}}, 4},
		{"single note", []Note{{ID: 5}}, 6},
		{"no notes", []Note{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nextID(tt.notes); got != tt.want {
				t.Errorf("nextID() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	original := notesFile
	notesFile = dir + "/notes.json"
	defer func() { notesFile = original }()

	note := Note{Title: "Test", Body: "Body", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := addNote(note); err != nil {
		t.Fatalf("addNote() error = %v", err)
	}

	notes, err := loadNotes()
	if err != nil {
		t.Fatalf("load() error = %v", err)
	}
	if len(notes) != 1 {
		t.Fatalf("expected 1 note, got %d", len(notes))
	}
	if notes[0].Title != note.Title || notes[0].Body != note.Body {
		t.Errorf("loaded note = %+v, want title=%q body=%q", notes[0], note.Title, note.Body)
	}
}

func TestLoadMissingFile(t *testing.T) {
	original := notesFile
	notesFile = t.TempDir() + "/nonexistent.json"
	defer func() { notesFile = original }()

	notes, err := loadNotes()
	if err != nil {
		t.Fatalf("load() on missing file should return empty slice, got error: %v", err)
	}
	if len(notes) != 0 {
		t.Errorf("expected empty slice, got %d notes", len(notes))
	}
}

func TestFindNote(t *testing.T) {
	dir := t.TempDir()
	original := notesFile
	notesFile = dir + "/notes.json"
	defer func() { notesFile = original }()

	n := Note{Title: "Find me", Body: "body", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := addNote(n); err != nil {
		t.Fatalf("addNote() error = %v", err)
	}

	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{"existing note", 0, false},
		{"non-existing note", 99, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findNote(tt.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("findNote(%d) error = %v, wantErr %v", tt.id, err, tt.wantErr)
			}
			if !tt.wantErr && got == nil {
				t.Error("findNote() returned nil without error")
			}
		})
	}
}
