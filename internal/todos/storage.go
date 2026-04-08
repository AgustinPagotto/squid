package todos

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"

	"github.com/AgustinPagotto/squid/internal/config"
)

type TodoStorageInterface interface {
	findTodo(id int) (*Todo, error)
	editTodo(todo Todo) error
	addTodo(todo Todo) error
	delTodo(id int) error
	loadTodos() ([]Todo, error)
	clearTodos() error
}

type TodoStorage struct {
	Path string
}

func (ts *TodoStorage) findTodo(id int) (*Todo, error) {
	todos, err := ts.loadTodos()
	if err != nil {
		return nil, err
	}
	for i := range todos {
		if todos[i].ID == id {
			return &todos[i], nil
		}
	}
	return nil, fmt.Errorf("todos with id %d not found", id)
}

func (ts *TodoStorage) editTodo(todo Todo) error {
	todos, err := ts.loadTodos()
	if err != nil {
		return err
	}
	for i := range todos {
		if todos[i].ID == todo.ID {
			todos[i] = todo
			return ts.persistTodos(todos)
		}
	}
	return fmt.Errorf("todo with id %d not found", todo.ID)
}

func (ts *TodoStorage) addTodo(todo Todo) error {
	todos, err := ts.loadTodos()
	if err != nil {
		return err
	}
	todo.ID = nextID(todos)
	todos = append(todos, todo)
	return ts.persistTodos(todos)
}

func (ts *TodoStorage) delTodo(id int) error {
	todos, err := ts.loadTodos()
	if err != nil {
		return err
	}
	for i := range todos {
		if todos[i].ID == id {
			return ts.persistTodos(slices.Delete(todos, i, i+1))
		}
	}
	return fmt.Errorf("todo with id %d not found", id)
}

func (ts *TodoStorage) persistTodos(todos []Todo) error {
	jsonTodos, err := json.MarshalIndent(todos, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(ts.Path, jsonTodos, config.FilePerm)
}

func (ts *TodoStorage) loadTodos() ([]Todo, error) {
	file, err := os.ReadFile(ts.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Todo{}, nil
		}
		return nil, err
	}
	var todos []Todo
	err = json.Unmarshal(file, &todos)
	return todos, err
}

func (ts *TodoStorage) clearTodos() error {
	todos, err := ts.loadTodos()
	if err != nil {
		return err
	}
	var pendingTodos []Todo
	for _, t := range todos {
		if !t.IsDone {
			pendingTodos = append(pendingTodos, t)
		}
	}
	return ts.persistTodos(pendingTodos)
}

func nextID(todos []Todo) int {
	if len(todos) == 0 {
		return 0
	}
	maxID := 0
	for _, todo := range todos {
		if todo.ID > maxID {
			maxID = todo.ID
		}
	}
	return maxID + 1
}
