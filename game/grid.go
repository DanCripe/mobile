package game

import (
	"fmt"
	"strings"
	"sync"
)

type Grid struct {
	Size   int
	Pixels [][]Pixel
	Off    int
}

type Pixel struct {
	On    bool
	Color Color
}

func (p *Pixel) String() string {
	val := p.Color.String()
	if !p.On {
		val = strings.ToLower(val)
	}
	return fmt.Sprintf("%s", val)
}

type Color interface {
	String() string
	Description() string
	Transform(g *Grid, l Location)
}

func (g *Grid) GetDotColor(x, y int) string {
	return g.Pixels[x][y].String()
}

func (g *Grid) Print() {
	for x := 0; x < g.Size; x++ {
		for y := 0; y < g.Size; y++ {
			fmt.Printf("%s ", g.Pixels[x][y].String())
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}

func (g *Grid) Click(l Location) bool {
	if l.X < 0 || l.Y < 0 || l.X >= g.Size || l.Y >= g.Size {
		return false
	}
	g.Pixels[l.X][l.Y].Color.Transform(g, l)
	if g.Off == 0 {
		return true
	}
	return false
}

func (g *Grid) Toggle(x, y int) {
	if x < 0 || y < 0 || x >= g.Size || y >= g.Size {
		return
	}
	g.Pixels[x][y].On = !g.Pixels[x][y].On
	if g.Pixels[x][y].On {
		g.Off--
	} else {
		g.Off++
	}
}

func (g *Grid) Reset() {
	for x := 0; x < g.Size; x++ {
		for y := 0; y < g.Size; y++ {
			g.Pixels[x][y].On = false
		}
	}
	g.Off = g.Size * g.Size
}

type ColorDef struct {
	OffStr   string
	OnStr    string
	OffPixel func() Pixel
	OnPixel  func() Pixel
}

var o sync.Once
var ColorMap map[string]func() Pixel
var ColorSlice []ColorDef
var ColorIndex map[string]int

func RegisterColor(offStr, onStr string, offPixel, onPixel func() Pixel, offIndex, onIndex int) {
	o.Do(func() {
		ColorMap = make(map[string]func() Pixel)
		ColorIndex = make(map[string]int)
	})
	ColorMap[offStr] = offPixel
	ColorMap[onStr] = onPixel
	ColorIndex[offStr] = offIndex
	ColorIndex[onStr] = onIndex
	ColorSlice = append(ColorSlice, ColorDef{offStr, onStr, offPixel, onPixel})
	if len(ColorSlice)*2 != len(ColorMap) {
		panic(fmt.Sprintf("Duplicate color %s/%s", offStr, onStr))
	}
}

func GetColorIndex() map[string]int {
	return ColorIndex
}
