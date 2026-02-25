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
//	bit 6: isDora
//	bit 7: isUra
type Tile uint8

const (
	maskRank uint8 = 0x0F
	maskSuit uint8 = 0x30
	bitDora  uint8 = 1 << 6
	bitUra   uint8 = 1 << 7
)

// NewTile creates a tile from suit and rank.
// For suits man/pin/sou: rank 1-9
// For honors: rank mapping is like Mahjong Soul
//
//	1=East, 2=South, 3=West, 4=North, 5=White, 6=Green, 7=Red
func NewTile(s Suit, rank int) (Tile, error) {
	if rank < 0 || rank > 9 {
		return 0, fmt.Errorf("invalid rank %d", rank)
	}
	if s == SuitHonor && (rank < 1 || rank > 7) {
		return 0, fmt.Errorf("invalid honor rank %d, expected 1–7", rank)
	}
	return Tile((uint8(s) << 4) | uint8(rank)), nil
}

// NewRedFive creates a red 5 tile for the given suit.
func NewRedFive(s Suit) (Tile, error) {
	if s == SuitHonor {
		return 0, fmt.Errorf("no red five for honor tiles")
	}
	return Tile((uint8(s) << 4) | 0), nil // rank 0 = red five
}

func (t Tile) Suit() Suit {
	return Suit((uint8(t) & maskSuit) >> 4)
}

func (t Tile) Rank() int {
	return int(uint8(t) & maskRank)
}

func (t Tile) IsDora() bool {
	return uint8(t)&bitDora != 0
}

func (t Tile) IsUra() bool {
	return uint8(t)&bitUra != 0
}

func (t Tile) IsRed() bool {
	return t.IsNumbered() && t.Rank() == 0
}

func (t Tile) SetDora(isDora bool) Tile {
	if isDora {
		return Tile(uint8(t) | bitDora)
	}
	return Tile(uint8(t) & ^bitDora)
}

func (t Tile) SetUra(isUra bool) Tile {
	if isUra {
		return Tile(uint8(t) | bitUra)
	}
	return Tile(uint8(t) & ^bitUra)
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
	s := t.Suit()
	r := t.Rank()

	if s != SuitHonor {
		if r == 0 {
			return "0" + s.String() // red five
		}
		return fmt.Sprintf("%d%s", r, s.String())
	}

	// Honor tiles
	switch r {
	case 1:
		return "E"
	case 2:
		return "S"
	case 3:
		return "W"
	case 4:
		return "N"
	case 5:
		return "G" // green Dragon
	case 6:
		return "R" // red Dragon
	case 7:
		return "Wh" // white Dragon
	}

	return "?"
}

// ParseTile accepts:
// Numbered: 1m, 0p, 9s
// Winds:   E, S, W, N
// Dragons: G, R, Wh
func ParseTile(s string) (Tile, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty tile")
	}

	up := strings.ToUpper(s)

	switch up {
	case "E":
		return NewTile(SuitHonor, 1)
	case "S":
		return NewTile(SuitHonor, 2)
	case "W":
		return NewTile(SuitHonor, 3)
	case "N":
		return NewTile(SuitHonor, 4)
	case "G":
		return NewTile(SuitHonor, 5)
	case "R":
		return NewTile(SuitHonor, 6)
	case "WH":
		return NewTile(SuitHonor, 7)

	}

	// Numbered suits (0m, 5p, 9s)
	if len(up) != 2 {
		return 0, fmt.Errorf("invalid tile format: %q", s)
	}

	rankChar := up[0]
	suitChar := up[1]

	// Rank parsing
	var rank int
	if rankChar == '0' {
		rank = 0
	} else if rankChar >= '1' && rankChar <= '9' {
		rank = int(rankChar - '0')
	} else {
		return 0, fmt.Errorf("invalid rank in %q", s)
	}

	// Suit parsing
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

	return NewTile(suit, rank)
}
