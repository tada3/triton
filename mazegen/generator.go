package mazegen

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"
	"unsafe"

	"github.com/tada3/triton/model"
)

// MazeGenerator generate a maze with Kruscal's method.
type MazeGenerator struct {
	size   int
	cells  []int
	forest [][]int
	walls  []*wall
	rnd    *rand.Rand
}

type wall struct {
	removed bool
	a       int
	b       int
}

// NewMazeGenerator create a new MazeGenerator
func NewMazeGenerator() *MazeGenerator {
	g := new(MazeGenerator)
	seed := time.Now().UnixNano()
	fmt.Printf("seed=%v\n", seed)
	g.rnd = rand.New(rand.NewSource(seed))
	return g
}

// For testing only
func (g *MazeGenerator) SetRandSeed(s int64) {
	g.rnd = rand.New(rand.NewSource(s))
}

func (g *MazeGenerator) numOfCells() int {
	return g.size * g.size
}

func (g *MazeGenerator) inFastRow(i int) bool {
	return i/g.size == 0
}

func (g *MazeGenerator) inLastRow(i int) bool {
	return i/g.size == g.size-1
}

func (g *MazeGenerator) inFastCol(i int) bool {
	return i%g.size == 0
}

func (g *MazeGenerator) inLastCol(i int) bool {
	return i%g.size == g.size-1
}

func (g *MazeGenerator) addSouthWall(cell int) {
	w := newWall(cell, cell+g.size)
	g.walls = append(g.walls, w)
}

func (g *MazeGenerator) addEastWall(cell int) {
	w := newWall(cell, cell+1)
	g.walls = append(g.walls, w)
}

func (g *MazeGenerator) Generate(size int) *model.Maze {
	g.initCells(size)

	g.initSetsAndWalls()

	g.build()

	m := g.genMaze()

	return m
}

func (g *MazeGenerator) initCells(size int) {
	g.size = size
	g.cells = make([]int, g.numOfCells())
	for i := range g.cells {
		g.cells[i] = i
	}
}

func (g *MazeGenerator) initSetsAndWalls() {
	for _, c := range g.cells {
		// sets
		x := []int{c}
		g.forest = append(g.forest, x)

		// walls
		if !g.inLastRow(c) {
			g.addSouthWall(c)
		}
		if !g.inLastCol(c) {
			g.addEastWall(c)
		}
	}
}

func (g *MazeGenerator) build() {
	wallIdxs := g.rnd.Perm(len(g.walls))
	for _, i := range wallIdxs {
		wall := g.walls[i]
		aSetIdx := -1
		bSetIdx := -1
		for j, set := range g.forest {
			if exists(set, wall.a) {
				aSetIdx = j
			}
			if exists(set, wall.b) {
				bSetIdx = j
			}
			if aSetIdx >= 0 && bSetIdx >= 0 {
				break
			}
		}
		if aSetIdx < 0 || bSetIdx < 0 {
			fmt.Printf("ERROR! Cannot set setIdx: %v, %v\n", aSetIdx, bSetIdx)
			panic("Invalid state!")
		}
		if aSetIdx != bSetIdx {
			aSet := g.forest[aSetIdx]
			bSet := g.forest[bSetIdx]
			if aSetIdx < bSetIdx {
				g.forest = removeElement(g.forest, aSetIdx)
				g.forest = removeElement(g.forest, bSetIdx-1)
			} else {
				g.forest = removeElement(g.forest, bSetIdx)
				g.forest = removeElement(g.forest, aSetIdx-1)
			}
			g.forest = append(g.forest, union(aSet, bSet))

			wall.removed = true
		}
	}
}

func (g *MazeGenerator) genMaze() *model.Maze {

	m := model.NewMaze(g.size)

	for _, wall := range g.walls {
		if !wall.removed {
			if wall.isHorizontal() {
				m.Cells[wall.a].SouthWall = true
				m.Cells[wall.b].NorthWall = true
			} else {
				m.Cells[wall.a].EastWall = true
				m.Cells[wall.b].WestWall = true
			}
		}
	}

	return m
}

func newWall(a int, b int) *wall {
	w := new(wall)
	w.a = a
	w.b = b
	return w
}

// Added as south wall
func (w *wall) isHorizontal() bool {
	return (w.b - w.a) > 1
}

func exists(s []int, w int) bool {
	for _, e := range s {
		if e == w {
			return true
		}
	}
	return false
}

func isIdentical(s1 []int, s2 []int) bool {
	sh1 := (*reflect.SliceHeader)(unsafe.Pointer(&s1))
	p1 := unsafe.Pointer(sh1.Data)
	sh2 := (*reflect.SliceHeader)(unsafe.Pointer(&s2))
	p2 := unsafe.Pointer(sh2.Data)
	return p1 == p2
}

func union(s1 []int, s2 []int) []int {
	return append(s1, s2...)
}

func removeElement(s [][]int, i int) [][]int {
	// Always i < len(s)
	return append(s[:i], s[i+1:]...)
}
