package engine

import (
	"slices"
	"testing"
)

func TestBuildWall(t *testing.T) {
	t.Run("default rules (3 red fives)", func(t *testing.T) {
		rules := DefaultRules()
		wall, err := BuildWall(rules)
		if err != nil {
			t.Fatalf("BuildWall(DefaultRules()) failed: %v", err)
		}

		if len(wall) != 136 {
			t.Errorf("expected 136 tiles, got %d", len(wall))
		}

		// Verify red fives: 1m, 1p, 1s
		redCount := 0
		for _, tile := range wall {
			if tile.IsRed() {
				redCount++
			}
		}
		if redCount != 3 {
			t.Errorf("expected 3 red fives, got %d", redCount)
		}

		// Check specific suits
		var redMan, redPin, redSou int
		for _, tile := range wall {
			if tile.IsRed() {
				switch tile.Suit() {
				case SuitManzu:
					redMan++
				case SuitPinzu:
					redPin++
				case SuitSouzu:
					redSou++
				default:
					panic("unhandled default case")
				}
			}
		}
		if redMan != 1 || redPin != 1 || redSou != 1 {
			t.Errorf("expected 1 of each red five, got m=%d, p=%d, s=%d", redMan, redPin, redSou)
		}

		// Verify distribution: each of 34 kinds has 4 copies
		counts := make(map[Tile]int)
		for _, tile := range wall {
			counts[tile]++
		}

		// There are 31 kinds with 4 copies each, and 3 kinds (5m/p/s) split between red and normal
		// 5m: 1 red, 3 normal
		// 5p: 1 red, 3 normal
		// 5s: 1 red, 3 normal
		// Total tiles: 31*4 + 3*1 + 3*3 = 124 + 3 + 9 = 136. Correct.

		// Let's verify each kind exists and has correct count
		for suit := SuitManzu; suit <= SuitSouzu; suit++ {
			for rank := 1; rank <= 9; rank++ {
				if rank == 5 {
					red, _ := NewRedFive(suit)
					normal, _ := NewTile(suit, 5)
					if counts[red] != 1 {
						t.Errorf("expected 1 red %v, got %d", red, counts[red])
					}
					if counts[normal] != 3 {
						t.Errorf("expected 3 normal %v, got %d", normal, counts[normal])
					}
				} else {
					tile, _ := NewTile(suit, rank)
					if counts[tile] != 4 {
						t.Errorf("expected 4 copies of %v, got %d", tile, counts[tile])
					}
				}
			}
		}
		for rank := 1; rank <= 7; rank++ {
			tile, _ := NewTile(SuitHonor, rank)
			if counts[tile] != 4 {
				t.Errorf("expected 4 copies of honor %v, got %d", tile, counts[tile])
			}
		}
	})

	t.Run("no red fives", func(t *testing.T) {
		rules := Rules{RedFivesMan: 0, RedFivesPin: 0, RedFivesSou: 0}
		wall, err := BuildWall(rules)
		if err != nil {
			t.Fatalf("BuildWall failed: %v", err)
		}

		for _, tile := range wall {
			if tile.IsRed() {
				t.Errorf("found red five in a wall with no red fives: %v", tile)
			}
		}

		counts := make(map[Tile]int)
		for _, tile := range wall {
			counts[tile]++
		}
		if len(counts) != 34 {
			t.Errorf("expected 34 distinct tile kinds, got %d", len(counts))
		}
		for tile, count := range counts {
			if count != 4 {
				t.Errorf("tile %v has %d copies, want 4", tile, count)
			}
		}
	})

	t.Run("max red fives (4 per suit)", func(t *testing.T) {
		rules := Rules{RedFivesMan: 4, RedFivesPin: 4, RedFivesSou: 4}
		wall, err := BuildWall(rules)
		if err != nil {
			t.Fatalf("BuildWall failed: %v", err)
		}

		redCount := 0
		for _, tile := range wall {
			if tile.IsRed() {
				redCount++
			}
		}
		if redCount != 12 {
			t.Errorf("expected 12 red fives, got %d", redCount)
		}

		counts := make(map[Tile]int)
		for _, tile := range wall {
			counts[tile]++
		}
		for _, suit := range []Suit{SuitManzu, SuitPinzu, SuitSouzu} {
			red, _ := NewRedFive(suit)
			normal, _ := NewTile(suit, 5)
			if counts[red] != 4 {
				t.Errorf("expected 4 red fives for %v, got %d", suit, counts[red])
			}
			if counts[normal] != 0 {
				t.Errorf("expected 0 normal fives for %v, got %d", suit, counts[normal])
			}
		}
	})

	t.Run("invalid red five counts", func(t *testing.T) {
		invalidRules := []Rules{
			{RedFivesMan: -1},
			{RedFivesMan: 5},
			{RedFivesPin: -1},
			{RedFivesPin: 5},
			{RedFivesSou: -1},
			{RedFivesSou: 5},
		}
		for _, rules := range invalidRules {
			_, err := BuildWall(rules)
			if err == nil {
				t.Errorf("expected error for rules %+v, got nil", rules)
			}
		}
	})
}

func TestShuffleWall(t *testing.T) {
	rules := DefaultRules()
	wall, _ := BuildWall(rules)

	shuffled, err := ShuffleWall(wall)
	if err != nil {
		t.Fatalf("ShuffleWall failed: %v", err)
	}

	if len(shuffled) != len(wall) {
		t.Errorf("ShuffleWall changed length: %d -> %d", len(wall), len(shuffled))
	}

	// Verify content preservation: same tiles, same counts
	countsOriginal := make(map[Tile]int)
	for _, tile := range wall {
		countsOriginal[tile]++
	}

	countsShuffled := make(map[Tile]int)
	for _, tile := range shuffled {
		countsShuffled[tile]++
	}

	if len(countsOriginal) != len(countsShuffled) {
		t.Errorf("ShuffleWall changed number of distinct tiles: %d -> %d", len(countsOriginal), len(countsShuffled))
	}

	for tile, count := range countsOriginal {
		if countsShuffled[tile] != count {
			t.Errorf("ShuffleWall changed count of tile %v: %d -> %d", tile, count, countsShuffled[tile])
		}
	}

	// Verify it's actually shuffled (very unlikely to be identical)
	if slices.Equal(wall, shuffled) {
		t.Errorf("ShuffleWall returned identical wall (very low probability event)")
	}
}
