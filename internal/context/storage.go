package context

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"

	"github.com/AgustinPagotto/squid/internal/config"
)

type ContextStorageInterface interface {
	loadAliases() ([]Alias, error)
	persistContext([]Alias) error
	findAlias(int) (*Alias, error)
	delAlias(int) error
}

type ContextStorage struct {
	Path string
}

func (cs *ContextStorage) loadAliases() ([]Alias, error) {
	file, err := os.ReadFile(cs.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Alias{}, nil
		}
		return nil, err
	}
	var aliases []Alias
	err = json.Unmarshal(file, &aliases)
	return aliases, err
}

func (cs *ContextStorage) persistContext(aliases []Alias) error {
	jsonNotes, err := json.MarshalIndent(aliases, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(cs.Path, jsonNotes, config.FilePerm)
}

func (cs *ContextStorage) findAlias(id int) (*Alias, error) {
	aliases, err := cs.loadAliases()
	if err != nil {
		return nil, err
	}
	for i := range aliases {
		if aliases[i].ID == id {
			return &aliases[i], nil
		}
	}
	return nil, fmt.Errorf("alias with id %d not found", id)
}

func (cs *ContextStorage) delAlias(id int) error {
	notes, err := cs.loadAliases()
	if err != nil {
		return err
	}
	for i := range notes {
		if notes[i].ID == id {
			return cs.persistContext(slices.Delete(notes, i, i+1))
		}
	}
	return fmt.Errorf("alias with id %d not found", id)
}

func nextID(aliases []Alias) int {
	if len(aliases) == 0 {
		return 0
	}
	maxID := 0
	for _, a := range aliases {
		if a.ID > maxID {
			maxID = a.ID
		}
	}
	return maxID + 1
}
