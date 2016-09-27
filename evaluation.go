// ghess Evaluation is trying to take all possible moves,
// And somehow pick the best? Yikes
//
package ghess

import (
	"math/rand"
	"unicode"
)

/*
Yikes, Notes:
what are the scoring going to be?
TODO: See how many times each side is attacking a square.

*/

// MoveRandom picks move from lists of valid moves.
// Return an error, such as checkmate or draw.
func (b *Board) MoveRandom(origs, dests []int) error {
	randomMove := rand.Intn(len(origs))
	e := b.Move(origs[randomMove], dests[randomMove])
	if e != nil {
		return e
	}
	return nil
}

// MoveBest finds the best move of all valid moves.
// This method is not operational.
func (b *Board) MoveBest() {
	origs, dests := b.SearchForValid()
	bests := b.EvaluateMoves(origs, dests)
	// get best of bests and return the index
	var best int
	var i int
	for idx, val := range bests {
		if val > best {
			best = val
			i = idx // save index
		} else {
			continue
		}

	}
	b.Move(origs[i], dests[i])
}

// EvaluateMoves scores all valid moves.
// Example of Possible moves:
// [21 21 21 21 43]
// [11 11 12 31 23]
func (b *Board) EvaluateMoves(origs, dests []int) []int {
	var bests []int
	for i := range origs {
		o := origs[i]
		d := dests[i]
		//p := byte(unicode.ToUpper(rune(o)))
		s := b.Evaluate(o, d)
		bests = append(bests, s)
	}
	return bests
}

// Evaluate scores a move based on the piece
// and its destination.
// TODO: Must I acknowledge castling?
func (b *Board) Evaluate(orig, dest int) int {
	var score int
	var trade int
	isWhite := b.toMove == "w"
	piece := byte(unicode.ToUpper(rune(orig)))
	target := byte(unicode.ToUpper(rune(b.board[dest])))
	// If is Capture
	isCapture := target != '.'

	// Default Scores of pieces:
	switch piece {
	case 'P':
		trade = 10
	case 'N':
		trade = 30
	case 'B':
		trade = 30
	case 'R':
		trade = 50
	case 'Q':
		trade = 90
	case 'K':
		trade = 100
	}

	if !isCapture {
		trade = -1
	}

	// Default Scores of pieces:
	switch target {
	case 'P':
		trade -= 10
	case 'N':
		trade -= 30
	case 'B':
		trade -= 30
	case 'R':
		trade -= 50
	case 'Q':
		trade -= 90
	// Doesn't make sense to take king?
	case '.':
		// Redundant, But who cares?
		trade = -1 // No trade
	}

	// trade, Well, it's not necessarily a trade
	// This int tracks whether the attacker is valued
	// more or less than the target. It is a good thing
	// if it is valued less. Right?
	if isCapture {
		score += -trade
	}

	// isCenter and isBorder checks the destinations
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
			// Pawn in center is good
			score += 15
		case piece == 'N':
			// Knight in center is good
			score += 10
		}
	}

	switch piece {
	case 'P':
		// Check if protecting a piece
		// or attacking
		if isWhite {
			pot1 := dest + 9
			pot2 := dest + 11
		} else {
			pot1 := dest - 9
			pot2 := dest - 11

		}
		if pot1 != '.' {
			if (isUpper(pot1) &&
				isWhite) || (!isUpper(pot1) && !isWhite) {
				// protecting
				score += 20
			} else if (isUpper(pot1) &&
				!isWhite) || (!isUpper(pot1) && isWhite) {
				// attacking
				score += 15
			}
		}
		if pot2 != '.' {
			if (isUpper(pot2) &&
				isWhite) || (!isUpper(pot2) && !isWhite) {
				// protecting
				score += 20
			} else if (isUpper(pot2) &&
				!isWhite) || (!isUpper(pot2) && isWhite) {
				// attacking
				score += 15
			}
		}
	case 'N':
		// Check for outpost Knight
		if isWhite {
			if dest > 50 {
				score += 30
			}
		} else {
			if dest < 50 {
				score += 30
			}
		}
	case 'B':
		// Check if on long access Bishop
		if dest%11 == 0 {
			score += 30
		}
	case 'R':
		// If seventh Rank
		if isWhite {
			if dest < 79 || dest > 70 {
				score += 50
			}
		} else {
			if dest < 29 || dest > 20 {
				score += 50
			}
		}
		//remainder := dest - orig
		// Check for open file of Rook?
		//if remainder < 10 && remainder > -10 {
		// Horizontal

		//} else {
		// verticle
		// TODO:
		// check if on open file
		//}
	case 'Q':
		// Fuck lol
	case 'K':
		// idunno fuck

	}

	return score
}
