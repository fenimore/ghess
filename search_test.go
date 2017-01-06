package ghess

import (
	"fmt"
	"reflect"
	"testing"
)

// TODO: search for empassant

func TestTension(t *testing.T) {
	game := NewBoard()
	fen := `rnbqkbnr/ppp2ppp/4p3/3p4/4P3/2N5/PPPP1PPP/R1BQKBNR w KQkq - 0 3`
	var err error
	err = game.LoadFen(fen)
	if err != nil {
		fmt.Println(err)
	}

	tension := game.Tension()
	if tension[55] != 0 {
		t.Error("d5 should be equally tense")
	}
}

func TestSearchValid(t *testing.T) {
	var err error
	game := NewBoard()
	fen := "6k1/5p2/7p/1R1r4/P2P1R2/6P1/2r4K/8 w ---- - 0 42"
	_ = game.LoadFen(fen)
	o, d := game.SearchValid()
	exO := []int{21, 21, 21, 43}
	exD := []int{31, 11, 12, 23}
	//exD := []int{11, 12, 31, 23} // old searchvalid function
	if !reflect.DeepEqual(o, exO) || !reflect.DeepEqual(d, exD) {
		fmt.Println(o)
		fmt.Println(d)

		t.Error("1 Search doesn't return the correct/valid moves")
	}

	fen = `rn1q1kbr/ppNNpppp/8/3p4/4P3/8/PPPP1PPP/R1BQKB1R b KQkq - 0 2`
	err = game.LoadFen(fen)
	if err != nil {
		t.Error("Fen error")
	}

	o, d = game.SearchValid()
	exO = []int{85, 87}
	exD = []int{75, 75}
	if !reflect.DeepEqual(o, exO) || !reflect.DeepEqual(d, exD) {
		fmt.Println(o, exO, d, exD)
		t.Error("2 Search doesn't return the correct/valid moves")
	}

	fen = `rn1qkb1r/ppNppppp/8/8/4P3/8/PPPP1PPP/R1BQKBNR b KQkq - 0 2`
	err = game.LoadFen(fen)
	if err != nil {
		t.Error("Fen error")
	}

	o, d = game.SearchValid()
	exO = []int{85}
	exD = []int{76}
	if !reflect.DeepEqual(o, exO) || !reflect.DeepEqual(d, exD) {
		t.Error("3 Search doesn't return the correct/valid moves")
	}
}

func TestSearchValidPawn(t *testing.T) {
	var err error
	game := NewBoard()
	fen := `7k/pppppppp/8/8/8/8/PPPPPPPP/K7 w KQkq - 0 1`
	err = game.LoadFen(fen)
	if err != nil {
		t.Error("Fen error")
	}

	o, d := game.SearchValid()
	exO := []int{18, 21, 21, 22, 22, 23, 23, 24, 24, 25, 25, 26, 26, 27, 27, 28, 28}
	exD := []int{17, 31, 41, 32, 42, 33, 43, 34, 44, 35, 45, 36, 46, 37, 47, 38, 48}
	if !reflect.DeepEqual(o, exO) || !reflect.DeepEqual(d, exD) {
		t.Error("Search doesn't find White Pawn Moves")
	}
}

func TestSearchValidCastle(t *testing.T) {
	var err error
	game := NewBoard()
	// White to move
	fen := `r3k2r/p6p/8/8/8/8/P6P/R3K2R w KQkq - 0 1`
	err = game.LoadFen(fen)
	if err != nil {
		t.Error("Fen error")
	}
	o, d := game.SearchValid()

	exO := []int{11, 11, 14, 14, 14, 14, 14, 14, 14, 18, 18, 18, 21, 21, 28, 28}
	exD := []int{12, 13, 25, 23, 24, 15, 13, 18, 11, 17, 16, 15, 31, 41, 38, 48}
	if !reflect.DeepEqual(o, exO) || !reflect.DeepEqual(d, exD) {
		t.Error("Search doesn't find white castle")
	}

	// For black:
	fen = `r3k2r/p6p/8/8/8/8/P6P/R3K2R b KQkq - 0 1`
	err = game.LoadFen(fen)
	if err != nil {
		t.Error("Fen error")
	}

	o, d = game.SearchValid()
	exO = []int{71, 71, 78, 78, 81, 81, 84, 84, 84, 84, 84, 84, 84, 88, 88, 88}
	exD = []int{61, 51, 68, 58, 82, 83, 85, 73, 74, 75, 83, 88, 81, 87, 86, 85}
	if !reflect.DeepEqual(o, exO) || !reflect.DeepEqual(d, exD) {
		t.Error("Search doesn't find Black castle")
	}
}

func ExampleBoard_SearchValid() {
	game := NewBoard()
	o, d := game.SearchValid()
	fmt.Println(o)
	fmt.Println(d)
	fen := "6k1/5p2/7p/1R1r4/P2P1R2/6P1/2r4K/8 w ---- - 0 42"
	_ = game.LoadFen(fen)
	o, d = game.SearchValid()
	fmt.Println(o)
	fmt.Println(d)

	// Output:
	//[12 12 17 17 21 21 22 22 23 23 24 24 25 25 26 26 27 27 28 28]
	//[33 31 38 36 31 41 32 42 33 43 34 44 35 45 36 46 37 47 38 48]
	//[21 21 21 43]
	//[31 11 12 23]
}
