package model

import (
	"fmt"
	"testing"
)

func Test_ToString(t *testing.T) {
	maze := NewMaze(2)
	fmt.Println(maze.ToString())
}
