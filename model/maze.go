package model

import (
	"bytes"
	"fmt"
)

type Maze struct {
	Size  int
	Cells []*Cell
	Start int
	Goal  int
}

type Cell struct {
	NorthWall bool
	EastWall  bool
	SouthWall bool
	WestWall  bool
}

func NewMaze(size int) *Maze {
	m := new(Maze)
	m.Size = size
	m.Start = size * (size - 1) // bottom left
	m.Goal = size - 1           // top right
	m.initCells()
	return m
}

func newCell() *Cell {
	return new(Cell)
}

func (m *Maze) numOfCells() int {
	return m.Size * m.Size
}

func (m *Maze) inFirstRow(i int) bool {
	return i/m.Size == 0
}

func (m *Maze) inLastRow(i int) bool {
	return i/m.Size == m.Size-1
}

func (m *Maze) inFirstCol(i int) bool {
	return i%m.Size == 0
}

func (m *Maze) inLastCol(i int) bool {
	return i%m.Size == m.Size-1
}

func (m *Maze) HasNorthWall(i int) bool {
	cell := m.Cells[i]
	return cell.NorthWall
}

func (m *Maze) HasSouthWall(i int) bool {
	cell := m.Cells[i]
	return cell.SouthWall
}

func (m *Maze) HasEastWall(i int) bool {
	cell := m.Cells[i]
	return cell.EastWall
}

func (m *Maze) HasWestWall(i int) bool {
	cell := m.Cells[i]
	return cell.WestWall
}

func (m *Maze) GetNorthCell(i int) int {
	if m.inFirstRow(i) {
		return i
	}
	return i - m.Size
}

func (m *Maze) GetSouthCell(i int) int {
	if m.inLastRow(i) {
		return i
	}
	return i + m.Size
}

func (m *Maze) GetEastCell(i int) int {
	if m.inLastCol(i) {
		return i
	}
	return i + 1
}

func (m *Maze) GetWestCell(i int) int {
	if m.inFirstCol(i) {
		return i
	}
	return i - 1
}

func (m *Maze) initCells() {
	m.Cells = make([]*Cell, m.numOfCells())
	for i := 0; i < m.numOfCells(); i++ {
		m.Cells[i] = newCell()

	}
	m.buildOuterWall()
}

func (m *Maze) buildOuterWall() {
	for i, c := range m.Cells {
		if m.inFirstRow(i) {
			c.NorthWall = true
		} else if m.inLastRow(i) {
			c.SouthWall = true
		}

		if m.inFirstCol(i) {
			c.WestWall = true
		} else if m.inLastCol(i) {
			c.EastWall = true
		}
	}
}

// For debugging
func (m *Maze) checkCells() {
	for i, c := range m.Cells {
		fmt.Printf("%v: %+v\n", i, c)
	}
}

func (m *Maze) ToString() string {
	var buf bytes.Buffer
	// North walls of 1st row
	for i := 0; i < m.Size; i++ {
		buf.WriteString(" _")
	}
	// West wall of (0, 0)
	buf.WriteString("\n|")

	// Body
	for i, c := range m.Cells {
		if c.SouthWall {
			buf.WriteString("_")
		} else {
			buf.WriteString(" ")
		}
		if c.EastWall {
			buf.WriteString("|")
		} else {
			buf.WriteString(" ")
		}
		if m.inLastCol(i) {
			buf.WriteString("\n")
			if !m.inLastRow(i) {
				buf.WriteString("|")
			}
		}
	}
	return buf.String()
}
