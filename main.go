package main

import (
	"fmt"

	"github.com/KevinHaeusler/go-mahjong-engine/internal/engine"
)

func main() {
	debugTiles()
}

func debugTiles() {
	t, _ := engine.ParseTile("1m")
	fmt.Println(t, t.Suit(), t.Rank(), t.IsTerminal(), t.IsHonor(), t.String())

	h, _ := engine.ParseTile("E")
	fmt.Println(h, h.Suit(), h.Rank(), h.IsHonor(), h.IsWind(), h.String())
}
