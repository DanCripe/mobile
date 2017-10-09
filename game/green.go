package game

type Green struct {
}

func (c *Green) String() string {
	return "G"
}

func (c *Green) Description() string {
	return "Green pixels toggle column."
}

func (c *Green) Transform(g *Grid, l Location) {
	for b := 0; b < g.Size; b++ {
		g.Toggle(l.X, b)
	}
}

func GreenPixelOn() Pixel {
	return Pixel{On: true, Color: &Green{}}
}

func GreenPixelOff() Pixel {
	return Pixel{On: false, Color: &Green{}}
}

func init() {
	RegisterColor("g", "G", GreenPixelOff, GreenPixelOn, 6, 7)
}
