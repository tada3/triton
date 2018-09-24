package mazegen

import (
	"fmt"
	"testing"
)

func Test_MazeGenerator_Generate(t *testing.T) {
	gen := NewMazeGenerator()
	gen.SetRandSeed(1)
	maze := gen.Generate(6)
	fmt.Println(maze.ToString())
}
