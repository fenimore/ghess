package ghess

import (
	"errors"
	"unicode"
)

// Problem: 1nbq1knr/1rNpppb1/pp4pp/4N3/3P4/P6P/1PP1PPP1/R1BQKB1R w KQ-- - 0 9

// Move is the basic validation.
// The origin and destination square are tested
// in a dereferenced b.board to keep from moving
// into check. The individual pieces are validated
// in separate methods. Finally this method updates
// the board, updateBoard().
func (b *Board) Move(orig, dest int) error {
	if b.Checkmate {
		return errors.New("Cannot Move in Checkmate")
	}
	val := b.board[orig]
	var o byte           // supposed starting square
	var d byte           // supposed destination
	var isEmpassant bool // refactor?
	var isCastle bool
	if b.toMove == "w" {
		//b[pos] = byte(unicode.ToUpper(rune(b[pos])))
		// check that orig is Upper
		// and dest is Enemy or Empty
		// Use hash map TODO
		//		o = ByteToUpper[b.board[orig]]
		//		d = ByteToLower[b.board[dest]]
		o = byte(unicode.ToUpper(rune(b.board[orig])))
		//[]byte(bytes.ToUpper(b.board[orig : orig+1]))[0]
		d = byte(unicode.ToLower(rune(b.board[dest])))
		//[]byte(bytes.ToLower(b.board[dest : dest+1]))[0]
	} else if b.toMove == "b" {
		// check if orig is Lower
		// and dest is Enemy or Empty
		//		o = ByteToLower[b.board[orig]]
		//		d = ByteToUpper[b.board[dest]]
		o = byte(unicode.ToLower(rune(b.board[orig])))
		//[]byte(bytes.ToLower(b.board[orig : orig+1]))[0]
		d = byte(unicode.ToUpper(rune(b.board[dest])))

		//[]byte(bytes.ToUpper(b.board[dest : dest+1]))[0]
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

	p := b.board[orig]
	switch p {
	case 'p', 'P':
		e := b.validPawn(orig, dest)
		if e != nil {
			return e
		}
		emp := dest - orig
		if emp > 11 || emp < -11 {
			isEmpassant = true
		}

	case 'n', 'N':
		e := b.validKnight(orig, dest)
		if e != nil {
			return e
		}
	case 'b', 'B':
		e := b.validBishop(orig, dest)
		if e != nil {
			return e
		}
	case 'R', 'r':
		e := b.validRook(orig, dest)
		if e != nil {
			return e
		}
	case 'Q', 'q':
		e := b.validQueen(orig, dest)
		if e != nil {
			return e
		}
	case 'k', 'K': // is castle?
		if !isCastle {
			e := b.validKing(orig, dest, false)
			if e != nil {
				return e
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
	possible := CopyBoard(b)
	// Check possibilities
	possible.updateBoard(orig, dest, val, isEmpassant, isCastle)
	putIntoCheck, isCheck := possible.checkCheck(isWhite) // because turn updateso
	//isCheck := possible.isOpponentInCheck()
	if isCheck {
		return errors.New("Cannot move into Check")
	}
	if isCastle {
		// FIXME, reuse this value?
		isCheck = b.isPlayerInCheck()
		if isCheck {
			return errors.New("Cannot Castle in Check")
		}
		possible = CopyBoard(b)
		switch {
		case isWhite && dest < orig:
			possible.updateBoard(orig, 13, 'K',
				false, false) //King side, 13
		case isWhite && dest > orig:
			possible.updateBoard(orig, 15, 'K',
				false, false) // Queen side, 15
		case !isWhite && dest < orig:
			possible.updateBoard(orig, 83, 'k',
				false, false) // King 83
		case !isWhite && dest > orig:
			possible.updateBoard(orig, 85, 'k',
				false, false) // Queen 85
		}

		isCheck = possible.isOpponentInCheck()
		if isCheck {
			return errors.New("Cannot Castle through check")
		}
	}
	// If all goes well:
	// update real board
	b.updateBoard(orig, dest, val, isEmpassant, isCastle)

	if putIntoCheck {
		b.Check = true
		_ = b.PlayerCheckMate()
	} else {
		b.Check = false
	}

	// Check if it is draw
	if orig == b.history[6] && orig == b.history[3] && b.history[0] == b.history[5] {
		// origins all match upppp... suspicious
		if dest == b.history[7] && dest == b.history[2] && b.history[1] == b.history[4] {
			b.Score = "1/2 - 1/2"
			b.Draw = true
		}
	}
	// For draw?
	b.cycleHistory(orig, dest)
	return nil
}

// updateBoard changes the byte values of board.
// It is useless without validation from Move().
// This method checks, and sets, Check for Board.board.
func (b *Board) updateBoard(orig, dest int,
	val byte, isEmpassant, isCastle bool) {
	isWhite := b.toMove == "w"
	var isPromotion bool
	// Check for Promotion
	switch {
	case val == 'p' && dest < 20:
		isPromotion = true
	case val == 'P' && dest > 80:
		isPromotion = true
	}
	// Check for castle deactivation
	switch {
	case val == 'r' || val == 'R':
		switch { // Castle
		case orig == 18:
			b.castle[1] = '-'
		case orig == 11:
			b.castle[3] = '-'
		case orig == 88:
			b.castle[0] = '-'
		case orig == 81:
			b.castle[2] = '-'
		}
	case orig == 14 || orig == 84:
		switch {
		case val == 'K':
			b.castle[0], b.castle[1] = '-', '-'
		case val == 'k':
			b.castle[2], b.castle[3] = '-', '-'
		default:
			break
		}
	case isCastle:
		switch {
		case isWhite:
			b.castle[0], b.castle[1] = '-', '-'
		case !isWhite:
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

}

// PlayerCheck Method checks if Current player is in Check
// and updates score accordingly
func (b *Board) PlayerCheck() bool {
	isCheck := b.isPlayerInCheck()
	if isCheck {
		b.Check = true
	} else {
		b.Check = false
	}

	return isCheck
}

func (b *Board) GameOver() bool {
	origs, _ := b.SearchValid()
	return len(origs) < 1
}

// PlayerCheckMate checks if current player is checkmated
// and updates score accordingly. Only call after setting
// the Check field in Board.
func (b *Board) PlayerCheckMate() bool {
	origs, _ := b.SearchValid()
	if len(origs) < 1 {
		b.Checkmate = true
		if b.toMove == "w" {
			b.Score = "0-1"
			return true
		} else {
			b.Score = "1-0"
			return true
		}
	}
	return false
}

// isPlayerInCheck, current player is in Check.
// TODO: Change to upper case?
func (b *Board) isPlayerInCheck() bool {
	isWhite := b.toMove == "w"
	for idx, val := range b.board {
		if val == 'K' && isWhite {
			//fmt.Println(b.board[idx])
			return b.isInCheck(idx)
		}
		if val == 'k' && !isWhite {
			return b.isInCheck(idx)
		}
	}
	return false
}

// isOpponentInCheck checks if move put the other in Check
func (b *Board) isOpponentInCheck() bool {
	isWhite := b.toMove == "w"
	for idx, val := range b.board {
		if val == 'K' && !isWhite {
			//fmt.Println(b.board[idx])
			return b.isInCheck(idx)
		}
		if val == 'k' && isWhite {
			return b.isInCheck(idx)
		}
	}
	return false
}

func (b *Board) checkCheck(isBlack bool) (bool, bool) {
	// This checks who is in check. IsWhite is the inverse of
	// the validation check. So that means that if isWhite, these
	// black player CANNOT be in check, and if isBlack,
	// the white player cannot be in check
	//var moveIntoCheck bool
	var putIntoCheck bool
CheckLoop:
	for idx, val := range b.board {
		switch val {
		case 'K': // white player
			check := b.isInCheck(idx)
			if isBlack && check {
				//moveIntoCheck = true
				return false, true
			} else if !isBlack && check {
				putIntoCheck = true
			} else {
				continue CheckLoop
			}
		case 'k':
			check := b.isInCheck(idx)
			if !isBlack && check {
				//moveIntoCheck = true
				return false, true
			} else if isBlack && check {
				putIntoCheck = true
			} else {
				continue CheckLoop
			}
		default:
			continue CheckLoop
		}
	}
	return putIntoCheck, false //moveIntoCheck
}

// TODO: Break these into functions
func (b *Board) isInCheck(target int) bool {

	isWhite := b.isUpper(target)

	switch {
	case b.checkVerticalAxis(target, isWhite):
		return true
	case b.checkHorizontalAsix(target, isWhite):
		return true
	case b.checkA1Diagonal(target, isWhite):
		return true
	case b.checkH1Diagonal(target, isWhite):
		return true
	case b.checkProximity(target, isWhite):
		return true
	default:
		return false
	}
}

// isInCheck checks if target King is in Check.
// Automaticaly checks for turn by the target King.
func (b *Board) isInCheckDep(target int) bool {
	isWhite := b.isUpper(target)
	//k := b.board[target]

	// store all the orig of the opponents pieces
	attackers := make([]int, 0, 16)

	for idx := range b.board {
		whitePiece := b.isUpper(idx)
		if isWhite && !whitePiece {
			if b.board[idx] == 'p' {
				if (idx - target) > 11 {
					continue
				}
			}
			attackers = append(attackers, idx)
		} else if !isWhite && whitePiece { // black
			if b.board[idx] == 'P' {
				if (target - idx) > 11 {
					continue
				}
			}
			attackers = append(attackers, idx)
		}
	}
	// Check for Attackers
	// and only add the oproprate ones.
	for _, val := range attackers {
		p := b.board[val]
		switch p {
		case 'P', 'p':
			e := b.validPawn(val, target)
			if e == nil {
				return true
			}
		case 'N', 'n':
			e := b.validKnight(val, target)
			if e == nil {
				//fmt.Println("Knight check")
				return true
			}
		case 'B', 'b':
			e := b.validBishop(val, target)
			if e == nil {
				return true
			}
		case 'R', 'r':
			e := b.validRook(val, target)
			if e == nil {
				//fmt.Println("Rook check")
				return true
			}
		case 'Q', 'q':
			e := b.validQueen(val, target)
			if e == nil {
				return true
			}
		case 'K', 'k':
			e := b.validKing(val, target, false)
			if e == nil {
				return true
			}
		}
	}
	return false // Not in Check
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
	switch remainder {
	case 10:
		if b.board[dest] != '.' {
			return err
		}
	case 20:
		if b.board[dest] != '.' {
			return err
		}
		if isWhite && b.board[dest-10] != '.' {
			return err
		} else if !isWhite && b.board[dest+10] != '.' {
			return err
		}
		if orig > 28 && b.toMove == "w" { // Only from 2nd rank
			return err
		} else if orig < 70 && b.toMove == "b" {
			return err
		}
	case 9, 11:
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
	if a1h8 == 0 {
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
	} else if a8h1 == 0 {
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
	a8h1 := remainder%9 == 0                      // Diag a8h1 A8
	a1h8 := remainder%11 == 0                     // Diag a1h8 A1
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
	} else if a8h1 {
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
	} else if a1h8 {
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
	inValidCastle := dest != 88 && dest != 81 && dest != 11 && dest != 18
	validCastle := !inValidCastle && (orig == 14 || orig == 84)
	// validCastle is a not so valid castle position
	if !validCastle && castle {
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
		isQueenSide := dest > orig
		isKingSide := dest < orig
		validQueenSide := !(g[orig+1] != '.' ||
			g[orig+2] != '.' || g[orig+3] != '.')
		validKingSide := !(g[orig-1] != '.' ||
			g[orig-2] != '.')
		isWhite := b.toMove == "w"
		// Is castle Left?
		whiteQueenCastle := b.castle[1] == 'Q'
		whiteKingCastle := b.castle[0] == 'K'
		blackQueenCastle := b.castle[3] == 'q'
		blackKingCastle := b.castle[2] == 'k'
		if isQueenSide && validQueenSide {
			if isWhite && whiteQueenCastle {
				return nil
			} else if !isWhite && blackQueenCastle {
				return nil
			} else {
				return noCastle
			}
		} else if isKingSide && validKingSide {
			if isWhite && whiteKingCastle {
				return nil
			} else if !isWhite && blackKingCastle {
				return nil
			} else {
				return noCastle
			}
		} else {
			return castlerr
		}

	} else {
		return errors.New("Illegal King Move")
	}
	//return nil
}

// Check if King is in Check on Horizontal axis
func (b *Board) checkHorizontalAsix(target int, isWhite bool) bool {
RightLoop:
	for i := target + 1; !((i+1)%10 == 0); i = i + 1 {
		switch b.board[i] {
		case 'r', 'q':
			if isWhite {
				return true
			}
			break RightLoop
		case 'R', 'Q':
			if !isWhite {
				return true
			}
			break RightLoop
		case '.':
			continue RightLoop
		default:
			break RightLoop
		}
	}
LeftLoop:
	for i := target - 1; !(i%10 == 0); i = i - 1 {
		switch b.board[i] {
		case 'r', 'q':
			if isWhite {
				return true
			}
			break LeftLoop
		case 'R', 'Q':
			if !isWhite {
				return true
			}
			break LeftLoop
		case '.':
			continue LeftLoop
		default:
			break LeftLoop
		}
	}
	return false
}

// Check if King (piece) is in check on Vertical Axis
func (b *Board) checkVerticalAxis(target int, isWhite bool) bool {
UpVerLoop:
	for i := target + 10; i < 88; i = i + 10 {

		switch b.board[i] {
		case 'r', 'q':
			if isWhite {
				return true
			}
			break UpVerLoop
		case 'R', 'Q':
			if !isWhite {
				return true
			}
			break UpVerLoop
		case '.':
			continue UpVerLoop
		default:
			break UpVerLoop
		}
	}
DownVerLoop:
	for i := target - 10; i > 10; i = i - 10 {
		// Should stop when off the board
		switch b.board[i] {
		case 'r', 'q':
			if isWhite {
				return true
			}
			break DownVerLoop
		case 'R', 'Q':
			if !isWhite {
				return true
			}
			break DownVerLoop
		case '.':
			continue DownVerLoop
		default:
			break DownVerLoop
		}
	}
	return false
}

func (b *Board) checkA1Diagonal(target int, isWhite bool) bool {
a1h8Loop:
	for i := target + 9; i < 89; i = i + 9 {
		// Should stop when off the board
		switch b.board[i] {
		case 'b', 'q':
			if isWhite {
				return true
			}
			break a1h8Loop
		case 'B', 'Q':
			if !isWhite {
				return true
			}
			break a1h8Loop
		case '.':
			continue a1h8Loop
		default:
			break a1h8Loop
		}
	}
h8a1Loop:
	for i := target - 9; i > 10; i = i - 9 {
		// Should stop when off the board
		switch b.board[i] {
		case 'b', 'q':
			if isWhite {
				return true
			}
			break h8a1Loop
		case 'B', 'Q':
			if !isWhite {
				return true
			}
			break h8a1Loop
		case '.':
			continue h8a1Loop
		default:
			break h8a1Loop
		}
	}
	return false
}

func (b *Board) checkH1Diagonal(target int, isWhite bool) bool {
h1a8Loop:
	for i := target + 11; i < 89; i = i + 11 {
		// Should stop when off the board
		switch b.board[i] {
		case 'b', 'q':
			if isWhite {
				return true
			}
			break h1a8Loop
		case 'B', 'Q':
			if !isWhite {
				return true
			}
			break h1a8Loop
		case '.':
			continue h1a8Loop
		default:
			break h1a8Loop
		}
	}
a8h1Loop:
	for i := target - 11; i > 10; i = i - 11 {
		// Should stop when off the board
		switch b.board[i] {
		case 'b', 'q':
			if isWhite {
				return true
			}
			break a8h1Loop
		case 'B', 'Q':
			if !isWhite {
				return true
			}
			break a8h1Loop
		case '.':
			continue a8h1Loop
		default:
			break a8h1Loop
		}
	}
	return false
}

// Check for Knight, Kings or Pawns in Proximity of Piece
func (b *Board) checkProximity(target int, isWhite bool) bool {

	// The immediate Proximity
	A := b.board[target-1]
	B := b.board[target+1]
	C := b.board[target-11] // Pawn
	D := b.board[target+11] // Pawn
	E := b.board[target-9]  // Pawn
	F := b.board[target+9]  // Pawn
	G := b.board[target-10]
	H := b.board[target+10]

	var N, M, O, P, S, T byte
	// Possible Knight Positions
	if target > 18 {
		O = b.board[target-21]
		P = b.board[target-19]
		T = b.board[target-12]
	} else {
		O = '.'
		P = '.'
		T = '.'
	}

	if target < 80 {
		N = b.board[target+21]
		M = b.board[target+19]
		S = b.board[target+12]
	} else {
		N = '.'
		M = '.'
		S = '.'
	}
	Q := b.board[target+8]
	R := b.board[target-8]

	if isWhite {
		// for white player
		switch byte('k') {
		case A, B, C, D, E, F, G, H:
			return true
		}

		switch byte('p') {
		case D, F:
			return true
		}
		switch byte('n') {
		case N, M, O, P, Q, R, S, T:
			return true
		}
	} else {
		// for white player
		switch byte('K') {
		case A, B, C, D, E, F, G, H:
			return true
		}

		switch byte('P') {
		case E, C:
			return true
		}
		switch byte('N') {
		case N, M, O, P, Q, R, S, T:
			return true
		}
	}
	return false
}
