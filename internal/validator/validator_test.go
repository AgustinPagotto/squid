package validator

import "testing"

func TestHasArgAmount(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		amount int
		want   bool
	}{
		{"exact match", []string{"edit", "1"}, 2, true},
		{"more than required", []string{"edit", "1", "extra"}, 2, true},
		{"one short", []string{"edit"}, 2, false},
		{"empty args", []string{}, 1, false},
		{"zero required", []string{}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasArgAmount(tt.args, tt.amount); got != tt.want {
				t.Errorf("HasArgAmount(%v, %d) = %v, want %v", tt.args, tt.amount, got, tt.want)
			}
		})
	}
}
