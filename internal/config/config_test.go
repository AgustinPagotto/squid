package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func chdirTemp(t *testing.T) string {
	t.Helper()
	original, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not get cwd: %v", err)
	}
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("could not chdir to temp dir: %v", err)
	}
	t.Cleanup(func() { os.Chdir(original) })
	return tmp
}

func TestFindRoot(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(root string)
		wantErr error
		wantDir bool
	}{
		{
			name:    ".squid exists in cwd",
			setup:   func(root string) { os.Mkdir(filepath.Join(root, ".squid"), DirPerm) },
			wantDir: true,
		},
		{
			name: ".squid exists in parent dir",
			setup: func(root string) {
				os.Mkdir(filepath.Join(root, ".squid"), DirPerm)
				os.Mkdir(filepath.Join(root, "subdir"), DirPerm)
				os.Chdir(filepath.Join(root, "subdir"))
			},
			wantDir: true,
		},
		{
			name:    "no .squid anywhere",
			setup:   func(root string) {},
			wantErr: ErrRootNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := chdirTemp(t)
			tt.setup(root)

			got, err := FindRoot()

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("FindRoot() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("FindRoot() unexpected error = %v", err)
			}
			if !tt.wantDir || got == "" {
				t.Errorf("FindRoot() returned empty path, want a valid .squid path")
			}
		})
	}
}

func TestInit(t *testing.T) {
	tests := []struct {
		name  string
		setup func(root string)
	}{
		{
			name:  "creates .squid when not present",
			setup: func(root string) {},
		},
		{
			name:  "no error when .squid already exists",
			setup: func(root string) { os.Mkdir(filepath.Join(root, ".squid"), DirPerm) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := chdirTemp(t)
			tt.setup(root)

			err := Init([]string{})
			if err != nil {
				t.Fatalf("Init() error = %v", err)
			}

			info, statErr := os.Stat(filepath.Join(root, ".squid"))
			if statErr != nil {
				t.Errorf(".squid was not created: %v", statErr)
			} else if info.Mode().Perm() != DirPerm {
				t.Errorf(".squid permissions = %v, want %v", info.Mode().Perm(), DirPerm)
			}
		})
	}
}
