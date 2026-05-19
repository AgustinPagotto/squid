package validator

import "testing"

func TestHasArgAmount(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		amount int
		want   bool
	}{
		{
			name:   "exact match",
			args:   []string{"edit", "1"},
			amount: 2,
			want:   true,
		},
		{
			name:   "more than required",
			args:   []string{"edit", "1", "extra"},
			amount: 2,
			want:   true,
		},
		{
			name:   "one short",
			args:   []string{"edit"},
			amount: 2,
			want:   false,
		},
		{
			name:   "empty args",
			args:   []string{},
			amount: 1,
			want:   false,
		},
		{
			name:   "zero required",
			args:   []string{},
			amount: 0,
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasArgAmount(tt.args, tt.amount); got != tt.want {
				t.Errorf("HasArgAmount(%v, %d) = %v, want %v", tt.args, tt.amount, got, tt.want)
			}
		})
	}
}

func TestValidateAndParseID(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		loc       int
		wantId    int
		wantError bool
	}{
		{
			name:      "valid id",
			args:      []string{"edit", "squid", "hello", "2"},
			loc:       3,
			wantId:    2,
			wantError: false,
		},
		{
			name:      "id is not a number",
			args:      []string{"edit", "squid", "hello"},
			loc:       1,
			wantError: true,
		},
		{
			name:      "not enough args",
			args:      []string{},
			loc:       1,
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := ValidateAndParseID(tt.args, tt.loc)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateAndParseID() error = %v, wantError %v", err, tt.wantError)
			}
			if !tt.wantError && id != tt.wantId {
				t.Errorf("ValidateAndParseID() id = %d, want %d", id, tt.wantId)
			}
		})
	}
}

func TestValidateTodoText(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{
			name:      "valid text",
			input:     "buy milk",
			wantError: false,
		},
		{
			name:      "empty text",
			input:     "",
			wantError: true,
		},
		{
			name:      "whitespace only",
			input:     "   ",
			wantError: true,
		},
		{
			name:      "exactly 280 chars",
			input:     string(make([]byte, 280)),
			wantError: false,
		},
		{
			name:      "exceeds 280 chars",
			input:     string(make([]byte, 281)),
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTodoText(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateTodoText() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateAndExtractSearch(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		want      string
		wantError bool
	}{
		{
			name:      "valid query",
			args:      []string{"search", "my query"},
			want:      "my query",
			wantError: false,
		},
		{
			name:      "multi-word query joined",
			args:      []string{"search", "my", "query"},
			want:      "my query",
			wantError: false,
		},
		{
			name:      "missing query arg",
			args:      []string{"search"},
			wantError: true,
		},
		{
			name:      "query too short",
			args:      []string{"search", "a"},
			wantError: true,
		},
		{
			name:      "query too long",
			args:      []string{"search", string(make([]byte, 101))},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateAndExtractSearch(tt.args)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateAndExtractSearch() error = %v, wantError %v", err, tt.wantError)
			}
			if !tt.wantError && got != tt.want {
				t.Errorf("ValidateAndExtractSearch() = %q, want %q", got, tt.want)
			}
		})
	}
}
