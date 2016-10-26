// Package ghess Evaluation scores a certain position.
package ghess

import (
	"bytes"
	"errors"
	"math/rand"
	"strings"
	"time"
	"unicode"
)

/*
State Functions ##################################3
*/

// GetState turns a Board into a copy and it's state.
// The Init value is nil.
func GetState(b *Board) State {
	c := *b                        // dereference the pointer
	boardCopy := make([]byte, 120) // []bytes are slices
	castleCopy := make([]byte, 4)
	copy(boardCopy, b.board)
	copy(castleCopy, b.castle)
	c.board = boardCopy
	c.castle = castleCopy
	s := State{board: &c, eval: c.Evaluate()}
	return s
}

// TryState takes in a *Board and valid move and returns
// a State struct.
func TryState(b *Board, o, d int) (State, error) {
	state := State{}
	possible := CopyBoard(b)
	err := possible.Move(o, d)
	if err != nil {
		return state, err
	}
	state.board = possible
	state.eval = possible.Evaluate()
	return state, nil
}

// GetPossibleStates returns a slice of State structs
// Each with a score and the move that got there.
func GetPossibleStates(state State) (States, error) {
	states := make(States, 0)
	origs, dests := state.board.SearchValid()
	for i := 0; i < len(origs); i++ {
		s, err := TryState(state.board, origs[i], dests[i])
		if err != nil {
			return states, err
		}
		if state.Init[0] == 0 {
			s.Init[0], s.Init[1] = origs[i], dests[i]
		} else {
			s.Init[0], s.Init[1] = state.Init[0], state.Init[1]
		}
		s.isMax = state.isMax // Basically is White
		states = append(states, s)
	}
	return states, nil
}

// DictionaryAttack looks up common openings
// for less stupid opening moves.
func DictionaryAttack(s State) (State, error) {
	position := s.board.Position()
	// I don't need castling empassant or move number
	posits := strings.Split(position, " ")
	key := posits[0] + " " + posits[1]
	// Check if opening exists
	if val, ok := dict[key]; ok {
		state := State{Init: val}
		return state, nil
	}

	return s, errors.New("No Dictionary Attack Found")
}

/*
 Evaluation is HERE: #######################################
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

// Evaluate returns score based on position.
// When evaluating individual pieces, the boolean to pass
// in does not mean WHOSE turn it is but rather who owns the piece.
func (b *Board) Evaluate() int {
	// For position, if piece,
	var score int
	tension := b.Tension()
	var material int

	/*	// Find King
		for idx, _ := range b.board {
			// Only look for 64 squares
			if idx%10 == 0 || (idx+1)%10 == 0 ||
				idx > 88 || idx < 11 {
				continue
			}
		}
	*/

	if b.Check && b.toMove == "w" {
		score -= 200
	} else if b.Check && b.toMove == "b" {
		score += 200
	}

	if b.Checkmate {
		if b.score == "0-1" {
			score -= 9000
		} else if b.score == "1-0" {
			score += 9000
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
		piece := []byte(bytes.ToUpper(b.board[idx : idx+1]))[0]
		// Tensions values are negative if the majority of attackers
		// are black (etc). So if a piece moves to a 'controlled' square
		// I'm giving 5 times the over protection.
		if isWhitePiece {
			if tension[idx] < 1 {
				score -= 20
				// Controlling a square doesn't mean much
				// if it is attacked by pawns
			} else if !b.pawnThreat(idx, isWhitePiece) {
				score += (3 * tension[idx])
			}
		} else {
			if tension[idx] > -1 {
				score += 20
				// Controlling a square doesn't mean much
				// if it is attacked by pawns
			} else if !b.pawnThreat(idx, !isWhitePiece) {
				score -= (-3 * tension[idx])
			}
		}
		switch piece {
		case 'P':
			score += b.evalPawn(idx, isWhitePiece)
			if isWhitePiece {
				material += 10
			} else {
				material -= 10
			}
		case 'N':
			score += b.evalKnight(idx, isWhitePiece)
			if isWhitePiece {
				material += 30
			} else {
				material -= 30
			}
		case 'B':
			score += b.evalBishop(idx, isWhitePiece)
			if isWhitePiece {
				material += 30
			} else {
				material -= 30
			}
		case 'R':
			score += b.evalRook(idx, isWhitePiece)
			if isWhitePiece {
				material += 50
			} else {
				material -= 50
			}
		case 'Q':
			score += b.evalQueen(idx, isWhitePiece)
			if isWhitePiece {
				material += 100
			} else {
				material -= 100
			}
		case 'K':
			score += b.evalKing(idx, isWhitePiece)
		default:
			//wtf default?
			score += 0
		}
		// If threatens queen, maybe not so snazzy
		/*if b.queenThreaten(idx, piece, isWhitePiece) {
			if isWhitePiece {
				score += 200
			} else {
				score -= 200
			}
		}*/
	}
	// Take the material advantage
	// and multiply by two for greater weight.
	score += (material * 3)
	return score

}

/*
Evaluations:
- Inverted score for Black
- Typically 20 for Good position
- Piece value * 10 for piece itself
- 10 for support position
- minus piece value for pawnThreatened
- 50 for Awesome position
*/

// evalPawn returns a score for pawn position. It'll be
// negative for black and positive for white.
// Currently check for center, support, and seventh rank pawns.
// TODO: is attacking a piece.
func (b *Board) evalPawn(pos int, isWhite bool) int {
	var score int
	var isCenter bool
	score += 10 // Score for simply having a pawn
	// if in center
	if pos == 44 || pos == 45 || pos == 54 || pos == 55 {
		isCenter = true
		score += 20
	} else if pos == 51 || pos == 41 || pos == 48 || pos == 58 {
		score -= 20
	}
	if b.pawnProtect(pos, isWhite) {
		if !b.pawnThreat(pos, isWhite) {
			score += 2
			if isCenter {
				score += 6
			}
		}
	} else if b.pawnThreat(pos, isWhite) {
		score -= 11
	}

	// Invert for black
	// Position values
	if isWhite {
		switch {
		case pos == 31 || pos == 38:
			score += 10
		case pos > 70:
			score += 50 // seventh rank
		case pos == 46 || pos == 43 || pos == 35 || pos == 34:
			// support pawns
			score += 6
		}
	} else {
		score = -score
		switch {
		case pos == 61 || pos == 68:
			score -= 10
		case pos < 20:
			score -= 50 // seventh rank
		case pos == 56 || pos == 53 || pos == 65 || pos == 64:
			score -= 6
		}
	}

	return score
}

// evalKnight evaluates for knight position.
func (b *Board) evalKnight(pos int, isWhite bool) int {
	var score int
	score += 30 // just for being a knight
	if b.pawnThreat(pos, isWhite) {
		score -= 30 // attacked by opponent
	}
	// is on edge, yuck
	if (pos-1)%10 == 0 || (pos+2)%10 == 0 {
		score -= 20
	}
	// The score is inverted for Black
	if isWhite {
		if pos == 33 || pos == 36 {
			score += 20
		} else if pos > 48 {
			score += 30
		}
		// Outpost Knigth
		if pos > 38 && b.pawnProtect(pos, isWhite) {
			score += 3
		}
	} else {
		score = -score
		if pos == 63 || pos == 66 {
			score -= 20
		} else if pos < 58 {
			score -= 30
		}
		// Outpost Knigth
		if pos < 68 && b.pawnProtect(pos, isWhite) {
			score -= 3
		}
	}
	return score
}

// evalBishop evaluates bishop position.
func (b *Board) evalBishop(pos int, isWhite bool) int {
	var score int
	score += 30 // just for being a knight
	if b.pawnThreat(pos, isWhite) {
		score -= 30 // attacked by opponent
	}
	if b.pawnProtect(pos, isWhite) {
		score += 2
	}
	// Score inverted for Black
	if isWhite {
		if pos == 46 || pos == 43 || pos == 22 || pos == 27 {
			score += 20
		} else if pos == 61 || pos == 68 {
			score -= 20
		}
		// check if checks ?
	} else {
		score = -score
		if pos == 56 || pos == 53 || pos == 72 || pos == 77 {
			score -= 20
		} else if pos == 31 || pos == 38 {
			score += 20
		}
	}
	return score
}

// evalRook evaluates the rook position.
func (b *Board) evalRook(pos int, isWhite bool) int {
	var score int
	score += 50
	// Invert for Black
	if b.pawnThreat(pos, isWhite) {
		score -= 50
	}
	// TODO:
	// Check for Open File
	// Check for Castle Possibility
	if isWhite {
		if pos < 80 && pos > 70 {
			score += 50
		}
	} else {
		score = -score
		if pos < 29 && pos > 20 {
			score -= 50
		}
	}

	return score
}

// evalQueen evaluates the queen position.
func (b *Board) evalQueen(pos int, isWhite bool) int {
	var score int
	score += 150 // I up the value because queens are worth taking..
	if b.pawnThreat(pos, isWhite) {
		score -= 200 // Because this is real dumb
	}
	if !isWhite {
		score = -score
	}
	return score
}

// evalKing evaluates king Position. Checks if in Check.
func (b *Board) evalKing(pos int, isWhite bool) int {
	var score int
	if isWhite {
		score += 100
	} else {
		score -= 100
	}
	// TODO check if castle
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

/*
DEPRICATED EVALUATION / ISN'T RELEVANT #####################
*/

// SumEval returns the sum of evaluations of one side.
// Depricated
func (b *Board) SumEval() int {
	origs, dests := b.SearchValid()
	evaluations := b.EvaluateMoves(origs, dests)
	var sum int
	for _, eval := range evaluations {
		sum += eval
	}
	return sum
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
		s := b.EvaluateMove(o, d)
		bests = append(bests, s)
	}
	return bests
}

// EvaluateMove scores a move based on the piece
// and its destination.
// TODO: Must I acknowledge castling?
func (b *Board) EvaluateMove(orig, dest int) int {
	var score int
	//var trade int
	isWhite := b.toMove == "w"
	piece := byte(unicode.ToUpper(rune(b.board[orig])))
	target := byte(unicode.ToUpper(rune(b.board[dest])))
	// If is Capture
	isCapture := target != '.'
	if b.pawnThreat(dest, isWhite) && piece != 'P' {
		score -= 70
	}

	// trade, Well, it's not necessarily a trade
	// This int tracks whether the attacker is valued
	// more or less than the target. It is a good thing
	// if it is valued less. Right?
	if isCapture {
		if target != 'P' {
			score += 40
		} else if target == 'P' && piece == 'P' && !b.pawnThreat(dest, isWhite) {
			score += 25
		} else if !b.pawnThreat(dest, isWhite) {
			score += 25
		}
	}
	if isCapture && target != 'P' {
		score += 20
	} else if isCapture && target == 'P' && piece == 'P' {
		score += 20
	}

	// isCenter and isBorder checks the destinations
	isCenter := dest == 55 || dest == 44 ||
		dest == 45 || dest == 54
	isBorder := (dest-1)%10 == 0 || (dest+2)%10 == 0
	if isBorder {
		switch {
		case piece == 'N':
			score -= 10
		case piece == 'P' && (dest > 40 && dest < 60):
			score -= 5
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
		var pot1, pot2 int
		if isWhite {
			pot1 = dest + 9
			pot2 = dest + 11
		} else {
			pot1 = dest - 9
			pot2 = dest - 11

		}
		if pot1 != '.' {
			if (b.isUpper(pot1) &&
				isWhite) || (!b.isUpper(pot1) && !isWhite) {
				// protecting
				score += 20
			} else if (b.isUpper(pot1) &&
				!isWhite) || (!b.isUpper(pot1) && isWhite) {
				// attacking
				score += 15
			}
		}
		if pot2 != '.' {
			if (b.isUpper(pot2) &&
				isWhite) || (!b.isUpper(pot2) && !isWhite) {
				// protecting
				score += 10
			} else if (b.isUpper(pot2) &&
				!isWhite) || (!b.isUpper(pot2) && isWhite) {
				// attacking
				score += 5
			}
		}
	case 'N':
		// Check for outpost Knight
		if isWhite {
			if dest > 50 {
				score += 30
			}
			if dest == 33 || dest == 36 {
				score += 30
			}
			if dest == 17 || dest == 12 {
				// Don't move back lol
				score -= 10
			}
		} else {
			if dest < 50 {
				score += 30
			}
			if dest == 66 || dest == 63 {
				score += 30
			}
			if dest == 87 || dest == 82 {
				// Don't move back lol
				score -= 10
			}
		}
	case 'B':
		// Check if on long access Bishop
		if dest%11 == 0 {
			score += 30
		}
		if dest == 43 || dest == 46 || dest == 53 || dest == 56 {
			score += 40
		}
	case 'R':
		// If seventh Rank
		if isWhite {
			if dest < 79 && dest > 70 {
				score += 50
			}
		} else {
			if dest < 29 && dest > 20 {
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

// FindBestMove returns the best move from evaluation.
func (b *Board) FindBestMove() (int, int) {
	origs, dests := b.SearchValid()
	bests := b.EvaluateMoves(origs, dests)
	// get best of bests and return the index
	var best int
	indexes := make([]int, 0)
	var i int
	for idx, val := range bests {
		if val > best {
			best = val
			indexes = bests[:0]
			indexes = append(indexes, idx)
		} else if val == best {
			indexes = append(indexes, idx)
		} else {
			continue
		}

	}
	//fmt.Println(best) // Lol this shit is crazy
	if len(indexes) > 1 {
		i = rand.Intn(len(indexes))
		//for _, x := range indexes {
		//	fmt.Println(origs[x])
		//}
	} else {
		i = indexes[0]
	}
	return origs[i], dests[i]
}

// MoveBest finds the best move of all valid moves.
// This method is not operational.
func (b *Board) MoveBest() {
	origs, dests := b.SearchValid()
	bests := b.EvaluateMoves(origs, dests)
	// get best of bests and return the index
	var best int
	indexes := make([]int, 0)
	var i int
	for idx, val := range bests {
		if val > best {
			best = val
			indexes = bests[:0]
			indexes = append(indexes, idx)
		} else if val == best {
			indexes = append(indexes, idx)
		} else {
			continue
		}

	}
	//fmt.Println(best) // Lol this shit is crazy
	if len(indexes) > 1 {
		i = rand.Intn(len(indexes))
		//for _, x := range indexes {
		//	fmt.Println(origs[x])
		//}
	} else {
		i = indexes[0]
	}
	b.Move(origs[i], dests[i])
}
