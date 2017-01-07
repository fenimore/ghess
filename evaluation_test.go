package ghess

import (
	"testing"
)

func TestEvaluateZero(t *testing.T) {
	game := NewBoard()
	if game.Evaluate() != 0 {
		t.Error("Init Position should be egal")
	}
}
