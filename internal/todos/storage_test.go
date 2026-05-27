package todos

import (
	"testing"
	"time"
)

func newTestStore(t *testing.T) *TodoStorage {
	t.Helper()
	return &TodoStorage{Path: t.TempDir() + "/todos.json"}
}

func newTestStoreWithTodos(t *testing.T, todos ...Todo) *TodoStorage {
	t.Helper()
	ts := newTestStore(t)
	for _, todo := range todos {
		if err := ts.addTodo(todo); err != nil {
			t.Fatalf("addTodo() error = %v", err)
		}
	}
	return ts
}

func TestNextID(t *testing.T) {
	tests := []struct {
		name  string
		todos []Todo
		want  int
	}{
		{
			name:  "no todos",
			todos: []Todo{},
			want:  0,
		},
		{
			name:  "single todo",
			todos: []Todo{{ID: 5}},
			want:  6,
		},
		{
			name:  "continues after max id",
			todos: []Todo{{ID: 1}, {ID: 3}, {ID: 2}},
			want:  4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nextID(tt.todos); got != tt.want {
				t.Errorf("nextID() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestAddAndLoad(t *testing.T) {
	ts := newTestStore(t)
	todo := Todo{Title: "Buy milk", IsDone: false, CreatedAt: time.Now()}

	if err := ts.addTodo(todo); err != nil {
		t.Fatalf("addTodo() error = %v", err)
	}

	todos, err := ts.loadTodos()
	if err != nil {
		t.Fatalf("loadTodos() error = %v", err)
	}
	if len(todos) != 1 {
		t.Fatalf("expected 1 todo, got %d", len(todos))
	}
	if todos[0].Title != todo.Title || todos[0].IsDone != todo.IsDone {
		t.Errorf("loaded todo = %+v, want title=%q isDone=%v", todos[0], todo.Title, todo.IsDone)
	}
}

func TestLoadMissingFile(t *testing.T) {
	ts := newTestStore(t)

	todos, err := ts.loadTodos()
	if err != nil {
		t.Fatalf("loadTodos() on missing file should return empty slice, got error: %v", err)
	}
	if len(todos) != 0 {
		t.Errorf("expected empty slice, got %d todos", len(todos))
	}
}

func TestFindTodo(t *testing.T) {
	ts := newTestStoreWithTodos(t,
		Todo{Title: "First", CreatedAt: time.Now()},
		Todo{Title: "Second", CreatedAt: time.Now()},
	)

	todos, _ := ts.loadTodos()
	firstID := todos[0].ID

	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			name:    "existing todo",
			id:      firstID,
			wantErr: false,
		},
		{
			name:    "non-existing todo",
			id:      99,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ts.findTodo(tt.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("findTodo(%d) error = %v, wantErr %v", tt.id, err, tt.wantErr)
			}
			if !tt.wantErr && got == nil {
				t.Error("findTodo() returned nil without error")
			}
		})
	}
}

func TestEditTodo(t *testing.T) {
	tests := []struct {
		name       string
		mutateFn   func(Todo) Todo
		wantErr    bool
		wantIsDone bool
	}{
		{
			name: "toggles IsDone to true",
			mutateFn: func(todo Todo) Todo {
				todo.IsDone = true
				return todo
			},
			wantErr:    false,
			wantIsDone: true,
		},
		{
			name: "non-existing id returns error",
			mutateFn: func(_ Todo) Todo {
				return Todo{ID: 99, Title: "Ghost"}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := newTestStoreWithTodos(t, Todo{Title: "Edit me", IsDone: false, CreatedAt: time.Now()})
			todos, _ := ts.loadTodos()
			updated := tt.mutateFn(todos[0])

			err := ts.editTodo(updated)
			if (err != nil) != tt.wantErr {
				t.Fatalf("editTodo() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				reloaded, _ := ts.loadTodos()
				if reloaded[0].IsDone != tt.wantIsDone {
					t.Errorf("editTodo() IsDone = %v, want %v", reloaded[0].IsDone, tt.wantIsDone)
				}
			}
		})
	}
}

func TestDelTodo(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		wantErr    bool
		wantRemain int
	}{
		{
			name:       "deletes existing todo",
			id:         0,
			wantErr:    false,
			wantRemain: 1,
		},
		{
			name:       "non-existing id returns error",
			id:         99,
			wantErr:    true,
			wantRemain: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := newTestStoreWithTodos(t,
				Todo{Title: "Keep", CreatedAt: time.Now()},
				Todo{Title: "Also keep", CreatedAt: time.Now()},
			)
			todos, _ := ts.loadTodos()
			id := tt.id
			if tt.id == 0 {
				id = todos[0].ID
			}
			err := ts.delTodo(id)
			if (err != nil) != tt.wantErr {
				t.Errorf("delTodo(%d) error = %v, wantErr %v", id, err, tt.wantErr)
			}
			remaining, _ := ts.loadTodos()
			if len(remaining) != tt.wantRemain {
				t.Errorf("remaining todos = %d, want %d", len(remaining), tt.wantRemain)
			}
		})
	}
}

func TestClearTodos(t *testing.T) {
	tests := []struct {
		name       string
		todos      []Todo
		wantRemain int
	}{
		{
			name: "removes done, keeps pending",
			todos: []Todo{
				{Title: "Pending", IsDone: false, CreatedAt: time.Now()},
				{Title: "Done", IsDone: true, CreatedAt: time.Now()},
			},
			wantRemain: 1,
		},
		{
			name: "all done removed",
			todos: []Todo{
				{Title: "Done A", IsDone: true, CreatedAt: time.Now()},
				{Title: "Done B", IsDone: true, CreatedAt: time.Now()},
			},
			wantRemain: 0,
		},
		{
			name: "no done todos leaves list unchanged",
			todos: []Todo{
				{Title: "Pending A", IsDone: false, CreatedAt: time.Now()},
			},
			wantRemain: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := newTestStoreWithTodos(t, tt.todos...)
			if err := ts.clearTodos(); err != nil {
				t.Fatalf("clearTodos() error = %v", err)
			}
			remaining, _ := ts.loadTodos()
			if len(remaining) != tt.wantRemain {
				t.Errorf("remaining todos = %d, want %d", len(remaining), tt.wantRemain)
			}
			for _, r := range remaining {
				if r.IsDone {
					t.Errorf("clearTodos() left a done todo: %q", r.Title)
				}
			}
		})
	}
}
