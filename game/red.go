package game

type Red struct {
}

func (c *Red) String() string {
	return "R"
}

func (c *Red) Description() string {
	return "Red pixels toggle themselves."
}

func (c *Red) Transform(g *Grid, l Location) {
	g.Toggle(l.X, l.Y)
}

func RedPixelOn() Pixel {
	return Pixel{On: true, Color: &Red{}}
}

func RedPixelOff() Pixel {
	return Pixel{On: false, Color: &Red{}}
}

func init() {
	RegisterColor("r", "R", RedPixelOff, RedPixelOn, 0, 1)
}
