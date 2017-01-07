package ghess

import (
	"fmt"
	"testing"
)

func TestNewBoard(t *testing.T) {
	game := NewBoard()
	expected := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	if game.Position() != expected {
		t.Error("Unexpected FEN:", game.Position())
	}
}

func TestTurnChange(t *testing.T) {
	game := NewBoard()
	//var err error
	_ = game.Move(24, 44)
	if game.toMove != "b" {
		t.Error("Turn did not change")
	}

	s, _ := MiniMaxPruning(0, 2, GetState(&game))
	_ = game.Move(s.Init[0], s.Init[1])
	if game.toMove != "w" {
		t.Error("Minimax does not change turn?")
	}
}

/**********************************
Examples
***********************************/
func Example_NewBoard() {
	game := NewBoard()
	fmt.Print(game.Position())
	// Output:
	// rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1
}
