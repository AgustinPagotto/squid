package context

import (
	"strings"
	"testing"
)

func TestParseAliases(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		{
			name:  "valid single alias",
			input: `start="npm run dev"`,
			want:  []string{`start="npm run dev"`},
		},
		{
			name:  "multiple valid aliases",
			input: "start=\"npm run dev\"\ntest=\"go test ./...\"",
			want:  []string{`start="npm run dev"`, `test="go test ./..."`},
		},
		{
			name:  "comment lines are skipped",
			input: "# this is a comment\nstart=\"npm run dev\"",
			want:  []string{`start="npm run dev"`},
		},
		{
			name:  "blank lines are skipped",
			input: "\n\nstart=\"npm run dev\"\n\n",
			want:  []string{`start="npm run dev"`},
		},
		{
			name:  "invalid lines are silently skipped",
			input: "not-valid\nstart=\"npm run dev\"",
			want:  []string{`start="npm run dev"`},
		},
		{
			name:  "single-quoted value is valid",
			input: `start='npm run dev'`,
			want:  []string{`start='npm run dev'`},
		},
		{
			name:    "empty double-quoted value is invalid",
			input:   `start=""`,
			wantErr: true,
		},
		{
			name:    "empty single-quoted value is invalid",
			input:   `start=''`,
			wantErr: true,
		},
		{
			name:    "alias name exceeds max length",
			input:   strings.Repeat("a", MaxAliasNameLen+1) + `="cmd"`,
			wantErr: true,
		},
		{
			name:    "command value exceeds max length",
			input:   `start="` + strings.Repeat("x", MaxAliasCommandLen+1) + `"`,
			wantErr: true,
		},
		{
			name:    "only comments returns error",
			input:   "# just a comment\n# another",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAliases(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseAliases() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Fatalf("ParseAliases() len = %d, want %d", len(got), len(tt.want))
				}
				for i := range got {
					if got[i] != tt.want[i] {
						t.Errorf("ParseAliases()[%d] = %q, want %q", i, got[i], tt.want[i])
					}
				}
			}
		})
	}
}

func TestHandleContextCreation(t *testing.T) {
	tests := []struct {
		name    string
		entries []string
		wantLen int
		wantErr bool
	}{
		{
			name:    "creates aliases with sequential IDs",
			entries: []string{`start="npm run dev"`, `test="go test ./..."`},
			wantLen: 2,
		},
		{
			name:    "empty entries persists empty slice",
			entries: []string{},
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &ContextStorage{Path: t.TempDir() + "/context.json"}
			err := HandleContextCreation(tt.entries, cs)
			if (err != nil) != tt.wantErr {
				t.Fatalf("HandleContextCreation() error = %v, wantErr %v", err, tt.wantErr)
			}
			aliases, _ := cs.loadAliases()
			if len(aliases) != tt.wantLen {
				t.Errorf("aliases len = %d, want %d", len(aliases), tt.wantLen)
			}
			for i, a := range aliases {
				if a.ID != i {
					t.Errorf("alias[%d].ID = %d, want %d", i, a.ID, i)
				}
				if a.AliasCommand != tt.entries[i] {
					t.Errorf("alias[%d].AliasCommand = %q, want %q", i, a.AliasCommand, tt.entries[i])
				}
			}
		})
	}
}

func TestHandleList(t *testing.T) {
	tests := []struct {
		name    string
		aliases []Alias
		wantErr bool
	}{
		{
			name:    "lists aliases",
			aliases: []Alias{{ID: 0, AliasCommand: `start="npm run dev"`}},
			wantErr: false,
		},
		{
			name:    "empty storage prints message",
			aliases: []Alias{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &ContextStorage{Path: t.TempDir() + "/context.json"}
			if err := cs.persistContext(tt.aliases); err != nil {
				t.Fatalf("persistContext() error = %v", err)
			}
			err := handleList(cs)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleActivate(t *testing.T) {
	tests := []struct {
		name    string
		aliases []Alias
		wantErr bool
	}{
		{
			name:    "prints alias export lines",
			aliases: []Alias{{ID: 0, AliasCommand: `start="npm run dev"`}},
			wantErr: false,
		},
		{
			name:    "empty aliases prints nothing",
			aliases: []Alias{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &ContextStorage{Path: t.TempDir() + "/context.json"}
			if err := cs.persistContext(tt.aliases); err != nil {
				t.Fatalf("persistContext() error = %v", err)
			}
			if err := handleActivate(cs); (err != nil) != tt.wantErr {
				t.Errorf("handleActivate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// --- ContextStorage (real file-backed) ---

func TestContextLoadAliases(t *testing.T) {
	tests := []struct {
		name              string
		aliases           []Alias
		wantAmountAliases int
		wantErr           bool
	}{
		{
			name:              "one valid alias to save",
			aliases:           []Alias{{ID: 0, AliasCommand: `start="npm run dev"`}},
			wantAmountAliases: 1,
			wantErr:           false,
		},
		{
			name:              "empty file returns nothing",
			aliases:           []Alias{},
			wantAmountAliases: 0,
			wantErr:           false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &ContextStorage{Path: t.TempDir() + "/context.json"}
			if err := cs.persistContext(tt.aliases); err != nil {
				t.Fatalf("persistContext() error = %v", err)
			}
			aliases, err := cs.loadAliases()
			if tt.wantAmountAliases != len(aliases) {
				t.Errorf("amount of aliases = %d, want %d", len(aliases), tt.wantAmountAliases)
			}
			if tt.wantErr != (err != nil) {
				t.Errorf("error= %q", err)
			}
		})
	}
}

func TestContextStorageFindAlias(t *testing.T) {
	aliases := []Alias{
		{ID: 0, AliasCommand: `start="npm run dev"`},
		{ID: 1, AliasCommand: `test="go test ./..."`},
	}
	cs := &ContextStorage{Path: t.TempDir() + "/context.json"}
	if err := cs.persistContext(aliases); err != nil {
		t.Fatalf("persistContext() error = %v", err)
	}

	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			name:    "existing alias",
			id:      0,
			wantErr: false,
		},
		{
			name:    "another existing alias",
			id:      1,
			wantErr: false,
		},
		{
			name:    "non-existing alias",
			id:      88,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cs.findAlias(tt.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("findAlias(%d) error = %v, wantErr %v", tt.id, err, tt.wantErr)
			}
			if !tt.wantErr && got == nil {
				t.Error("findAlias() returned nil without error")
			}
		})
	}
}

func TestContextStorageDelAlias(t *testing.T) {
	aliases := []Alias{
		{ID: 0, AliasCommand: `start="npm run dev"`},
		{ID: 1, AliasCommand: `test="go test ./..."`},
	}
	cs := &ContextStorage{Path: t.TempDir() + "/context.json"}
	if err := cs.persistContext(aliases); err != nil {
		t.Fatalf("persistContext() error = %v", err)
	}

	tests := []struct {
		name       string
		id         int
		wantErr    bool
		wantRemain int
	}{
		{
			name:       "deletes existing alias",
			id:         0,
			wantErr:    false,
			wantRemain: 1,
		},
		{
			name:       "non-existing alias returns error",
			id:         99,
			wantErr:    true,
			wantRemain: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cs.delAlias(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("delAlias(%d) error = %v, wantErr %v", tt.id, err, tt.wantErr)
			}
			aliases, _ := cs.loadAliases()
			if len(aliases) != tt.wantRemain {
				t.Errorf("remaining aliases = %d, want %d", len(aliases), tt.wantRemain)
			}
		})
	}
}

func TestNextID(t *testing.T) {
	tests := []struct {
		name    string
		aliases []Alias
		want    int
	}{
		{
			name:    "no aliases",
			aliases: []Alias{},
			want:    0,
		},
		{
			name:    "single alias",
			aliases: []Alias{{ID: 3}},
			want:    4,
		},
		{
			name:    "continues after max id",
			aliases: []Alias{{ID: 0}, {ID: 5}, {ID: 2}},
			want:    6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nextID(tt.aliases); got != tt.want {
				t.Errorf("nextID() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestAssignIDs(t *testing.T) {
	tests := []struct {
		name     string
		existing []Alias
		parsed   []string
		wantIDs  []int
	}{
		{
			name:     "fresh set assigns 0-based ids",
			existing: []Alias{},
			parsed:   []string{`start="npm run dev"`, `test="go test ./..."`},
			wantIDs:  []int{0, 1},
		},
		{
			name: "existing aliases preserve their ids",
			existing: []Alias{
				{ID: 0, AliasCommand: `start="npm run dev"`},
				{ID: 1, AliasCommand: `test="go test ./..."`},
			},
			parsed:  []string{`start="npm run dev"`, `test="go test ./..."`},
			wantIDs: []int{0, 1},
		},
		{
			name: "deleted alias leaves gap, new alias does not reuse id",
			existing: []Alias{
				{ID: 0, AliasCommand: `start="npm run dev"`},
				{ID: 1, AliasCommand: `test="go test ./..."`},
				{ID: 2, AliasCommand: `build="go build ./..."`},
			},
			// user deleted id=1 (test) in the editor and added a new one
			parsed:  []string{`start="npm run dev"`, `build="go build ./..."`, `lint="golangci-lint run"`},
			wantIDs: []int{0, 2, 3},
		},
		{
			name: "new alias gets max+1 even with gaps",
			existing: []Alias{
				{ID: 0, AliasCommand: `start="npm run dev"`},
				{ID: 5, AliasCommand: `test="go test ./..."`},
			},
			parsed:  []string{`start="npm run dev"`, `test="go test ./..."`, `new="echo hi"`},
			wantIDs: []int{0, 5, 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := assignIDs(tt.existing, tt.parsed)
			if len(got) != len(tt.wantIDs) {
				t.Fatalf("assignIDs() len = %d, want %d", len(got), len(tt.wantIDs))
			}
			for i, a := range got {
				if a.ID != tt.wantIDs[i] {
					t.Errorf("assignIDs()[%d].ID = %d, want %d", i, a.ID, tt.wantIDs[i])
				}
				if a.AliasCommand != tt.parsed[i] {
					t.Errorf("assignIDs()[%d].AliasCommand = %q, want %q", i, a.AliasCommand, tt.parsed[i])
				}
			}
		})
	}
}
