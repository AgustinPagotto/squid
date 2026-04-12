package validator

import (
	"fmt"
	"strconv"
	"strings"
)

func HasArgAmount(args []string, amount int) bool {
	return len(args) >= amount
}

func ValidateAndParseID(args []string, loc int) (int, error) {
	if !HasArgAmount(args, loc+1) {
		return 0, fmt.Errorf("please provide an id")
	}
	id, err := strconv.Atoi(args[loc])
	if err != nil {
		return 0, fmt.Errorf("invalid id")
	}
	return id, nil
}

func ValidateTodoText(text string) error {
	text = strings.TrimSpace(text)
	if text == "" {
		return fmt.Errorf("todo text cannot be empty")
	}
	if len(text) > 280 {
		return fmt.Errorf("todo text cannot exceed 280 characters")
	}
	return nil
}

func ValidateAndExtractSearch(args []string) (string, error) {
	if !HasArgAmount(args, 2) {
		return "", fmt.Errorf("please provide a search query")
	}
	query := strings.TrimSpace(strings.Join(args[1:], " "))
	if query == "" {
		return "", fmt.Errorf("search query cannot be empty")
	}
	if len(query) < 2 {
		return "", fmt.Errorf("search query must be at least 2 characters")
	}
	if len(query) > 100 {
		return "", fmt.Errorf("search query cannot exceed 100 characters")
	}
	return query, nil
}
