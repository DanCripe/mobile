package game

import (
	"math/rand"
)

func Random(size int, index int) (*Grid, int) {
	var c int
	rand.Seed(int64(index))
	g := new(Grid)
	g.Size = size
	g.Off = size * size
	for {
		c++
		g.random()
		if g.Solvable() && !g.TooEasy() {
			return g, c
		}
	}
}

func (g *Grid) random() {
	g.Pixels = [][]Pixel{}
	r := len(ColorSlice)
	for x := 0; x < g.Size; x++ {
		g.Pixels = append(g.Pixels, []Pixel{})
		for y := 0; y < g.Size; y++ {
			i := rand.Intn(r)
			g.Pixels[x] = append(g.Pixels[x], ColorSlice[i].OffPixel())
		}
	}
}
