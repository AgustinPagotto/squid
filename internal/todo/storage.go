package todo

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
}

type TodoStorage struct {
	Path string
}

func (ns *TodoStorage) findTodo(id int) (*Todo, error) {
	todos, err := ns.loadTodos()
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

func (ns *TodoStorage) editTodo(todo Todo) error {
	todos, err := ns.loadTodos()
	if err != nil {
		return err
	}
	for i := range todos {
		if todos[i].ID == todo.ID {
			todos[i] = todo
			return ns.persistTodos(todos)
		}
	}
	return fmt.Errorf("todo with id %d not found", todo.ID)
}

func (ns *TodoStorage) addTodo(todo Todo) error {
	todos, err := ns.loadTodos()
	if err != nil {
		return err
	}
	todo.ID = nextID(todos)
	todos = append(todos, todo)
	return ns.persistTodos(todos)
}

func (ns *TodoStorage) delTodo(id int) error {
	todos, err := ns.loadTodos()
	if err != nil {
		return err
	}
	for i := range todos {
		if todos[i].ID == id {
			return ns.persistTodos(slices.Delete(todos, i, i+1))
		}
	}
	return fmt.Errorf("todo with id %d not found", id)
}

func (ns *TodoStorage) persistTodos(todos []Todo) error {
	jsonTodos, err := json.MarshalIndent(todos, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(ns.Path, jsonTodos, config.FilePerm)
}

func (ns *TodoStorage) loadTodos() ([]Todo, error) {
	file, err := os.ReadFile(ns.Path)
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
