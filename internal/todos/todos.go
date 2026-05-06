package todos

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/AgustinPagotto/squid/internal/validator"
)

type Todo struct {
	ID        int
	Title     string
	IsDone    bool
	CreatedAt time.Time
}

func Handle(args []string, ts TodoStorageInterface) {
	if len(args) == 0 {
		fmt.Println("expected subcommand (add, list, etc.)")
		return
	}
	switch args[0] {
	case "-h":
		fmt.Print(commandHelp)
	case "add", "a":
		if !validator.HasArgAmount(args, 2) {
			fmt.Println("expected a text for the new todo item")
			return
		}
		err := handleAdd(args, ts)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("todo added successfully")
	case "toggle", "t":
		err := handleToggle(ts)
		if err != nil {
			fmt.Println(err)
		}
	case "del", "d":
		id, err := validator.ValidateAndParseID(args, 1)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = handleDelete(ts, id)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("todo deleted")
	case "list", "l":
		err := handleList(ts)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "show", "s":
		id, err := validator.ValidateAndParseID(args, 1)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = handleShow(ts, id)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "clear", "c":
		err := handleClear(ts)
		if err != nil {
			fmt.Println(err)
			return
		}
	default:
		fmt.Println("unknown subcommand:", args[0])
	}
}

func handleAdd(args []string, ts TodoStorageInterface) error {
	if err := validator.ValidateTodoText(args[1]); err != nil {
		return err
	}
	todo := Todo{Title: args[1], IsDone: false, CreatedAt: time.Now()}
	return ts.addTodo(todo)
}

func handleToggle(ts TodoStorageInterface) error {
	todos, err := ts.loadTodos()
	for _, t := range todos {
		if t.IsDone {
			fmt.Printf("  %2d  [✓] %s\n", t.ID, t.Title)
		} else {
			fmt.Printf("  %2d  [ ] %s\n", t.ID, t.Title)
		}
	}
	fmt.Print("\nToggle todo (id): ")
	var input string
	fmt.Scanln(&input)
	parsedId, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		return fmt.Errorf("invalid id")
	}
	todo, err := ts.findTodo(parsedId)
	if err != nil {
		return err
	}
	todo.IsDone = !todo.IsDone
	return ts.editTodo(*todo)
}

func handleDelete(ts TodoStorageInterface, id int) error {
	todo, err := ts.findTodo(id)
	if err != nil {
		return err
	}
	fmt.Printf("Are you sure you want to delete todo: %q? [Y/n]: ", todo.Title)
	var input string
	fmt.Scanln(&input)
	input = strings.TrimSpace(strings.ToLower(input))
	if input != "" && input != "y" && input != "yes" {
		return fmt.Errorf("delete aborted")
	}
	return ts.delTodo(id)
}

func handleList(ts TodoStorageInterface) error {
	todos, err := ts.loadTodos()
	if err != nil {
		return err
	}
	if len(todos) == 0 {
		fmt.Println("no todos yet — run 'squid todo add <text>' to create one")
		return nil
	}

	var pending, done []Todo
	for _, t := range todos {
		if t.IsDone {
			done = append(done, t)
		} else {
			pending = append(pending, t)
		}
	}

	fmt.Println()
	for _, t := range pending {
		fmt.Printf("  %2d  [ ] %s\n", t.ID, t.Title)
	}

	if len(done) > 0 {
		fmt.Println("\n  ── done ──────────────────────────")
		for _, t := range done {
			fmt.Printf("  %2d  [✓] %s\n", t.ID, t.Title)
		}
	}

	fmt.Printf("\n  %d done · %d pending\n\n", len(done), len(pending))
	return nil
}

func handleShow(ts TodoStorageInterface, id int) error {
	todo, err := ts.findTodo(id)
	if err != nil {
		return err
	}
	status := "[ ] pending"
	if todo.IsDone {
		status = "[✓] done"
	}
	fmt.Printf("%s\n\n", todo.Title)
	fmt.Printf("Status:  %s\n", status)
	fmt.Printf("Created: %s\n", todo.CreatedAt.Format("2006-01-02"))
	return nil
}

func handleClear(ts TodoStorageInterface) error {
	fmt.Print("This will delete all of the done todos [Y/n]: ")
	var input string
	fmt.Scanln(&input)
	input = strings.TrimSpace(strings.ToLower(input))
	if input != "" && input != "y" && input != "yes" {
		return fmt.Errorf("clear aborted")
	}
	return ts.clearTodos()
}
