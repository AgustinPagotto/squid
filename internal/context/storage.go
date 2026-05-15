package context

import (
	"encoding/json"
	"os"

	"github.com/AgustinPagotto/squid/internal/config"
)

type ContextStorageInterface interface {
	loadAliases() ([]Alias, error)
	persistContext([]Alias) error
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
