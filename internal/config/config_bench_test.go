package config

import (
	"os"
	"path/filepath"
	"testing"
)

// findWithReadDir checks for .squid/ by reading all directory entries.
func findWithReadDir(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}
	for _, entry := range entries {
		if entry.Name() == ".squid" && entry.IsDir() {
			return dir, nil
		}
	}
	return "", ErrRootNotFound
}

// findWithStat checks for .squid/ with a single targeted stat call.
func findWithStat(dir string) (string, error) {
	target := filepath.Join(dir, ".squid")
	info, err := os.Stat(target)
	if err == nil && info.IsDir() {
		return dir, nil
	}
	return "", ErrRootNotFound
}

// setupFixture creates a temp dir tree and places .squid/ inside it.
// Returns the root dir that contains .squid/.
func setupFixture(b *testing.B) string {
	b.Helper()
	root := b.TempDir()
	if err := os.Mkdir(filepath.Join(root, ".squid"), 0755); err != nil {
		b.Fatalf("failed to create .squid fixture: %v", err)
	}
	// Add sibling entries to make the directory realistic
	for _, name := range []string{"README.md", "main.go", "go.mod", ".git"} {
		if err := os.Mkdir(filepath.Join(root, name), 0755); err != nil {
			b.Fatalf("failed to create fixture entry %s: %v", name, err)
		}
	}
	return root
}

func BenchmarkFindWithReadDir(b *testing.B) {
	root := setupFixture(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := findWithReadDir(root); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFindWithStat(b *testing.B) {
	root := setupFixture(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := findWithStat(root); err != nil {
			b.Fatal(err)
		}
	}
}
