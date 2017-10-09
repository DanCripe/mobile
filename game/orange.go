package game

type Orange struct {
}

func (c *Orange) String() string {
	return "O"
}

func (c *Orange) Description() string {
	return "Orange pixels toggle diagonals."
}

func (c *Orange) Transform(g *Grid, l Location) {
	for a := 0; a < g.Size; a++ {
		g.Toggle(a, l.Y-l.X+a)
		if a != l.Y {
			g.Toggle(l.X+l.Y-a, a)
		}
	}
}

func OrangePixelOn() Pixel {
	return Pixel{On: true, Color: &Orange{}}
}

func OrangePixelOff() Pixel {
	return Pixel{On: false, Color: &Orange{}}
}

func init() {
	RegisterColor("o", "O", OrangePixelOff, OrangePixelOn, 2, 3)
}
