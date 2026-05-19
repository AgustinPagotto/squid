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
		{
			name:    "lists notes",
			notes:   []Note{{ID: 1, Title: "A", Body: "B"}},
			wantErr: false,
		},
		{
			name:    "empty storage prints message",
			notes:   []Note{},
			wantErr: false,
		},
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
		{
			name:    "existing note",
			id:      0,
			wantErr: false,
		},
		{
			name:    "missing note",
			id:      99,
			wantErr: true,
		},
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
		{
			name:       "deletes existing note",
			id:         0,
			wantErr:    false,
			wantRemain: 1,
		},
		{
			name:       "missing note returns error",
			id:         99,
			wantErr:    true,
			wantRemain: 2,
		},
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

func TestTruncate(t *testing.T) {
	tests := []struct {
		name  string
		input string
		max   int
		want  string
	}{
		{
			name:  "shorter than max",
			input: "hello",
			max:   10,
			want:  "hello",
		},
		{
			name:  "exact length",
			input: "hello",
			max:   5,
			want:  "hello",
		},
		{
			name:  "truncated with ellipsis",
			input: "hello world",
			max:   8,
			want:  "hello...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := truncate(tt.input, tt.max); got != tt.want {
				t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.max, got, tt.want)
			}
		})
	}
}
