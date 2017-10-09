package game

import (
	"fmt"
	"os"
)

func Save(file string, g *Grid) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "%d\n", g.Size)

	for x := 0; x < g.Size; x++ {
		for y := 0; y < g.Size; y++ {
			fmt.Fprintf(f, "%s", g.Pixels[x][y].String())
		}
		fmt.Fprintf(f, "\n")
	}

	return nil
}
