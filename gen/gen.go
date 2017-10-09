package main

import (
	"fmt"

	"github.com/dancripe/tmo/game"
)

func main() {
	g, c := game.Random(3)
	fmt.Printf("Count: %d\n", c)
	g.Print()
	/*
		g, c := game.Random(4)
		fmt.Printf("Count: %d\n", c)
		g.Print()
		g, c = game.Random(4)
		fmt.Printf("Count: %d\n", c)
		g.Print()
		g, c = game.Random(4)
		fmt.Printf("Count: %d\n", c)
		g.Print()
		g, c = game.Random(4)
		fmt.Printf("Count: %d\n", c)
		g.Print()
		g, c = game.Random(4)
		fmt.Printf("Count: %d\n", c)
		g.Print()
		g, c = game.Random(5)
		fmt.Printf("Count: %d\n", c)
		g.Print()
	*/
}
