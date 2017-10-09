package game

type Blue struct {
}

func (c *Blue) String() string {
	return "B"
}

func (c *Blue) Description() string {
	return "Blue pixels toggle neighbors."
}

func (c *Blue) Transform(g *Grid, l Location) {
	for a := -1; a <= 1; a++ {
		for b := -1; b <= 1; b++ {
			if a == 0 && b == 0 {
				continue
			}
			g.Toggle(l.X+a, l.Y+b)
		}
	}
}

func BluePixelOn() Pixel {
	return Pixel{On: true, Color: &Blue{}}
}

func BluePixelOff() Pixel {
	return Pixel{On: false, Color: &Blue{}}
}

func init() {
	RegisterColor("b", "B", BluePixelOff, BluePixelOn, 8, 9)
}
