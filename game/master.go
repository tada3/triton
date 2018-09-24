package game

import (
	"errors"
	"fmt"

	"github.com/tada3/triton/mazegen"
	"github.com/tada3/triton/model"
)

type Direction int

const (
	NORTH = iota
	EAST
	SOUTH
	WEST
)

type GameState int

const (
	INIT = iota
	STARTED
	SEARCHING
	GOALED
	STOPPED
	DEAD
)

const (
	DEFAULT_MAZE_SIZE = 3
	MAX_HIT_COUNT     = 3
)

// Representation of location from user's viewpoint
//  bottom left: (1, 1)
//  top right: (size, size)
//  go right: (x, y) -> (x+1, y)
//  go up: (x, y) -> (x, y+1)
type Location struct {
	X int
	Y int
}

func (loc Location) String() string {
	return fmt.Sprintf("%dの%d", loc.X, loc.Y)
}

type GameMaster struct {
	maze            *model.Maze
	state           GameState
	mgen            *mazegen.MazeGenerator
	currentLocation int
	mazeSize        int
	moveCount       int
	locateCount     int
	hitCount        int
}

func NewGameMaster() *GameMaster {
	gm := new(GameMaster)
	gm.mgen = mazegen.NewMazeGenerator()
	gm.mazeSize = DEFAULT_MAZE_SIZE
	return gm
}

func (gm *GameMaster) SetRandSeed(s int64) {
	gm.mgen.SetRandSeed(s)
}

func (gm *GameMaster) State() GameState {
	return gm.state
}

func (gm *GameMaster) MoveCount() int {
	return gm.moveCount
}

func (gm *GameMaster) LocateCount() int {
	return gm.locateCount
}

func (gm *GameMaster) HitCount() int {
	return gm.hitCount
}

func (gm *GameMaster) IsAlive() bool {
	return gm.hitCount < MAX_HIT_COUNT
}

func (gm *GameMaster) StartNew() error {
	if gm.state != INIT && gm.state != STOPPED {
		return fmt.Errorf("Invalid state: %v", gm.state)
	}

	gm.maze = gm.mgen.Generate(gm.mazeSize)
	fmt.Printf("%v\n", gm.maze.ToString())

	gm.currentLocation = gm.maze.Start
	gm.moveCount = 0
	gm.locateCount = 0
	gm.hitCount = 0
	gm.state = STARTED
	return nil
}

func (gm *GameMaster) StartOver() error {
	if gm.state != STOPPED {
		return fmt.Errorf("Invalid state: %v", gm.state)
	}

	if gm.maze == nil {
		return errors.New("Maze not exist.")
	}

	gm.currentLocation = gm.maze.Start
	gm.moveCount = 0
	gm.locateCount = 0
	gm.hitCount = 0
	gm.state = STARTED
	return nil
}

func (gm *GameMaster) Stop() {
	if gm.state == INIT || gm.state == STOPPED {
		return
	}
	gm.state = STOPPED
}

func (gm *GameMaster) Move(dir Direction) (bool, error) {
	if gm.state != STARTED && gm.state != SEARCHING {
		return false, fmt.Errorf("Invalid state: %v", gm.state)
	}

	fmt.Printf("Move %d, currentLocation=%d\n", dir, gm.currentLocation)

	gm.moveCount++

	fmt.Printf("Move moveCount=%v\n", gm.moveCount)

	var hasWall bool
	if dir == NORTH {
		hasWall = gm.maze.HasNorthWall(gm.currentLocation)
	} else if dir == EAST {
		hasWall = gm.maze.HasEastWall(gm.currentLocation)
	} else if dir == SOUTH {
		hasWall = gm.maze.HasSouthWall(gm.currentLocation)
	} else if dir == WEST {
		hasWall = gm.maze.HasWestWall(gm.currentLocation)
	}

	fmt.Printf("Move hasWall=%v\n", hasWall)

	if hasWall {
		gm.hitCount++
		if !gm.IsAlive() {
			gm.state = DEAD
		} else {
			gm.state = SEARCHING
		}

		fmt.Printf("Move hasWall state=%v, hitCount=%v\n", gm.state, gm.hitCount)

		return false, nil
	}

	if dir == NORTH {
		gm.currentLocation = gm.maze.GetNorthCell(gm.currentLocation)
	} else if dir == EAST {
		gm.currentLocation = gm.maze.GetEastCell(gm.currentLocation)
	} else if dir == SOUTH {
		gm.currentLocation = gm.maze.GetSouthCell(gm.currentLocation)
	} else if dir == WEST {
		gm.currentLocation = gm.maze.GetWestCell(gm.currentLocation)
	}

	fmt.Printf("Move after move currentLocation=%d\n", gm.currentLocation)

	if gm.currentLocation == gm.maze.Goal {
		gm.state = GOALED
	} else {
		gm.state = SEARCHING
	}
	return true, nil
}

func (gm *GameMaster) GetSize() (int, error) {
	if gm.maze == nil {
		return 0, errors.New("Maze not exist.")
	}
	return gm.maze.Size, nil
}

func (gm *GameMaster) GetStart() (Location, error) {
	if gm.maze == nil {
		return Location{}, errors.New("Maze not exist.")
	}
	return gm.getLocation(gm.maze.Start), nil
}

func (gm *GameMaster) GetGoal() (Location, error) {
	if gm.maze == nil {
		return Location{}, errors.New("Maze not exist.")
	}
	return gm.getLocation(gm.maze.Goal), nil
}

func (gm *GameMaster) Locate() (Location, error) {
	if gm.maze == nil {
		return Location{}, errors.New("Maze not exist.")
	}
	gm.locateCount++
	return gm.getLocation(gm.currentLocation), nil
}

func (gm *GameMaster) GetCurrentLocation() (Location, error) {
	if gm.maze == nil {
		return Location{}, errors.New("Maze not exist.")
	}
	return gm.getLocation(gm.currentLocation), nil
}

func (gm *GameMaster) getLocation(i int) Location {
	size := gm.maze.Size
	return Location{i%size + 1, size - i/size}
}
