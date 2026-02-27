package main

import (
	"fmt"
	"log"

	"github.com/KevinHaeusler/go-mahjong-engine/internal/engine"
)

func main() {
	debugTiles()
	rules := engine.DefaultRules()
	wall, err := engine.BuildWall(rules)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Unshuffled wall size:", len(wall))

	for i := 0; i < len(wall); i++ {
		fmt.Printf("%d: %s\n", i, wall[i].String())
	}

	fmt.Println("Shuffled wall size:", len(wall))
	wall, err = engine.ShuffleWall(wall)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(wall); i++ {
		fmt.Printf("%d: %s\n", i, wall[i].String())
	}

}

func debugTiles() {
	t, _ := engine.ParseTile("1m")
	fmt.Println(t, t.Suit(), t.Rank(), t.IsTerminal(), t.IsHonor(), t.String())

	h, _ := engine.ParseTile("E")
	fmt.Println(h, h.Suit(), h.Rank(), h.IsHonor(), h.IsWind(), h.String())
}
