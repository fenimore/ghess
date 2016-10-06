package ghess

import (
	"bytes"
	"errors"
)

// Move is the basic validation.
// The origin and destination square are tested
// in a dereferenced b.board to keep from moving
// into check. The individual pieces are validated
// in separate methods. Finally this method updates
// the board, updateBoard().
func (b *Board) Move(orig, dest int) error {
	if b.checkmate {
		return errors.New("Cannot Move in Checkmate")
	}
	val := b.board[orig]
	var o byte           // supposed starting square
	var d byte           // supposed destination
	var isEmpassant bool // refactor?
	var isCastle bool
	if b.toMove == "w" {
		// check that orig is Upper
		// and dest is Enemy or Empty
		o = []byte(bytes.ToUpper(b.board[orig : orig+1]))[0]
		d = []byte(bytes.ToLower(b.board[dest : dest+1]))[0]
	} else if b.toMove == "b" {
		// check if orig is Lower
		// and dest is Enemy or Empty
		o = []byte(bytes.ToLower(b.board[orig : orig+1]))[0]
		d = []byte(bytes.ToUpper(b.board[dest : dest+1]))[0]
	}
	// Check for Castle
	if orig == 14 {
		isCastle = b.board[dest] == 'R'
	} else if orig == 84 {
		isCastle = b.board[dest] == 'r'
	}

	err := b.basicValidation(orig, dest, o, d, isCastle)
	if err != nil {
		return err
	}

	p := string(bytes.ToUpper(b.board[orig : orig+1]))
	switch {
	case p == "P":
		e := b.validPawn(orig, dest)
		if e != nil {
			return e
		}
		emp := dest - orig
		if emp > 11 || emp < -11 {
			isEmpassant = true
		}

	case p == "N":
		e := b.validKnight(orig, dest)
		if e != nil {
			return e
		}
	case p == "B":
		e := b.validBishop(orig, dest)
		if e != nil {
			return e
		}
	case p == "R":
		e := b.validRook(orig, dest)
		if e != nil {
			return e
		}
	case p == "Q":
		e := b.validQueen(orig, dest)
		if e != nil {
			return e
		}
	case p == "K": // is castle?
		if !isCastle {
			e := b.validKing(orig, dest, false)
			if e != nil {
				return e
			}
			if orig == 14 || orig == 84 { // starting pos
				switch {
				case o == 'K':
					b.castle[0], b.castle[1] = '-', '-'
				case o == 'k':
					b.castle[2], b.castle[3] = '-', '-'
				}
			}
		} else {
			e := b.validKing(orig, dest, true)
			if e != nil {
				return e
			}
		}
	}
	// Make sure new position doesn't put in check
	isWhite := b.toMove == "w"
	possible := *b                 // slices are  still pointing...
	boardCopy := make([]byte, 120) // b.board is Pointer
	castleCopy := make([]byte, 4)
	copy(boardCopy, b.board)
	copy(castleCopy, b.castle)
	possible.board = boardCopy
	possible.castle = castleCopy
	// Check possibilities
	possible.updateBoard(orig, dest, val, isEmpassant, isCastle)
	// find mover's king
	var king int
	for idx, val := range possible.board {
		if isWhite && val == 'K' {
			king = idx
			break
		} else if !isWhite && val == 'k' {
			king = idx
			break
		}
	}
	isCheck := possible.isInCheck(king)
	if isCheck {
		return errors.New("Cannot move into Check")
	}
	if isCastle {
		copy2 := make([]byte, 120)
		copy(copy2, b.board)
		possible.board = copy2
		switch {
		case isWhite && dest < orig:
			possible.updateBoard(orig, 13, 'K',
				false, false) //King side, 13
			king = 13
		case isWhite && dest > orig:
			possible.updateBoard(orig, 15, 'K',
				false, false) // Queen side, 15
			king = 15
		case !isWhite && dest < orig:
			possible.updateBoard(orig, 83, 'k',
				false, false) // King 83
			king = 83
		case !isWhite && dest > orig:
			possible.updateBoard(orig, 85, 'k',
				false, false) // Queen 85
			king = 85
		}

		isCheck = possible.isInCheck(king)
		if isCheck {
			return errors.New("Cannot Castle through check")
		}
	}
	// update real board
	b.updateBoard(orig, dest, val, isEmpassant, isCastle)
	// Check if it is draw// If not TODO
	if orig == b.history[6] && orig == b.history[3] && b.history[0] == b.history[5] {
		// origins all match upppp... suspicious
		if dest == b.history[7] && dest == b.history[2] && b.history[1] == b.history[4] {
			b.score = "1/2 - 1/2"
		}
	}
	b.cycleHistory(orig, dest)
	// Look for Checkmate
	// Check all possibl moves after a check?
	isCheck = b.isPlayerInCheck()
	//isCheck = b.isInCheck(king)
	if isCheck {
		isCheckMate := false
		origs, _ := b.SearchValid()
		if len(origs) < 1 {
			isCheckMate = true
		}
		if isCheckMate {
			b.checkmate = true
			if b.toMove == "w" {
				b.score = "0-1"
			} else {
				b.score = "1-0"
			}
		}
	}

	return nil
}

// updateBaord changes the byte values of board.
// It is useless without validation from Move().
// This method checks, and sets, Check for Board.board.
func (b *Board) updateBoard(orig, dest int,
	val byte, isEmpassant, isCastle bool) {
	isWhite := b.toMove == "w"
	var isPromotion bool
	// Check for Promotion
	switch {
	case b.board[orig] == 'p' && dest < 20:
		isPromotion = true
	case b.board[orig] == 'P' && dest > 80:
		isPromotion = true
	}
	// Check for castle deactivation
	switch {
	case b.board[orig] == 'r' || b.board[orig] == 'R':
		switch { // Castle
		case orig == b.pgnMap["a1"]:
			b.castle[1] = '-'
		case orig == b.pgnMap["a8"]:
			b.castle[3] = '-'
		case orig == b.pgnMap["h1"]:
			b.castle[0] = '-'
		case orig == b.pgnMap["h8"]:
			b.castle[2] = '-'
		}
	case isCastle:
		kingSide := orig > dest
		queenSide := orig < dest
		switch {
		case isWhite && kingSide:
			b.castle[0], b.castle[1] = '-', '-'
		case isWhite && queenSide:
			b.castle[0], b.castle[1] = '-', '-'
		case !isWhite && kingSide:
			b.castle[2], b.castle[3] = '-', '-'
		case !isWhite && queenSide:
			b.castle[2], b.castle[3] = '-', '-'
		}
	}

	// Check for Attack on Empassant
	if val == 'p' || val == 'P' {
		switch {
		case dest-orig == 9 || dest-orig == 11:
			if b.board[dest] == '.' {
				// White offset
				b.board[dest-10] = '.'

			}
		case orig-dest == 9 || orig-dest == 11:
			if b.board[dest] == '.' {
				// Black Offset
				b.board[dest+10] = '.'
			}
		}
	}

	// Set origin
	b.board[orig] = '.'

	// Set destination
	if isCastle {
		if dest > orig { // queen side
			b.board[dest-2],
				b.board[dest-3] = val, b.board[dest]
		} else { // king side
			b.board[dest+1],
				b.board[dest+2] = val, b.board[dest]
		}
		b.board[dest] = '.'
	} else if isPromotion {
		switch {
		case dest < 20:
			b.board[dest] = 'q'
		case dest > 80:
			b.board[dest] = 'Q'
		}
	} else { // Normal Move/Capture
		b.board[dest] = val
	}

	// TODO check for Check
	// Update Game variables
	if b.toMove == "w" {
		b.toMove = "b"
	} else {
		b.moves++ // add one to move count
		b.toMove = "w"
	}
	if isEmpassant {
		b.empassant = dest
	} else {
		b.empassant = 0
	}

	// Check if move put other player in Check
	isCheck := b.isPlayerInCheck()
	if isCheck {
		b.check = true
	} else {
		b.check = false
	}
}

// isPlayerInCheck, current player is in Check.
// TODO: Change to upper case?
func (b *Board) isPlayerInCheck() bool {
	isWhite := b.toMove == "w"
	for idx, val := range b.board {
		if val == 'K' && b.isUpper(idx) && isWhite {
			//fmt.Println(b.board[idx])
			return b.isInCheck(idx)
		}
		if val == 'k' && !b.isUpper(idx) && !isWhite {
			return b.isInCheck(idx)
		}
	}
	return false
}

// isInCheck checks if target King is in Check.
// Automaticaly checks for turn by the target King.
func (b *Board) isInCheck(target int) bool {
	isWhite := b.isUpper(target)
	//k := b.board[target]

	// store all the orig of the opponents pieces
	attackers := make([]int, 0, 16)

	for idx, _ := range b.board {
		whitePiece := b.isUpper(idx)
		if isWhite && !whitePiece {
			attackers = append(attackers, idx)
		} else if !isWhite && whitePiece { // black
			attackers = append(attackers, idx)
		}
	}
	//fmt.Println("white ", isWhite, "attackers ", attackers, "king", k)
	// check for valid attacks
	for _, val := range attackers {
		p := string(bytes.ToUpper(b.board[val : val+1]))
		switch {
		case p == "P":
			e := b.validPawn(val, target)
			if e == nil {
				return true
			}
		case p == "N":
			e := b.validKnight(val, target)
			if e == nil {
				//fmt.Println("Knight check")
				return true
			}
		case p == "B":
			e := b.validBishop(val, target)
			if e == nil {
				return true
			}
		case p == "R":
			e := b.validRook(val, target)
			if e == nil {
				//fmt.Println("Rook check")
				return true
			}
		case p == "Q":
			e := b.validQueen(val, target)
			if e == nil {
				return true
			}
		case p == "K":
			e := b.validKing(val, target, false)
			if e == nil {
				return true
			}
		}
	}
	// if nothing was valid, return false
	return false
}

// basicValidation assures basic chess rules:
// correct-color to move, origin is empty, and
// only attack an enemy.
func (b *Board) basicValidation(orig, dest int, o, d byte, isCastle bool) error {
	// Check if it is the right turn
	if b.board[orig] != o {
		return errors.New("Not your turn")
	}
	// Check if Origin is Empty
	if o == '.' {
		return errors.New("Empty square")
	}
	// Check if destination is Enemy
	if b.board[dest] != d && !isCastle { //
		return errors.New("Can't attack your own piece")
	}
	return nil
}

// validatePawn Move
// d param is destination byte, used mainly fro empassant
func (b *Board) validPawn(orig int, dest int) error {
	err := errors.New("Illegal Pawn Move")
	var remainder int
	var empOffset int
	var empTarget byte
	isWhite := b.isUpper(orig)
	// Whose turn
	if isWhite {
		remainder = dest - orig
		empOffset = -10 // where the empassant piece should be
		empTarget = 'p'
	} else {
		remainder = orig - dest
		empOffset = 10
		empTarget = 'P'
	}
	// What sort of move
	switch {
	case remainder == 10:
		if b.board[dest] != '.' {
			return err
		}
	case remainder == 20:
		if b.board[dest] != '.' {
			return err
		}
		if orig > 28 && b.toMove == "w" { // Only from 2nd rank
			return err
		} else if orig < 70 && b.toMove == "b" {
			return err
		}
	case remainder == 9 || remainder == 11:
		if b.board[dest] == '.' && dest+empOffset == b.empassant {
			// Empassant attack
			if b.board[dest+empOffset] != empTarget {
				return err
			}
		} else if b.board[dest] == '.' {
			return err
		}
	default:
		return errors.New("Not valid Pawn move.")
	}
	return nil
}

// validateKnight move.
func (b *Board) validKnight(orig int, dest int) error {
	var possibilities [8]int
	possibilities[0], possibilities[1],
		possibilities[2], possibilities[3],
		possibilities[4], possibilities[5],
		possibilities[6], possibilities[7] = orig+21,
		orig+19, orig+12, orig+8, orig-8,
		orig-12, orig-19, orig-21
	for _, possibility := range possibilities {
		if possibility == dest {
			return nil
		}
	}
	return errors.New("Illegal Knight Move")
}

// validateBishop move.
func (b *Board) validBishop(orig int, dest int) error {
	// Check if other pieces are in the way
	err := errors.New("Illegal Bishop Move")
	trajectory := orig - dest
	a1h8 := trajectory % 11 // if 0 remainder...
	a8h1 := trajectory % 9
	// Check which slope
	if a8h1 == 0 {
		if dest > orig { // go to bottom right
			for i := orig + 11; i <= dest-11; i += 11 {
				if b.board[i] != '.' {
					return err
				}
			}
		} else if dest < orig { // go to top left
			for i := orig - 11; i >= dest+11; i -= 11 {
				if b.board[i] != '.' {
					return err
				}
			}
		}
	} else if a1h8 == 0 {
		if dest > orig { // go to bottem left
			for i := orig + 9; i <= dest-9; i += 9 {
				if b.board[i] != '.' {
					return err
				}
			}
		} else if orig > dest { // go to top right
			for i := orig - 9; i >= dest+9; i -= 9 {
				if b.board[i] != '.' {
					return err
				}
			}
		}
	} else {
		return errors.New("Bishop move not valid")
	}
	return nil
}

// validRook validate rook move.
func (b *Board) validRook(orig int, dest int) error {
	// Check if pieces are in the way
	err := errors.New("Illegal Rook Move")
	remainder := dest - orig
	if remainder < 10 && remainder > -10 {
		// Horizontal
		if remainder < 0 {
			for i := orig - 1; i > dest; i-- {
				if b.board[i] != '.' {
					return err
				}
			}
		} else {
			for i := orig + 1; i < dest; i++ {
				if b.board[i] != '.' {
					return err
				}
			}
		}
	} else {
		if remainder%10 != 0 {
			return err
		}
		// Vertical
		if remainder < 0 { // descends
			for i := orig - 10; i > dest; i -= 10 {
				if b.board[i] != '.' {
					return err
				}
			}
		} else {
			for i := orig + 10; i < dest; i += 10 {
				if b.board[i] != '.' {
					return err
				}
			}
		}
	}
	return nil
}

// validQueen validates queen move.
func (b *Board) validQueen(orig int, dest int) error {
	err := errors.New("Illegal Queen Move")
	remainder := dest - orig
	vertical := remainder%10 == 0
	horizontal := remainder < 9 && remainder > -9 // Horizontal
	diagA8 := remainder%9 == 0                    // Diag a8h1
	diagA1 := remainder%11 == 0                   // Diag a1h8
	// Check if moves through not-empty squares
	if horizontal { // 1st
		if remainder < 0 {
			for i := orig - 1; i > dest; i-- {
				if b.board[i] != '.' {
					return err
				}
			}
		} else { // go right
			for i := orig + 1; i < dest; i++ {
				if b.board[i] != '.' {
					return err
				}
			}
		}
	} else if vertical {
		if remainder < 0 {
			for i := orig - 10; i > dest; i -= 10 {
				if b.board[i] != '.' {
					return err
				}
			}
		} else {
			for i := orig + 10; i < dest; i += 10 {
				if b.board[i] != '.' {
					return err
				}
			}
		}
	} else if diagA8 {
		if dest > orig { // go to bottem left
			for i := orig + 9; i <= dest-9; i += 9 {
				if b.board[i] != '.' {
					return err
				}
			}
		} else if orig > dest { // go to top right
			for i := orig - 9; i >= dest+9; i -= 9 {
				if b.board[i] != '.' {
					return err
				}
			}
		}
	} else if diagA1 {
		if dest > orig { // go to bottom right
			for i := orig + 11; i <= dest-11; i += 11 {
				if b.board[i] != '.' {
					return err
				}
			}
		} else if dest < orig { // go to top left
			for i := orig - 11; i >= dest+11; i -= 11 {
				if b.board[i] != '.' {
					return err
				}
			}
		}
	} else {
		return errors.New("Illegal Queen Move")
	}
	// check if anything is in between

	return nil
}

// validKing alidates king move.
// Check for castle.
func (b *Board) validKing(orig int, dest int, castle bool) error {
	validCastle := dest != 88 && dest != 81 && dest != 11 && dest != 18
	// validCastle is a not so valid castle position
	if validCastle && castle {
		return errors.New("Castle by moving K to R position")
	}
	castlerr := errors.New("Something is in your way")
	noCastle := errors.New("Castle on this side is foutu")
	var possibilities [8]int
	g := b.board // g for gameboard
	possibilities[0], possibilities[1],
		possibilities[2], possibilities[3],
		possibilities[4], possibilities[5],
		possibilities[6], possibilities[7] = orig+10,
		orig+11, orig+1, orig+9, orig-10,
		orig-11, orig-1, orig-9
	for _, possibility := range possibilities {
		if possibility == dest {
			return nil
		}
	}
	if castle {
		queenSideCastle := !(g[orig+1] != '.' || g[orig+2] != '.' || g[orig+3] != '.')
		kingSideCastle := !(g[orig-1] != '.' || g[orig-2] != '.')
		if dest > orig { // Queen side
			if !queenSideCastle {
				return castlerr
			}
			if b.toMove == "w" {
				if b.castle[1] != 'Q' {
					return noCastle
				}

			} else { // b
				if b.castle[3] != 'q' {
					return noCastle
				}
			}
		} else if orig > dest {
			if !kingSideCastle {
				return castlerr
			}
			if b.toMove == "w" {
				if b.castle[0] != 'K' {
					return noCastle
				}
			} else {
				if b.castle[2] != 'k' {
					return noCastle
				}
			}
		}

	} else {
		return errors.New("Illegal King Move")
	}
	return nil
}
