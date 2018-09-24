package game

import "testing"

func Test_GameMaster_PlayMaze(t *testing.T) {

	gm := NewGameMaster()
	gm.SetRandSeed(1)

	gm.StartNew()

	mr, err := gm.Move(SOUTH)
	if err != nil {
		t.Fatalf("Invalid err value: %v\n", err)
	}
	if mr != false {
		t.Fatalf("Invalid move result: %v\n", mr)
	}
	if gm.HitCount() != 1 {
		t.Errorf("Invalid hitCount: %v\n", gm.HitCount())
	}

	mr, err = gm.Move(EAST)
	if err != nil {
		t.Fatalf("Invalid err value: %v\n", err)
	}
	if mr != true {
		t.Fatalf("Invalid move result: %v\n", mr)
	}

	mr, err = gm.Move(EAST)
	if err != nil {
		t.Fatalf("Invalid err value: %v\n", err)
	}
	if mr != true {
		t.Fatalf("Invalid move result: %v\n", mr)
	}

	mr, err = gm.Move(NORTH)
	if err != nil {
		t.Fatalf("Invalid err value: %v\n", err)
	}
	if mr != true {
		t.Fatalf("Invalid move result: %v\n", mr)
	}

	mr, err = gm.Move(WEST)
	if err != nil {
		t.Fatalf("Invalid err value: %v\n", err)
	}
	if mr != false {
		t.Fatalf("Invalid move result: %v\n", mr)
	}

	mr, err = gm.Move(NORTH)
	if err != nil {
		t.Fatalf("Invalid err value: %v\n", err)
	}
	if mr != true {
		t.Fatalf("Invalid move result: %v\n", mr)
	}
	if gm.State() != GOALED {
		t.Fatalf("Invalid state: %v\n", gm.State())
	}

	if gm.MoveCount() != 6 {
		t.Errorf("Invalid tryCount: %v\n", gm.MoveCount())
	}

	gm.Stop()

}

func Test_GameMaster_Locate(t *testing.T) {

	gm := NewGameMaster()
	gm.SetRandSeed(1)

	gm.StartNew()

	loc, err := gm.Locate()
	if err != nil || loc.X != 1 || loc.Y != 1 {
		t.Fatalf("Invalid result: %v, %v\n", loc, err)
	}

	locCount := gm.LocateCount()
	if locCount != 1 {
		t.Fatalf("Wrong count: %v\n", locCount)
	}

	_, err = gm.Move(EAST)
	if err != nil {
		t.Fatalf("Invalid err value: %v\n", err)
	}

	loc, err = gm.Locate()
	if err != nil || loc.X != 2 || loc.Y != 1 {
		t.Fatalf("Invalid result: %v, %v\n", loc, err)
	}

	locCount = gm.LocateCount()
	if locCount != 2 {
		t.Fatalf("Wrong count: %v\n", locCount)
	}
}
