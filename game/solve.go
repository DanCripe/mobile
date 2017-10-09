package game

func (g *Grid) Solvable() bool {
	var coordinates []Location
	for x := 0; x < g.Size; x++ {
		for y := 0; y < g.Size; y++ {
			coordinates = append(coordinates, Location{x, y})
		}
	}
	return g.recurse(coordinates, 0)
}

func (g *Grid) TooEasy() bool {
	green := make(map[int]struct{})
	yellow := make(map[int]struct{})
	for x := 0; x < g.Size; x++ {
		for y := 0; y < g.Size; y++ {
			if g.Pixels[x][y].Color.String() == "G" {
				green[x] = struct{}{}
			}
			if g.Pixels[x][y].Color.String() == "Y" {
				yellow[y] = struct{}{}
			}
		}
	}
	if len(green) == g.Size || len(yellow) == g.Size {
		return true
	}
	return false
}

func (g *Grid) recurse(coords []Location, start int) bool {
	if start > g.Size*g.Size*3/4 {
		return false
	}
	for i := start; i < len(coords); i++ {
		c := coords[i]
		g.Click(c)
		if g.Off == 0 {
			g.Click(c)
			return true
		}
		if g.recurse(coords, i+1) {
			g.Click(c)
			return true
		}
		g.Click(c)
	}
	return false
}
