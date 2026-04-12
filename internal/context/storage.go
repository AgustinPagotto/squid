package context

import (
	"encoding/json"
	"os"
)

type ContextStorageInterface interface {
	loadAliases() ([]Alias, error)
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
