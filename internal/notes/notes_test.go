package notes

import (
	"testing"
	"time"
)

func newTestStorageWithNotes(t *testing.T, notes ...Note) *NoteStorage {
	t.Helper()
	newNoteStorage := &NoteStorage{Path: t.TempDir() + "/notes.json"}
	for _, note := range notes {
		newNoteStorage.addNote(note)
	}
	return newNoteStorage
}

func TestHandleList(t *testing.T) {
	tests := []struct {
		name    string
		notes   []Note
		wantErr bool
	}{
		{"lists notes", []Note{{ID: 1, Title: "A", Body: "B"}}, false},
		{"empty storage returns error", []Note{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleList(newTestStorageWithNotes(t, tt.notes...))
			if (err != nil) != tt.wantErr {
				t.Errorf("handleList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleShow(t *testing.T) {
	ns := newTestStorageWithNotes(t, Note{Title: "Hello", Body: "World", CreatedAt: time.Now(), UpdatedAt: time.Now()})

	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{"existing note", 0, false},
		{"missing note", 99, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleShow(ns, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleShow(%d) error = %v, wantErr %v", tt.id, err, tt.wantErr)
			}
		})
	}
}

func TestHandleDelete(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		wantErr    bool
		wantRemain int
	}{
		{"deletes existing note", 0, false, 1},
		{"missing note returns error", 99, true, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ns := newTestStorageWithNotes(t,
				Note{Title: "To delete", Body: "body", CreatedAt: time.Now(), UpdatedAt: time.Now()},
				Note{Title: "Not deleted", Body: "body", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			)
			err := ns.delNote(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("delNote(%d) error = %v, wantErr %v", tt.id, err, tt.wantErr)
			}
			notes, _ := ns.loadNotes()
			if len(notes) != tt.wantRemain {
				t.Errorf("remaining notes = %d, want %d", len(notes), tt.wantRemain)
			}
		})
	}
}

func TestParseNote(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantTitle string
		wantBody  string
		wantErr   bool
	}{
		{
			name:      "valid note",
			input:     "My Title\nThis is the body",
			wantTitle: "My Title",
			wantBody:  "This is the body",
		},
		{
			name:      "strips comment lines",
			input:     "# ignore this\nMy Title\nBody text",
			wantTitle: "My Title",
			wantBody:  "Body text",
		},
		{
			name:      "leading blank lines are skipped",
			input:     "\n\nMy Title\nBody text",
			wantTitle: "My Title",
			wantBody:  "Body text",
		},
		{
			name:    "only title, no body",
			input:   "Just a title",
			wantErr: true,
		},
		{
			name:    "empty input",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			title, body, err := parseNote(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseNote() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if title != tt.wantTitle {
					t.Errorf("title = %q, want %q", title, tt.wantTitle)
				}
				if body != tt.wantBody {
					t.Errorf("body = %q, want %q", body, tt.wantBody)
				}
			}
		})
	}
}

func TestParseID(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		loc     int
		want    int
		wantErr bool
	}{
		{"valid id at loc 1", []string{"edit", "42"}, 1, 42, false},
		{"valid id at loc 0", []string{"99"}, 0, 99, false},
		{"not a number", []string{"edit", "abc"}, 1, 0, true},
		{"args too short", []string{"edit"}, 1, 0, true},
		{"no args", []string{}, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseID(tt.args, tt.loc)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseID(%v, %d) error = %v, wantErr %v", tt.args, tt.loc, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("parseID(%v, %d) = %d, want %d", tt.args, tt.loc, got, tt.want)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name  string
		input string
		max   int
		want  string
	}{
		{"shorter than max", "hello", 10, "hello"},
		{"exact length", "hello", 5, "hello"},
		{"truncated with ellipsis", "hello world", 8, "hello..."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := truncate(tt.input, tt.max); got != tt.want {
				t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.max, got, tt.want)
			}
		})
	}
}
