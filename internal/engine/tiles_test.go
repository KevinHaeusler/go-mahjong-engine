package engine

import (
	"testing"
)

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
		tile := NewTile(SuitManzu, 1)
		if tile.Suit() != SuitManzu {
			t.Errorf("expected SuitManzu, got %v", tile.Suit())
		}
		if tile.Rank() != 1 {
			t.Errorf("expected rank 1, got %v", tile.Rank())
		}
	})

	t.Run("panic on invalid rank low", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		NewTile(SuitManzu, 0)
	})

	t.Run("panic on invalid rank high", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		NewTile(SuitManzu, 10)
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
		{NewTile(SuitManzu, 1), false, true, true, true, false, false, "1m"},
		{NewTile(SuitManzu, 5), false, true, false, false, false, false, "5m"},
		{NewTile(SuitManzu, 9), false, true, true, true, false, false, "9m"},
		{NewTile(SuitPinzu, 1), false, true, true, true, false, false, "1p"},
		{NewTile(SuitSouzu, 9), false, true, true, true, false, false, "9s"},
		{NewTile(SuitHonor, 1), true, false, false, true, true, false, "E"},
		{NewTile(SuitHonor, 4), true, false, false, true, true, false, "N"},
		{NewTile(SuitHonor, 5), true, false, false, true, false, true, "P"},
		{NewTile(SuitHonor, 7), true, false, false, true, false, true, "C"},
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
		})
	}
}

func TestParseTile(t *testing.T) {
	tests := []struct {
		input   string
		want    Tile
		wantErr bool
	}{
		{"1m", NewTile(SuitManzu, 1), false},
		{"9p", NewTile(SuitPinzu, 9), false},
		{"3s", NewTile(SuitSouzu, 3), false},
		{"E", NewTile(SuitHonor, 1), false},
		{"S", NewTile(SuitHonor, 2), false},
		{"W", NewTile(SuitHonor, 3), false},
		{"N", NewTile(SuitHonor, 4), false},
		{"P", NewTile(SuitHonor, 5), false},
		{"F", NewTile(SuitHonor, 6), false},
		{"C", NewTile(SuitHonor, 7), false},
		{" 1m ", NewTile(SuitManzu, 1), false},
		{"e", NewTile(SuitHonor, 1), false},
		{"", 0, true},
		{"X", 0, true},
		{"1x", 0, true},
		{"0m", 0, true},
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
