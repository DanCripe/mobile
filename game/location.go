package game

type Location struct {
	X, Y int
}

func (l *Location) Equal(r Location) bool {
	return l.X == r.X && l.Y == r.Y
}

func (g *Grid) CalculateLocation(x, y, scale int) (Location, bool) {
	size := g.Size
	diameter := scale / size
	var l Location
	if x < 0 || y < 0 {
		return Location{}, false
	}
	l.X = x / diameter
	l.Y = y / diameter
	if l.X >= size || l.Y >= size {
		return Location{}, false
	}
	return l, true
}

func (g *Grid) CalculateDot(X, Y, scale int) (int, int, int, int, string, bool) {
	p := g.Pixels[X][Y]
	size := g.Size
	radius := (scale / size) / 2
	X = X*radius*2 + radius
	Y = Y*radius*2 + radius
	inner := radius * 3 / 4 // inner circle is 3/4 of the size of the dot

	return X, Y, radius, inner, p.Color.String(), p.On
}
