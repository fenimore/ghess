package ghess

import (
    "unicode"
    "math/rand"
)

// MoveRandom, pick move from lists of valid moves.
// Return an error, such as checkmate or draw.
func (b *Board) MoveRandom(origs, dests []int) error {
	randomMove := rand.Intn(len(origs))
	e := b.Move(origs[randomMove], dests[randomMove])
	if e != nil {
		return e
	}
	return nil
}

func (b *Board) MoveBest() {
	origs, dests := b.SearchForValid()
	bests := b.EvaluateMoves(origs, dests)
	// get best of bests and return the index
	var best int
	var i int
	for idx, val := range bests {
		if val > best {
			i = idx // save index
		} else {
			continue
		}

	}
	b.Move(origs[i], dests[i])
}

// EvaluateMoves() scores all valid moves.
func (b *Board) EvaluateMoves(origs, dests []int) []int {
	var bests []int
	for i, _ := range origs {
		o := b.board[origs[i]]
		d := dests[i]
		p := byte(unicode.ToUpper(rune(o)))
		s := b.Evaluate(p, d)
		bests = append(bests, s)
	}
	return bests
}

// Evaluate() scores a move based on the piece
// and its destination.
// TODO: Must I acknowledge castling?
func (b *Board) Evaluate(piece byte, dest int) int {
	var score int
	isCenter := dest == 55 || dest == 44 ||
		dest == 45 || dest == 54
	isBorder := (dest-1)%10 == 0 || (dest+2)%10 == 0
	if isBorder {
		switch {
		case piece == 'N':
			score -= 10
		}
	}
	if isCenter {
		switch {
		case piece == 'P':
			score += 50
		case piece == 'N':
			score += 25
		}
	}

	// Check for open file of Rook?
	// Check if on long access Bishop
	// Check if puts in check?
	// Should I be effecting the move?
	// eg. possible.UpdateBoard()????

	// If is Capture
	if b.board[dest] != '.' {
		score += 100
	}

	return score
}
