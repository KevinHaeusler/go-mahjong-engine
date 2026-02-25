package engine

import (
	"fmt"
	"strings"
)

type Suit uint8

const (
	SuitManzu Suit = iota // characters / 万子
	SuitPinzu             // dots / 筒子
	SuitSouzu             // bamboo / 索子
	SuitHonor             // winds & dragons / 字牌
)

func (s Suit) String() string {
	switch s {
	case SuitManzu:
		return "m"
	case SuitPinzu:
		return "p"
	case SuitSouzu:
		return "s"
	case SuitHonor:
		return "z"
	default:
		return "?"
	}
}

// Tile represents a **tile kind**, not a physical tile instance.
// For example, "1m", "East", "Red dragon".
//
// Encoding scheme (1 byte):
//
//	bits 0-3: rank (1-9 for numbered suits, 1-7 for honors)
//	bits 4-5: suit (0=man,1=pin,2=sou,3=honor)
type Tile uint8

// NewTile creates a tile from suit and rank.
// For suits man/pin/sou: rank 1-9
// For honors: rank mapping is like Mahjong Soul
//
//	1=East, 2=South, 3=West, 4=North, 5=White, 6=Green, 7=Red
func NewTile(s Suit, rank int) Tile {
	if rank < 1 || rank > 9 {
		panic(fmt.Sprintf("invalid rank %d", rank))
	}
	return Tile((uint8(s) << 4) | uint8(rank))
}

func (t Tile) Suit() Suit {
	return Suit((t >> 4) & 0x03)
}

func (t Tile) Rank() int {
	return int(t & 0x0F)
}

func (t Tile) IsHonor() bool {
	return t.Suit() == SuitHonor
}

func (t Tile) IsNumbered() bool {
	s := t.Suit()
	return s == SuitManzu || s == SuitPinzu || s == SuitSouzu
}

func (t Tile) IsTerminal() bool {
	return t.IsNumbered() && (t.Rank() == 1 || t.Rank() == 9)
}

func (t Tile) IsTerminalOrHonor() bool {
	// Terminals or honors
	return t.IsTerminal() || t.IsHonor()
}

func (t Tile) IsWind() bool {
	return t.Suit() == SuitHonor && t.Rank() >= 1 && t.Rank() <= 4
}

func (t Tile) IsDragon() bool {
	return t.Suit() == SuitHonor && t.Rank() >= 5 && t.Rank() <= 7
}

func (t Tile) String() string {
	suit := t.Suit()
	r := t.Rank()

	if suit == SuitHonor {
		switch r {
		case 1:
			return "E" // East
		case 2:
			return "S" // South
		case 3:
			return "W" // West
		case 4:
			return "N" // North
		case 5:
			return "P" // White (haku/白) - often shown as "P" for "pai"
		case 6:
			return "F" // Green (hatsu/發) - often "F" for "fa"
		case 7:
			return "C" // Red (chun/中)
		default:
			return "?"
		}
	}

	// Numbered suits: "1m", "9p", etc.
	return fmt.Sprintf("%d%s", r, suit.String())
}

// ParseTile parses strings like "1m", "9p", "3s", "E", "S", "W", "N", "P", "F", "C".
func ParseTile(s string) (Tile, error) {
	s = strings.TrimSpace(strings.ToUpper(s))
	if s == "" {
		return 0, fmt.Errorf("empty tile string")
	}

	// Honors as single letters
	if len(s) == 1 {
		switch s {
		case "E":
			return NewTile(SuitHonor, 1), nil
		case "S":
			return NewTile(SuitHonor, 2), nil
		case "W":
			return NewTile(SuitHonor, 3), nil
		case "N":
			return NewTile(SuitHonor, 4), nil
		case "P": // haku / white
			return NewTile(SuitHonor, 5), nil
		case "F": // hatsu / green
			return NewTile(SuitHonor, 6), nil
		case "C": // chun / red
			return NewTile(SuitHonor, 7), nil
		default:
			return 0, fmt.Errorf("unknown honor tile %q", s)
		}
	}

	if len(s) != 2 {
		return 0, fmt.Errorf("invalid tile format %q", s)
	}

	rankChar := s[0]
	suitChar := s[1]

	var rank int
	if rankChar < '1' || rankChar > '9' {
		return 0, fmt.Errorf("invalid rank in %q", s)
	}
	rank = int(rankChar - '0')

	var suit Suit
	switch suitChar {
	case 'M':
		suit = SuitManzu
	case 'P':
		suit = SuitPinzu
	case 'S':
		suit = SuitSouzu
	default:
		return 0, fmt.Errorf("invalid suit in %q", s)
	}

	return NewTile(suit, rank), nil
}
