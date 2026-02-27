package engine

import (
	"testing"
)

func mustNewTile(t *testing.T, s Suit, rank int) Tile {
	t.Helper()
	tile, err := NewTile(s, rank)
	if err != nil {
		t.Fatalf("NewTile(%v, %d) failed: %v", s, rank, err)
	}
	return tile
}

func TestSuit_String(t *testing.T) {
	tests := []struct {
		suit Suit
		want string
	}{
		{SuitManzu, "m"},
		{SuitPinzu, "p"},
		{SuitSouzu, "s"},
		{SuitHonor, "z"},
		{Suit(99), "?"},
	}
	for _, tt := range tests {
		if got := tt.suit.String(); got != tt.want {
			t.Errorf("Suit.String() = %v, want %v", got, tt.want)
		}
	}
}

func TestNewTile(t *testing.T) {
	t.Run("valid tiles", func(t *testing.T) {
		tile, err := NewTile(SuitManzu, 1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if tile.Suit() != SuitManzu {
			t.Errorf("expected SuitManzu, got %v", tile.Suit())
		}
		if tile.Rank() != 1 {
			t.Errorf("expected rank 1, got %v", tile.Rank())
		}
	})

	t.Run("error on invalid rank low", func(t *testing.T) {
		_, err := NewTile(SuitManzu, -1)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("error on invalid rank high", func(t *testing.T) {
		_, err := NewTile(SuitManzu, 10)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}

func TestTile_Methods(t *testing.T) {
	tests := []struct {
		tile              Tile
		isHonor           bool
		isNum             bool
		isTerminal        bool
		isTerminalOrHonor bool
		isWind            bool
		isDragon          bool
		str               string
	}{
		{mustNewTile(t, SuitManzu, 1), false, true, true, true, false, false, "1m"},
		{mustNewTile(t, SuitManzu, 5), false, true, false, false, false, false, "5m"},
		{mustNewTile(t, SuitManzu, 9), false, true, true, true, false, false, "9m"},
		{mustNewTile(t, SuitPinzu, 1), false, true, true, true, false, false, "1p"},
		{mustNewTile(t, SuitSouzu, 9), false, true, true, true, false, false, "9s"},
		{mustNewTile(t, SuitHonor, 1), true, false, false, true, true, false, "1z"},
		{mustNewTile(t, SuitHonor, 4), true, false, false, true, true, false, "4z"},
		{mustNewTile(t, SuitHonor, 5), true, false, false, true, false, true, "5z"},
		{mustNewTile(t, SuitHonor, 6), true, false, false, true, false, true, "6z"},
		{mustNewTile(t, SuitHonor, 7), true, false, false, true, false, true, "7z"},
	}

	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			if got := tt.tile.IsHonor(); got != tt.isHonor {
				t.Errorf("IsHonor() = %v, want %v", got, tt.isHonor)
			}
			if got := tt.tile.IsNumbered(); got != tt.isNum {
				t.Errorf("IsNumbered() = %v, want %v", got, tt.isNum)
			}
			if got := tt.tile.IsTerminal(); got != tt.isTerminal {
				t.Errorf("IsTerminal() = %v, want %v", got, tt.isTerminal)
			}
			if got := tt.tile.IsTerminalOrHonor(); got != tt.isTerminalOrHonor {
				t.Errorf("IsJunchan() = %v, want %v", got, tt.isTerminalOrHonor)
			}
			if got := tt.tile.IsWind(); got != tt.isWind {
				t.Errorf("IsWind() = %v, want %v", got, tt.isWind)
			}
			if got := tt.tile.IsDragon(); got != tt.isDragon {
				t.Errorf("IsDragon() = %v, want %v", got, tt.isDragon)
			}
			if got := tt.tile.String(); got != tt.str {
				t.Errorf("String() = %v, want %v", got, tt.str)
			}

			// Test new flags
			if tt.tile.IsDora() {
				t.Errorf("expected IsDora() to be false initially")
			}
			if tt.tile.IsUra() {
				t.Errorf("expected IsUra() to be false initially")
			}
			isRed := tt.tile.IsNumbered() && tt.tile.Rank() == 0
			if got := tt.tile.IsRed(); got != isRed {
				t.Errorf("IsRed() = %v, want %v", got, isRed)
			}

			// Test SetDora/SetUra
			doraTile := tt.tile.SetDora(true)
			if !doraTile.IsDora() {
				t.Errorf("SetDora(true) failed")
			}
			if doraTile.SetDora(false).IsDora() {
				t.Errorf("SetDora(false) failed")
			}

			uraTile := tt.tile.SetUra(true)
			if !uraTile.IsUra() {
				t.Errorf("SetUra(true) failed")
			}
			if uraTile.SetUra(false).IsUra() {
				t.Errorf("SetUra(false) failed")
			}
		})
	}
}

func TestTile_RedFive(t *testing.T) {
	redFive, _ := NewRedFive(SuitManzu)
	if !redFive.IsRed() {
		t.Errorf("NewRedFive should be IsRed")
	}
	if redFive.Rank() != 0 {
		t.Errorf("NewRedFive rank should be 0, got %d", redFive.Rank())
	}

	normalFive := mustNewTile(t, SuitManzu, 5)
	if normalFive.IsRed() {
		t.Errorf("Normal five should not be IsRed")
	}
}

func TestParseTile(t *testing.T) {
	tests := []struct {
		input   string
		want    Tile
		wantErr bool
	}{
		{"1m", mustNewTile(t, SuitManzu, 1), false},
		{"9p", mustNewTile(t, SuitPinzu, 9), false},
		{"3s", mustNewTile(t, SuitSouzu, 3), false},
		{"E", mustNewTile(t, SuitHonor, 1), false},
		{"S", mustNewTile(t, SuitHonor, 2), false},
		{"W", mustNewTile(t, SuitHonor, 3), false},
		{"N", mustNewTile(t, SuitHonor, 4), false},
		{"7z", mustNewTile(t, SuitHonor, 7), false},
		{"6z", mustNewTile(t, SuitHonor, 6), false},
		{"5z", mustNewTile(t, SuitHonor, 5), false},
		{" 1m ", mustNewTile(t, SuitManzu, 1), false},
		{"e", mustNewTile(t, SuitHonor, 1), false},
		{"", 0, true},
		{"X", 0, true},
		{"1x", 0, true},
		{"0m", mustNewTile(t, SuitManzu, 0), false},
		{"10m", 0, true},
		{"ABC", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseTile(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseTile() = %v, want %v", got, tt.want)
			}
		})
	}
}
