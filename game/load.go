package game

import (
	"bufio"
	"fmt"
	"os"
)

func Load(file string) (*Grid, error) {
	var g Grid

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(f)

	line, err := reader.ReadBytes(byte('\n'))
	if err != nil {
		return nil, err
	}

	var size int
	fmt.Sscanf(string(line), "%d", &size)
	g.Size = size

	for x := 0; x < size; x++ {
		column := []Pixel{}
		line, err = reader.ReadBytes(byte('\n'))
		if len(line) != size+1 {
			return nil, fmt.Errorf("Incorrect line length at %s:%d.", file, x+2)
		}
		for y := 0; y < g.Size; y++ {
			pixel := NewPixel(line[y : y+1])
			if !pixel.On {
				g.Off++
			}
			column = append(column, pixel)
		}
		g.Pixels = append(g.Pixels, column)
	}

	return &g, nil
}

func NewPixel(p []byte) Pixel {
	pf, ok := ColorMap[string(p)]
	if ok {
		return pf()
	} else {
		return Pixel{On: false, Color: &Red{}}
	}
}
