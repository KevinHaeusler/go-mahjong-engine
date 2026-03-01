package engine

import (
	"fmt"
	"unicode"
)

// ParseHandCompact parses strings like "123m456789s12345z" into []Tile.
//
// Rules:
//   - Digits (0–9) accumulate until a suit letter appears: m/p/s/z
//   - '0' = red five (only for numbered suits)
//   - Honors (z): only 1–7 are valid, not 0 or 8/9
//   - On any invalid input, returns an error (never panics)
func ParseHandCompact(input string) ([]Tile, error) {
	var tiles []Tile
	var digits []rune // buffered rank digits

	flushDigits := func(suit Suit) error {
		if len(digits) == 0 {
			return fmt.Errorf("missing ranks before suit %s", suit.String())
		}

		for _, d := range digits {
			if d < '0' || d > '9' {
				return fmt.Errorf("invalid rank character %q", d)
			}

			rank := int(d - '0')
			var (
				t   Tile
				err error
			)

			if rank == 0 {
				// Red five only for m/p/s
				if suit == SuitHonor {
					return fmt.Errorf("red five (0) not allowed for honors")
				}
				t, err = NewRedFive(suit)
			} else {
				t, err = NewTile(suit, rank)
			}

			if err != nil {
				return err
			}

			tiles = append(tiles, t)
		}

		// Clear buffer
		digits = digits[:0]
		return nil
	}

	for _, r := range input {
		switch {
		case unicode.IsSpace(r):
			continue

		case r >= '0' && r <= '9':
			digits = append(digits, r)

		default:
			// Suit character
			l := unicode.ToLower(r)
			var suit Suit

			switch l {
			case 'm':
				suit = SuitManzu
			case 'p':
				suit = SuitPinzu
			case 's':
				suit = SuitSouzu
			case 'z':
				suit = SuitHonor
			default:
				return nil, fmt.Errorf("unexpected character %q in %q", r, input)
			}

			if err := flushDigits(suit); err != nil {
				return nil, err
			}
		}
	}

	if len(digits) > 0 {
		return nil, fmt.Errorf("dangling ranks at end of input (missing suit)")
	}

	return tiles, nil
}
