package engine

import (
	"slices"
	"testing"
)

func TestParseHandCompact(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantLen int
		wantErr bool
	}{
		{
			name:    "14-tile hand (example 1)",
			input:   "123m456789s12345z",
			wantLen: 14,
			wantErr: false,
		},
		{
			name:    "14-tile hand (example 2)",
			input:   "1m2m3m4s5s6s7s8s9s1z2z3z4z5z",
			wantLen: 14,
			wantErr: false,
		},
		{
			name:    "Short hand (not 14)",
			input:   "123m",
			wantLen: 3,
			wantErr: false,
		},
		{
			name:    "with red fives (rank 0)",
			input:   "0m0p0s",
			wantLen: 3,
			wantErr: false,
		},
		{
			name:    "with spaces",
			input:   " 123m 456p 789s 12345z ",
			wantLen: 14,
			wantErr: false,
		},
		{
			name:    "invalid rank: 8z",
			input:   "8z",
			wantErr: true,
		},
		{
			name:    "invalid rank: 0z",
			input:   "0z",
			wantErr: true,
		},
		{
			name:    "missing suit at end",
			input:   "123m456",
			wantErr: true,
		},
		{
			name:    "missing ranks before suit",
			input:   "m",
			wantErr: true,
		},
		{
			name:    "unexpected character",
			input:   "123x",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseHandCompact(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseHandCompact(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && len(got) != tt.wantLen {
				t.Errorf("ParseHandCompact(%q) length = %d, want %d", tt.input, len(got), tt.wantLen)
			}
		})
	}
}

func TestParseHandCompact_Content(t *testing.T) {
	input := "123m0p45z"
	got, err := ParseHandCompact(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Helper to create tiles without error checks in expectation
	mk := func(s Suit, r int) Tile {
		t, _ := NewTile(s, r)
		return t
	}

	want := []Tile{
		mk(SuitManzu, 1),
		mk(SuitManzu, 2),
		mk(SuitManzu, 3),
		mk(SuitPinzu, 0),
		mk(SuitHonor, 4),
		mk(SuitHonor, 5),
	}

	if !slices.Equal(got, want) {
		t.Errorf("ParseHandCompact(%q) = %v, want %v", input, got, want)
	}
}
