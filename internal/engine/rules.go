package engine

// Rules contains configurable game rules.
type Rules struct {
	// Number of red 5s in each suit.
	RedFivesMan  int
	RedFivesPin  int
	RedFivesSou  int
	StartingDora int
}

// DefaultRules returns standard Riichi rules (3 akadora) 1 starting Dora.
func DefaultRules() Rules {
	return Rules{
		RedFivesMan:  1,
		RedFivesPin:  1,
		RedFivesSou:  1,
		StartingDora: 1,
	}
}
