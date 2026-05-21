package todos

import (
	"strings"
	"testing"
	"time"
)

func TestHandleList(t *testing.T) {
	tests := []struct {
		name    string
		todos   []Todo
		wantErr bool
	}{
		{
			name:    "empty store prints message",
			todos:   []Todo{},
			wantErr: false,
		},
		{
			name:    "pending todos only",
			todos:   []Todo{{Title: "Pending A"}, {Title: "Pending B"}},
			wantErr: false,
		},
		{
			name:    "done todos only",
			todos:   []Todo{{Title: "Done A", IsDone: true}},
			wantErr: false,
		},
		{
			name: "mixed pending and done",
			todos: []Todo{
				{Title: "Pending", IsDone: false},
				{Title: "Done", IsDone: true},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := newTestStoreWithTodos(t, tt.todos...)
			err := handleList(ts)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleShow(t *testing.T) {
	ts := newTestStoreWithTodos(t,
		Todo{Title: "Pending task", IsDone: false, CreatedAt: time.Now()},
		Todo{Title: "Done task", IsDone: true, CreatedAt: time.Now()},
	)
	todos, _ := ts.loadTodos()
	pendingID := todos[0].ID
	doneID := todos[1].ID

	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			name:    "existing pending todo",
			id:      pendingID,
			wantErr: false,
		},
		{
			name:    "existing done todo",
			id:      doneID,
			wantErr: false,
		},
		{
			name:    "missing todo returns error",
			id:      99,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleShow(ts, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleShow(%d) error = %v, wantErr %v", tt.id, err, tt.wantErr)
			}
		})
	}
}

func TestHandleAdd(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "valid text",
			args:    []string{"add", "Buy oat milk"},
			wantErr: false,
		},
		{
			name:    "empty text",
			args:    []string{"add", ""},
			wantErr: true,
		},
		{
			name:    "whitespace only",
			args:    []string{"add", "   "},
			wantErr: true,
		},
		{
			name:    "text exceeds 280 chars",
			args:    []string{"add", strings.Repeat("x", 281)},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := newTestStore(t)
			err := handleAdd(tt.args, ts)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleAdd() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				todos, _ := ts.loadTodos()
				if len(todos) != 1 {
					t.Errorf("expected 1 todo after add, got %d", len(todos))
				}
			}
		})
	}
}

func TestHandleDeleteMissingTodo(t *testing.T) {
	ts := newTestStore(t)
	err := handleDelete(ts, 99)
	if err == nil {
		t.Error("handleDelete() expected error for non-existing id, got nil")
	}
}

func TestHandle(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "no args",
			args: []string{},
		},
		{
			name: "help flag",
			args: []string{"-h"},
		},
		{
			name: "add missing text arg",
			args: []string{"add"},
		},
		{
			name: "del with non-numeric id",
			args: []string{"del", "abc"},
		},
		{
			name: "show with non-numeric id",
			args: []string{"show", "abc"},
		},
		{
			name: "unknown subcommand",
			args: []string{"unknown"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := newTestStore(t)
			// Handle prints to stdout; we just verify it does not panic.
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Handle(%v) panicked: %v", tt.args, r)
				}
			}()
			Handle(tt.args, ts)
		})
	}
}
