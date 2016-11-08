// Package ghess Evaluation scores a certain position.
package ghess

import (
	"errors"
	"math/rand"
	"time"
)

/*
State Functions ##################################3
*/

/*
 Evaluation is HERE: #######################################

### Ideas for Evaluation:

- tension
- mobility (possible moves)
- Double/isolated pawns
- Open file Rook
- Outpost knight
- Center control


*/

// pawnProtect returns true if piece is protected by a pawn.
func (b *Board) pawnProtect(dest int, isWhite bool) bool {
	var pot1, pot2 int
	var pawn byte
	if !isWhite {
		pot1 = dest + 9
		pot2 = dest + 11
		pawn = 'p'
	} else {
		pot1 = dest - 9
		pot2 = dest - 11
		pawn = 'P'
	}

	switch pawn {
	case b.board[pot1]:
		return true
	case b.board[pot2]:
		return true
	default:
		return false
	}
}

// pawnThreat returns true if square is attacked
// by enemy pawn. According to turn
func (b *Board) pawnThreat(dest int, isWhite bool) bool {
	//isWhite := b.toMove == "w"

	var pot1, pot2 int
	var pawn byte
	if isWhite {
		pot1 = dest + 9
		pot2 = dest + 11
		pawn = 'p'
	} else {
		pot1 = dest - 9
		pot2 = dest - 11
		pawn = 'P'
	}

	switch pawn {
	case b.board[pot1]:
		return true
	case b.board[pot2]:
		return true
	default:
		return false
	}
}

// threatenQueen checks if the piece threatens the enemy queen.
func (b *Board) queenThreaten(pos int, piece byte, isWhite bool) bool {
	// find queen
	var queen int
	var e error
	for idx, val := range b.board {
		// Only look for 64 squares
		if idx%10 == 0 || (idx+1)%10 == 0 ||
			idx > 88 || idx < 11 {
			continue
		}
		if isWhite && val == 'q' {
			queen = idx
		} else if !isWhite && val == 'Q' {
			queen = idx
		}
	}
	switch piece {
	case 'P': // Cause pawns are weird
		return false
	case 'N':
		e = b.validKnight(pos, queen)
	case 'B':
		e = b.validBishop(pos, queen)
	case 'R':
		e = b.validRook(pos, queen)
	case 'Q':
		return false
	case 'K':
		e = b.validKing(pos, queen, false)
	}
	if e == nil {
		return true
	}
	return false
}

// TODO: implement, pawnThreatensPiece

// See chess programming wiki:
// http://chessprogramming.wikispaces.com/Simplified+evaluation+function

// Evaluate returns score based on position.
// When evaluating individual pieces, the boolean to pass
// in does not mean WHOSE turn it is but rather who owns the piece.
func (b *Board) Evaluate() int {
	// For position, if piece,
	var score int

	if b.Checkmate {
		if b.Score == "0-1" {
			score -= 1000000
		} else if b.Score == "1-0" {
			score += 1000000
		}
	} else if b.Draw {
		// Discourage the computer from draw,
		// But it should be better than checkmate
		if b.toMove == "w" {
			score -= 50000
		} else {
			score += 50000
		}
	}

	for idx, val := range b.board {
		// only look at 64 squares:
		if idx%10 == 0 || (idx+1)%10 == 0 || idx > 88 ||
			idx < 11 {
			continue
		} else if val == '.' {
			continue
		}
		isWhitePiece := b.isUpper(idx)
		//piece := b.board[idx]
		if isWhitePiece {
			score += matMap[val]
		} else {
			score -= matMap[val]
		}
		switch val {
		case 'P', 'p':
			if isWhitePiece {
				score += whitePawnMap[idx]
			} else {
				score -= blackPawnMap[idx]
			}
			//score += b.evalPawn(idx, isWhitePiece)
		case 'N', 'n':
			if isWhitePiece {
				score += whiteKnightMap[idx]
			} else {
				score -= blackKnightMap[idx]
			}
			//score += b.evalKnight(idx, isWhitePiece)
		case 'B', 'b':
			if isWhitePiece {
				score += whiteBishopMap[idx]
			} else {
				score -= blackBishopMap[idx]
			}
			//score += b.evalBishop(idx, isWhitePiece)
		case 'R', 'r':
			if isWhitePiece {
				//score += whiteRookMap[idx]
			} else {
				//score -= blackRookMap[idx]
			}
			//score += b.evalRook(idx, isWhitePiece)
		//case 'Q', 'q':
		//score += b.evalQueen(idx, isWhitePiece)
		//case 'K', 'k':
		//score += b.evalKing(idx, isWhitePiece)
		default:
			//wtf default?
			score += 0
		}
	}
	return score

}

/*
 LULZ EVALUATION ISN'T NECESSARY!!!1!#######################
*/

// MoveRandom picks move from lists of valid moves.
// Return an error, such as checkmate or draw.
func (b *Board) MoveRandom(origs, dests []int) error {
	if len(origs) < 1 {
		return errors.New("There are no valid moves left")
	}
	rand.Seed(time.Now().UTC().UnixNano())
	randomMove := rand.Intn(len(origs))
	e := b.Move(origs[randomMove], dests[randomMove])
	if e != nil {
		return e
	}
	return nil
}
