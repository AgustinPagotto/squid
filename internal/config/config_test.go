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
	t.Cleanup(func() {
		if err := os.Chdir(original); err != nil {
			t.Errorf("cleanup: failed to restore cwd: %v", err)
		}
	})
	return tmp
}

func TestFindRoot(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T, root string)
		wantErr error
		wantDir bool
	}{
		{
			name:    ".squid exists in cwd",
			setup:   func(t *testing.T, root string) {
				if err := os.Mkdir(filepath.Join(root, ".squid"), DirPerm); err != nil {
					t.Fatalf("setup: %v", err)
				}
			},
			wantDir: true,
		},
		{
			name: ".squid exists in parent dir",
			setup: func(t *testing.T, root string) {
				if err := os.Mkdir(filepath.Join(root, ".squid"), DirPerm); err != nil {
					t.Fatalf("setup: %v", err)
				}
				if err := os.Mkdir(filepath.Join(root, "subdir"), DirPerm); err != nil {
					t.Fatalf("setup: %v", err)
				}
				if err := os.Chdir(filepath.Join(root, "subdir")); err != nil {
					t.Fatalf("setup: %v", err)
				}
			},
			wantDir: true,
		},
		{
			name:    "no .squid anywhere",
			setup:   func(t *testing.T, root string) {},
			wantErr: ErrRootNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := chdirTemp(t)
			tt.setup(t, root)

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

func TestFindShellConfigurationFile(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T, home string)
		fileName string
		wantErr  error
		wantPath func(home string) string
	}{
		{
			name:     "file at home root",
			fileName: ".zshrc",
			setup: func(t *testing.T, home string) {
				if err := os.WriteFile(filepath.Join(home, ".zshrc"), []byte{}, FilePerm); err != nil {
					t.Fatalf("setup: %v", err)
				}
			},
			wantPath: func(home string) string { return filepath.Join(home, ".zshrc") },
		},
		{
			name:     "file inside .config",
			fileName: ".bashrc",
			setup: func(t *testing.T, home string) {
				if err := os.MkdirAll(filepath.Join(home, ".config"), DirPerm); err != nil {
					t.Fatalf("setup: %v", err)
				}
				if err := os.WriteFile(filepath.Join(home, ".config", ".bashrc"), []byte{}, FilePerm); err != nil {
					t.Fatalf("setup: %v", err)
				}
			},
			wantPath: func(home string) string { return filepath.Join(home, ".config", ".bashrc") },
		},
		{
			name:     "file inside .config/zsh",
			fileName: ".zshrc",
			setup: func(t *testing.T, home string) {
				if err := os.MkdirAll(filepath.Join(home, ".config", "zsh"), DirPerm); err != nil {
					t.Fatalf("setup: %v", err)
				}
				if err := os.WriteFile(filepath.Join(home, ".config", "zsh", ".zshrc"), []byte{}, FilePerm); err != nil {
					t.Fatalf("setup: %v", err)
				}
			},
			wantPath: func(home string) string {
				return filepath.Join(home, ".config", "zsh", ".zshrc")
			},
		},
		{
			name:     "file not found anywhere",
			fileName: ".zshrc",
			setup:    func(t *testing.T, home string) {},
			wantErr:  ErrConfigFileNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			home := t.TempDir()
			tt.setup(t, home)

			got, err := FindShellConfigurationFile(tt.fileName, home)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("FindShellConfigurationFile() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("FindShellConfigurationFile() unexpected error = %v", err)
			}
			if want := tt.wantPath(home); got != want {
				t.Errorf("FindShellConfigurationFile() = %q, want %q", got, want)
			}
		})
	}
}

func TestGetEditor(t *testing.T) {
	tests := []struct {
		name   string
		envVal string
		want   string
	}{
		{
			name:   "uses $EDITOR when set",
			envVal: "vim",
			want:   "vim",
		},
		{
			name:   "falls back to nano when unset",
			envVal: "",
			want:   "nano",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("EDITOR", tt.envVal)
			if got := getEditor(); got != tt.want {
				t.Errorf("getEditor() = %q, want %q", got, tt.want)
			}
		})
	}
}
