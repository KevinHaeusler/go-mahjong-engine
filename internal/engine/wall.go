package engine

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	copiesPerTileKind = 4
	totalTileKinds    = 34                                 // 9m+9p+9s+7 honors
	totalTiles        = copiesPerTileKind * totalTileKinds // 136
)

// BuildWall creates a full 136-tile wall based on the rules.
// - Uses red 5s (0m / 0p / 0s) according to RedFives* counts.
// - Returns tiles in a deterministic order.
func BuildWall(rules Rules) ([]Tile, error) {
	// Basic validation of red five counts
	if rules.RedFivesMan < 0 || rules.RedFivesMan > copiesPerTileKind {
		return nil, fmt.Errorf("RedFivesMan must be between 0 and %d", copiesPerTileKind)
	}
	if rules.RedFivesPin < 0 || rules.RedFivesPin > copiesPerTileKind {
		return nil, fmt.Errorf("RedFivesPin must be between 0 and %d", copiesPerTileKind)
	}
	if rules.RedFivesSou < 0 || rules.RedFivesSou > copiesPerTileKind {
		return nil, fmt.Errorf("RedFivesSou must be between 0 and %d", copiesPerTileKind)
	}

	wall := make([]Tile, 0, totalTiles)

	// Helper to add 4 copies of a tile kind
	addCopies := func(t Tile) {
		for i := 0; i < copiesPerTileKind; i++ {
			wall = append(wall, t)
		}
	}

	// Numbered suits with red fives
	type suitRed struct {
		suit     Suit
		redCount int
	}
	suits := []suitRed{
		{SuitManzu, rules.RedFivesMan},
		{SuitPinzu, rules.RedFivesPin},
		{SuitSouzu, rules.RedFivesSou},
	}

	for _, sr := range suits {
		for rank := 1; rank <= 9; rank++ {
			if rank == 5 {
				// Special handling for 5's: mix of normal 5 and red 5
				normalFive, err := NewTile(sr.suit, 5)
				if err != nil {
					return nil, err
				}
				redFive, err := NewRedFive(sr.suit)
				if err != nil {
					return nil, err
				}

				// e.g. if redCount = 1 -> [red, normal, normal, normal]
				for i := 0; i < copiesPerTileKind; i++ {
					if i < sr.redCount {
						wall = append(wall, redFive)
					} else {
						wall = append(wall, normalFive)
					}
				}
				continue
			}

			t, err := NewTile(sr.suit, rank)
			if err != nil {
				return nil, err
			}
			addCopies(t)
		}
	}

	// Honors: 4 winds + 3 dragons = 7 kinds, 4 copies each.
	// rank mapping:
	//   1=E, 2=S, 3=W, 4=N, 5=G, 6=R, 7=Wh
	for rank := 1; rank <= 7; rank++ {
		t, err := NewTile(SuitHonor, rank)
		if err != nil {
			return nil, err
		}
		addCopies(t)
	}

	if len(wall) != totalTiles {
		return nil, fmt.Errorf("internal error: expected %d tiles, got %d", totalTiles, len(wall))
	}

	return wall, nil
}

// ShuffleWall returns a shuffled copy of the given wall.
func ShuffleWall(wall []Tile) ([]Tile, error) {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	shuffled := make([]Tile, len(wall))
	copy(shuffled, wall)

	rng.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled, nil
}
