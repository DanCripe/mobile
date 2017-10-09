package game

type Yellow struct {
}

func (c *Yellow) String() string {
	return "Y"
}

func (c *Yellow) Description() string {
	return "Yellow pixels toggle row."
}

func (c *Yellow) Transform(g *Grid, l Location) {
	for b := 0; b < g.Size; b++ {
		g.Toggle(b, l.Y)
	}
}

func YellowPixelOn() Pixel {
	return Pixel{On: true, Color: &Yellow{}}
}

func YellowPixelOff() Pixel {
	return Pixel{On: false, Color: &Yellow{}}
}

func init() {
	RegisterColor("y", "Y", YellowPixelOff, YellowPixelOn, 4, 5)
}
