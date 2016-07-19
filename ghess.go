/*
Go Chess Engine - Ghess
Fenimore Love 2016
GPLv3

TODO: Evaluation

Exported Methods and Functions:

Position()
ParseMove()
Move()
LoadPgn()
LoadFen()
NewBoard()
Stats()
String()
PgnString()

*/
package ghess

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// Board is a chessboard type
// TODO: Make Upper Case? M-c for upper case
type Board struct {
	board []byte // piece position
	// Game Variables
	castle    []byte // castle possibility KQkq or ----
	empassant int    // square vulnerable to empassant
	score     string
	toMove    string // Next move is w or b
	moves     int    // the count of moves
	check     bool
	checkmate bool // start Capitalizing
	// Map for display grid
	pgnMap   map[string]int    // the pgn format
	pieceMap map[int]string    // coord to standard notation
	pieces   map[string]string // the unicode fonts
	rows     map[int][8]int       // rows for white/black squaring
	// Game Positions
	fen        string         // Game position
	pgn        string         // Game history
	headers    string         // Pgn format
	pgnPattern *regexp.Regexp // For parsing PGN
	fenPattern *regexp.Regexp
}

// Create a new Board in the starting position
func NewBoard() Board {
	b := make([]byte, 120)
	// starting position
	b = []byte(`           RNBKQBNR  PPPPPPPP  ........  ........  ........  ........  pppppppp  rnbkqbnr                                `)

	// Printed Board Notations
	b[91], b[92], b[93], b[94], b[95], b[96], b[97], b[98] = 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'
	b[19], b[29], b[39], b[49], b[59], b[69], b[79], b[89] = '1', '2', '3', '4', '5', '6', '7', '8'

	// Map of PGN notation
	m := make(map[string]int)
	m["a1"], m["b1"], m["c1"], m["d1"], m["e1"], m["f1"], m["g1"], m["h1"] = 18, 17, 16, 15, 14, 13, 12, 11
	m["a2"], m["b2"], m["c2"], m["d2"], m["e2"], m["f2"], m["g2"], m["h2"] = 28, 27, 26, 25, 24, 23, 22, 21
	m["a3"], m["b3"], m["c3"], m["d3"], m["e3"], m["f3"], m["g3"], m["h3"] = 38, 37, 36, 35, 34, 33, 32, 31
	m["a4"], m["b4"], m["c4"], m["d4"], m["e4"], m["f4"], m["g4"], m["h4"] = 48, 47, 46, 45, 44, 43, 42, 41
	m["a5"], m["b5"], m["c5"], m["d5"], m["e5"], m["f5"], m["g5"], m["h5"] = 58, 57, 56, 55, 54, 53, 52, 51
	m["a6"], m["b6"], m["c6"], m["d6"], m["e6"], m["f6"], m["g6"], m["h6"] = 68, 67, 66, 65, 64, 63, 62, 61
	m["a7"], m["b7"], m["c7"], m["d7"], m["e7"], m["f7"], m["g7"], m["h7"] = 78, 77, 76, 75, 74, 73, 72, 71
	m["a8"], m["b8"], m["c8"], m["d8"], m["e8"], m["f8"], m["g8"], m["h8"] = 88, 87, 86, 85, 84, 83, 82, 81
	// pieceMap
	p := make(map[int]string)
	p[18], p[17], p[16], p[15], p[14], p[13], p[12], p[11] = "a1", "b1", "c1", "d1", "e1", "f1", "g1", "h1"
	p[28], p[27], p[26], p[25], p[24], p[23], p[22], p[21] = "a2", "b2", "c2", "d2", "e2", "f2", "g2", "h2"
	p[38], p[37], p[36], p[35], p[34], p[33], p[32], p[31] = "a3", "b3", "c3", "d3", "e3", "f3", "g3", "h3"
	p[48], p[47], p[46], p[45], p[44], p[43], p[42], p[41] = "a4", "b4", "c4", "d4", "e4", "f4", "g4", "h4"
	p[58], p[57], p[56], p[55], p[54], p[53], p[52], p[51] = "a5", "b5", "c5", "d5", "e5", "f5", "g5", "h5"
	p[68], p[67], p[66], p[65], p[64], p[63], p[62], p[61] = "a6", "b6", "c6", "d6", "e6", "f6", "g6", "h6"
	p[78], p[77], p[76], p[75], p[74], p[73], p[72], p[71] = "a7", "b7", "c7", "d7", "e7", "f7", "g7", "h7"
	p[88], p[87], p[86], p[85], p[84], p[83], p[82], p[81] = "a8", "b8", "c8", "d8", "e8", "f8", "g8", "h8"
	// Map of unicode fonts
	r := make(map[string]string)
	r["p"], r["P"] = "\u2659", "\u265F"
	r["b"], r["B"] = "\u2657", "\u265D"
	r["n"], r["N"] = "\u2658", "\u265E"
	r["r"], r["R"] = "\u2656", "\u265C"
	r["q"], r["Q"] = "\u2655", "\u265B"
	r["k"], r["K"] = "\u2654", "\u265A"
	r["."] = "\u00B7"

	// Rows
	rows := make(map[int][8]int)
	rows[1] = [8]int{18, 17, 16, 15, 14, 13, 12, 11}
	rows[2] = [8]int{28, 27, 26, 25, 24, 23, 22, 21}
	rows[3] = [8]int{38, 37, 36, 35, 34, 33, 32, 31}
	rows[4] = [8]int{48, 47, 46, 45, 44, 43, 42, 41}
	rows[5] = [8]int{58, 57, 56, 55, 54, 53, 52, 51}
	rows[6] = [8]int{68, 67, 66, 65, 64, 63, 62, 61}
	rows[7] = [8]int{78, 77, 76, 75, 74, 73, 72, 71}
	rows[8] = [8]int{88, 87, 86, 85, 84, 83, 82, 81}
	
	// Regex Pattern for matching pgn moves
	pgnPattern, _ := regexp.Compile(`([PNBRQK]?[a-h]?[1-8]?)x?([a-h][1-8])([\+\?\!]?)|O(-?O){1,2}`)
	fenPattern, _ := regexp.Compile(`([PNBRQKpnbrqk\d]{1,8}/[PNBRQKpnbrqk\d]{1,8}/[PNBRQKpnbrqk\d]{1,8}/[PNBRQKpnbrqk\d]{1,8}/[PNBRQKpnbrqk\d]{1,8}/[PNBRQKpnbrqk\d]{1,8}/[PNBRQKpnbrqk\d]{1,8}/[PNBRQKpnbrqk\d]{1,8})\s(w|b)\s([KQkq-]{1,4})\s([a-h][36]|-)\s\d\s([1-9]?[1-9])`)
	return Board{
		board:      b,
		castle:     []byte(`KQkq`),
		pgnMap:     m,
		pieceMap:   p,
		pieces:     r,
		rows: rows,
		toMove:     "w",
		score:      "*",
		moves:      1,
		pgnPattern: pgnPattern,
		fenPattern: fenPattern,
	}
}

// Return PNG String.
// Concatnate headers and history
func (b *Board) PgnString() string {
	return b.headers + b.pgn
}

// Return a string printable board
// White position
// TODO: Rotation
func (b *Board) String() string {
	//white square
	// odds odd
	// evens even
	r := make(map[int]bool) // black squares
	r[17], r[15], r[13], r[11], r[28],r[26], r[24], r[22], r[37], r[35], r[33], r[31], r[48], r[46], r[44], r[42], r[57], r[55], r[53], r[51], r[68], r[66], r[64], r[62], r[77], r[75], r[73], r[71], r[88], r[86], r[84], r[82] = false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false
	game := b.board
	p := b.pieces
	var printBoard string
	for i := 89; i > 10; i-- {
		if i%10 == 0 {
			printBoard += "\n"
			continue
		} else if (i+1)%10 == 0 {
			printBoard += string(game[i]) + ": "
			continue
		}
		if game[i] == '.' {
			_, ok := r[i]
			if ok { // white square
				printBoard += "|" + " " + "|"
			} else { // black square
				printBoard += "|" + "\u2591" + "|"
			}
		} else {
			printBoard += "|" + p[string(game[i])] + "|"
		}
	}

	printBoard += "\n"
	printBoard += "   :a::b::c::d::e::f::g::h:\n"
	return printBoard
}

// Wrapper in for standard notation positions.
// TODO: use two coordinates, include piece value
// e2e4
func (b *Board) standardWrapper(orig, dest string) error {
	e := b.Move(b.pgnMap[orig], b.pgnMap[dest])
	if e != nil {
		return e
	}
	return nil
}

// Validate move
// Change byte values to new values.
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
		e := b.validPawn(orig, dest, d)
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

	// Look for Checkmate
	// Check all possibl moves after a check?
	isCheck = b.isPlayerInCheck()
	if isCheck {
		isCheckMate := false
		origs, _ := b.SearchForValid()
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

// Updates board, useless without validation
func (b *Board) updateBoard(orig, dest int,
	val byte, isEmpassant, isCastle bool) {
	isWhite := b.toMove == "w"
	var isPromotion bool
	//var attEmpassant bool

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

// Check if current player is in Check
// TODO: Change to upper case
func (b *Board) isPlayerInCheck() bool {
	isWhite := b.toMove == "w"
	for idx, val := range b.board {
		if val == 'K' && b.isUpper(idx) && isWhite {
			return b.isInCheck(idx)
		}
		if val == 'k' && !b.isUpper(idx) && !isWhite {
			return b.isInCheck(idx)
		}
	}
	return false
}

// Check if target King is in Check
func (b *Board) isInCheck(target int) bool {
	isWhite := b.isUpper(target)
	k := b.board[target]

	// store all the orig of the opponents pieces
	attackers := make([]int, 0, 16)

	for idx, val := range b.board {
		matchWhite, _ := regexp.MatchString(`[PNBRQK]`,
			string(val))
		matchBlack, _ := regexp.MatchString(`[pnbrqk]`,
			string(val))
		if isWhite && matchBlack {
			attackers = append(attackers, idx)
		} else if !isWhite && matchWhite { // black
			attackers = append(attackers, idx)
		}
	}
	//fmt.Println("white ", isWhite, "attackers ", attackers, "king", k)
	// check for valid attacks
	for _, val := range attackers {
		p := string(bytes.ToUpper(b.board[val : val+1]))
		switch {
		case p == "P":
			e := b.validPawn(val, target, k)
			if e == nil {
				//fmt.Println("Pawn check")
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
				//fmt.Println("Bishop check")
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

// Check: right-color, origin-empty, attack-enemy
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

// validate Pawn Move
func (b *Board) validPawn(orig int, dest int, d byte) error {
	err := errors.New("Illegal Pawn Move")
	var remainder int
	var empOffset int
	var empTarget byte
	// Whose turn
	if b.toMove == "w" {
		remainder = dest - orig
		empOffset = -10 // where the empassant piece should be
		empTarget = 'p'
	} else if b.toMove == "b" {
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
		if orig > 28 && b.toMove == "w" { // Only from 2nd rank
			return err
		} else if orig < 70 && b.toMove == "b" {
			return err
		}
	case remainder == 9 || remainder == 11:
		if b.board[dest] == d && d != '.' {
			// Proper attack
		} else if b.board[dest] == d && dest+empOffset == b.empassant {
			// Empassant attack
			if b.board[dest+empOffset] != empTarget {
				return err
			}
		} else {
			return err
		}
	default:
		return errors.New("Not valid Pawn move.")
	}
	return nil
}

// Validate Knight move.
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

// Validate Bishop move.
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

// Validate rook move.
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

// Validate queen move.
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
	// check if anything is inbetween

	return nil
}

// Validate king move.
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

// ParseMove infers origin and destination
// coordinates from a pgn notation move. Check
// and Check Mate notations will be added automatically.
// TODO: disambiguiation
func (b *Board) ParseMove(move string) error {
	move = strings.TrimRight(move, "\r\n") // prepare for input
	// Variables
	var piece string    // find move piece
	var orig int        // find origin coord of move
	var square string   // find pgnMap key of move
	var attacker string // left of x
	//var precise string // for multiple possibilities
	var target byte // the piece to move, in proper case
	var column int
	var row [8]int

	columns := make(map[string]int)
	columns["a"] = 8
	columns["b"] = 7
	columns["c"] = 6
	columns["d"] = 5
	columns["e"] = 4
	columns["f"] = 3
	columns["g"] = 2
	columns["h"] = 1

	// Status
	isCastle := false
	isWhite := b.toMove == "w"
	isCapture, _ := regexp.MatchString(`x`, move)

	res := b.pgnPattern.FindStringSubmatch(move)
	if res == nil && move != "O-O" && move != "O-O-O" {
		return errors.New("invalid input")
	} else if move == "O-O" || move == "O-O-O" {
		// be nice
		isCastle = true
	}
	// Either is capture or not
	if isCapture {
		attacker = res[1]
		if attacker == strings.ToLower(attacker) {
			piece = "P"
			column = columns[attacker]
		} else { // if  upper case, forcement a piece
			piece = string(attacker[0])
			if len(res[1]) > 1 {
				i, err := strconv.Atoi(string(res[1][1]))
				if err == nil && i != 0 {
					row = b.rows[i]
				} else {
					c := string(res[1][1])
					column = columns[c]
				}
			}
		}
		square = res[2]
	} else if isCastle {
		if move == "O-O" {
			piece = "K"
			if isWhite {
				square = "h1"
			} else {
				square = "h8"
			}
		} else if move == "O-O-O" {
			piece = "K"
			if isWhite {
				square = "a1"
			} else {
				square = "a8"
			}
		}
	} else { // No x
		chars := len(move)
		if chars == 2 {
			piece = "P"
			square = res[2]
		} else if chars == 3 && move != "0-0" {
			// Breaks when e44 is entered...
			piece = res[1]
			square = res[2] //move[0]
		} else if chars == 4 {
			if len(res[1]) > 1 {
				i, err := strconv.Atoi(string(res[1][1]))
				if err == nil && i != 0 {
					row = b.rows[i]
				} else {
					c := string(res[1][1])
					column = columns[c]
				}
			}
			piece = string(res[1][0]) // remove second char
			square = res[2]
		} else {
			return errors.New("Not enough input")
		}
	}
	// the presumed destination
	dest := b.pgnMap[square]
	// The piece will be saved as case sensitive byte
	if b.toMove == "b" {
		if piece != "" {
			target = []byte(strings.ToLower(piece))[0]
		} else {
			return errors.New("Invalid Input")
		}
	} else {
		if piece != "" {
			target = []byte(piece)[0]
		} else {
			return errors.New("Invalid Input")
		}
	}
	switch {
	case piece == "P": // Pawn Parse
		var possibilities [2]int // two potentional origins
		// TODO: Allow for empassant take
		if b.toMove == "w" {
			if isCapture {
				possibilities[0],
					possibilities[1] = dest-9,
					dest-11
			} else {
				possibilities[0],
					possibilities[1] = dest-10,
					dest-20
			}
		} else { // is black to move
			if isCapture {
				possibilities[0],
					possibilities[1] = dest+9,
					dest+11
			} else {
				possibilities[0],
					possibilities[1] = dest+10,
					dest+20
			}
		}
		// Disambiguate
		firstDisambig := column != 0 && possibilities[0]%10 == column
		secondDisambig := column != 0 && possibilities[1]%10 == column
		firstPoss := b.board[possibilities[0]] == target
		secondPoss := b.board[possibilities[1]] == target
		if firstPoss && firstDisambig && isCapture {
			orig = possibilities[0]
		} else if secondPoss && secondDisambig && isCapture {
			orig = possibilities[1]
		} else if firstPoss {
			orig = possibilities[0]
		} else if secondPoss {
			orig = possibilities[1]
		}
	case piece == "N": // Knight Parse
		var possibilities [8]int
		// TODO: assume no precision
		// Change to possibilities[]
		possibilities[0], possibilities[1],
			possibilities[2], possibilities[3],
			possibilities[4], possibilities[5],
			possibilities[6], possibilities[7] = dest+21,
			dest+19, dest+12, dest+8, dest-8,
			dest-12, dest-19, dest-21
	LoopKnight:
		for _, possibility := range possibilities {
			if column != 0 { // Disambiguate
				disambig := possibility%10 == column
				if b.board[possibility] == target && disambig {
					orig = possibility
					break LoopKnight
				}
			} else if row[0] != 0 { // Disambiguate
				for _, r := range row {
					disambig := b.board[possibility] == target &&
						b.board[possibility] == byte(r)
					if disambig {
						orig = possibility
						break LoopKnight
					}
				}
			} else {
				if b.board[possibility] == target {
					orig = possibility
					break LoopKnight
				}
			}
		}
	case piece == "B": // Bishop Parse
		var possibilities [14]int
		ticker := 0
		// a8 - h1
		for i := dest + 9; i < 90; i += 9 {
			if (i+1)%10 == 0 { // hits boarder
				break
			}
			possibilities[ticker] = i
			ticker++
		}
		for i := dest - 9; i > 10; i -= 9 {
			if (i+1)%10 == 0 { // hits boarder
				break
			}
			possibilities[ticker] = i
			ticker++
		}
		// a1 - h8 Vector
		for i := dest + 11; i < 90; i += 11 {
			if (i+1)%10 == 0 { // hits boarder
				break
			}
			possibilities[ticker] = i
			ticker++
		}
		for i := dest - 11; i > 10; i -= 11 {
			if i%10 == 0 {
				break
			}
			possibilities[ticker] = i
			ticker++
		}
		// Find piece origin
		for _, possibility := range possibilities {
			if b.board[possibility] == target {
				orig = possibility
				break
			}
		}
	case piece == "R": // Rook Parse
		var possibilities [14]int
		ticker := 0
		// Horizontal Vector
		for i := dest + 10; i < 90; i += 10 {
			possibilities[ticker] = i
			ticker++
		}
		for i := dest - 10; i > 10; i -= 10 {
			possibilities[ticker] = i
			ticker++
		}
		// Vertical Vector
		for i := dest + 1; i < 90; i++ {
			if (i+1)%10 == 0 { // hits boarder
				break
			}
			possibilities[ticker] = i
			ticker++
		}
		for i := dest - 1; i > 10; i-- {
			if i%10 == 0 {
				break
			}
			possibilities[ticker] = i
			ticker++
		}
	LoopRook:
		for _, possibility := range possibilities {
			if column != 0 { // Disambiguate
				disambig := possibility%10 == column
				if b.board[possibility] == target && disambig {
					orig = possibility
					err := b.validRook(orig, dest)
					if err != nil {
						continue
					}
					break LoopRook
				}
			} else if row[0] != 0 { // Disambiguate
				for _, r := range row {
					disambig := b.board[possibility] == target &&
						possibility == r
					if disambig {
						orig = possibility
						err := b.validRook(orig, dest)
						if err != nil {
							continue
						}
						break LoopRook
					}
				}
			} else {
				if b.board[possibility] == target {
					orig = possibility
					err := b.validRook(orig, dest)
					if err != nil {
						continue
					}
					break LoopRook
				}
			}
		}
	case piece == "Q": // Queen Parse
		for idx, possibility := range b.board {
			if possibility == target {
				orig = idx
				break
			}
		}
	case piece == "K": // King Parse
		var possibilities [8]int
		if isCastle {
			if isWhite {
				orig = 14
			} else {
				orig = 84
			}
			break
		}
		possibilities[0], possibilities[1],
			possibilities[2], possibilities[3],
			possibilities[4], possibilities[5],
			possibilities[6], possibilities[7] = dest+10,
			dest+11, dest+1, dest+9, dest-10,
			dest-11, dest-1, dest-9
		for _, possibility := range possibilities {
			if b.board[possibility] == target {
				orig = possibility
				break
			}
		}
	}
	// Move the Piece
	// - Validate Move in Board.Move()
	if b.board[dest] != '.' && !isCapture && !isCastle {
		return errors.New("Not the proper capture syntax")

	}
	if orig != 0 && dest != 0 {
		err := b.Move(orig, dest)
		if err == nil {
			matchCheck, _ := regexp.MatchString(`+`, move)
			matchMate, _  := regexp.MatchString(`#`, move)
			// Update pgn History
			if b.toMove == "b" {
				b.pgn += strconv.Itoa(b.moves) + ". "
			}
			b.pgn += (move)
			if b.checkmate && matchMate {
				b.pgn += "# "
			} else if b.check && matchCheck {

				b.pgn += "+ "
			} else {
				b.pgn += " " // add space
			}
		}
		return err
	} else {
		return errors.New("No such move")
	}
}

// LoadPgn reads a pgn match.
// TODO: ignore header strings eg [White].
func (b *Board) LoadPgn(match string) error {
	result := b.pgnPattern.FindAllString(match, -1)
	for _, val := range result {
		err := b.ParseMove(val)
		if err != nil {
			return err
		}
	}
	return nil
}

// LoadFen Parse FEN string and update Board.board.
func (b *Board) LoadFen(fen string) error {
	// Treat fen input
	fen = strings.TrimRight(fen, "\r\n")
	matches := b.fenPattern.MatchString(fen)
	if !matches {
		return errors.New("Invalid FEN")
	}
	res := b.fenPattern.FindStringSubmatch(fen)
	posCount := 88 // First position to fill
	//
	for _, val := range res[1] { // res[1] {
		// First, make sure it's a relevant position
		if (posCount % 10) == 0 { // || (posCount+1)%10 == 0 {
			posCount -= 2
		}
		// Check if there are Empty Squares
		num, e := strconv.Atoi(string(val))
		if e == nil {
			for j := 0; j < num; j++ {
				b.board[posCount] = '.'
				posCount--
			}
			continue
		}
		if val == '/' {
			continue
		}
		// Is regular Piece
		b.board[posCount] = byte(val)
		posCount--
	}
	// Update Board value
	// Castle
	b.castle = []byte("----")
	for _, v := range res[3] {
		switch {
		case v == 'K':
			b.castle[0] = 'K'
			break
		case v == 'Q':
			b.castle[1] = 'Q'
			break
		case v == 'k':
			b.castle[2] = 'k'
			break
		case v == 'q':
			b.castle[3] = 'q'
			break
		}
	}
	// Empassant
	if res[4] != "-" {
		emp := b.pgnMap[res[4]]
		if emp < 40 {
			// white pawn
		} else if emp > 60 {
		} else {
			return errors.New("Invalid FEN empassant")
		}
	} else {
		b.empassant = 0
	}
	// Turn
	turns, _ := strconv.Atoi(res[5])
	b.moves = turns
	b.fen = fen
	b.toMove = res[2]
	return nil
}

// Position returns string FEN position.
func (b *Board) Position() string {
	pos := ""
	emp := "-"
	zeroTicker := 0
	for i := 88; i > 10; i-- {
		if i%10 == 0 || (i+1)%10 == 0 {
			continue
		}
		// Cycle backwards and tally empty squares
		if b.board[i] == '.' {
			zeroTicker++
		} else if zeroTicker > 0 && b.board[i] != '.' {
			pos += strconv.Itoa(zeroTicker)
			pos += string(b.board[i])
			zeroTicker = 0
		} else {
			pos += string(b.board[i])
		}
		if (i-1)%10 == 0 && i > 10 { // hit edge
			if zeroTicker > 0 {
				pos += strconv.Itoa(zeroTicker)
			}
			zeroTicker = 0
			if i > 11 {
				pos += "/"
			}
		}

	}

	if b.empassant != 0 {
		if b.toMove == "w" {
			emp = b.pieceMap[b.empassant+10]
		} else {
			emp = b.pieceMap[b.empassant-10]
		}
	}
	b.fen = pos + " " + b.toMove + " " + string(b.castle[:4]) + " " + emp + " 0 " + strconv.Itoa(b.moves)
	return b.fen
}

// Play game in terminal
func main() {
	board := NewBoard()
	PlayGame(board)
}

// PlayGame takes user input and commands.
// See ui/clichess.go for more robust client.
func PlayGame(board Board) {
	var turn string
	welcome := `
********
    go-chess

    /~ |_ _  _ _
    \_|||(/__\_\

`
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(welcome)
	fmt.Print(board.String())

Loop:
	for {
		if board.toMove == "w" {
			turn = "White"
		} else {
			turn = "Black"
		}
		fmt.Print(turn, " to move: ")
		input, _ := reader.ReadString('\n')
		isCmd, _ := regexp.MatchString(`/`, input)
		if isCmd {
			input = strings.TrimRight(input, "\r\n")
			switch {
			case input == "/quit":
				break Loop //os.Exit(1)
			case input == "/new":
				board = NewBoard()
				fmt.Print(board.String())
			case input == "/print":
				fmt.Print(board.String())
			default:
				fmt.Println("Mysterious input")
			}
			continue
		}
		e := board.ParseMove(input)
		if board.toMove == "w" {
			turn = "White"
		} else {
			turn = "Black"
		}
		fmt.Println("\n-------------------")
		// TODO use formats.
		if e != nil {
			fmt.Printf("   [Error: %v]\n", e)
		}
		fmt.Print(board.String())
		if board.checkmate {
			fmt.Println("****CheckMate!****")
		} else if board.check {
			fmt.Println("****Check!****")
		}
	}
	fmt.Println("\nGood Game.")
}

// Stats returns program data of current game
// in map[string]string.
// Todo, replace with exported struct attirbutes.
func (b *Board) Stats() map[string]string {
	_ = b.Position()
	m := make(map[string]string)
	m["turn"] = b.toMove
	m["move"] = strconv.Itoa(b.moves)
	m["castling"] = string(b.castle)
	m["position"] = b.fen
	m["history"] = b.pgn
	m["check"] = strconv.FormatBool(b.check)
	m["headers"] = b.headers
	m["score"] = b.score
	m["checkmate"] = strconv.FormatBool(b.checkmate)
	return m
}

// SetHeaders sets pgnHeaders for a pgn export.
func (b *Board) SetHeaders(w, bl string) {
	w = strings.TrimRight(w, "\r\n")
	bl = strings.TrimRight(bl, "\r\n")
	y, m, d := time.Now().Date()
	ye, mo, da := strconv.Itoa(y), strconv.Itoa(int(m)),
		strconv.Itoa(d)
	white := "[White \"" + w + "\"]"
	black := "[Black \"" + bl + "\"]"
	date := "[Date \"" + ye + "." + mo + "." + da + "\"]"
	result := `[Result "*"]`
	b.headers = white + "\n" + black + "\n" + date + "\n" + result + "\n"
}

// Coordinates prints the int values used
// for Board.Move()
func (b *Board) Coordinates() {
	// TODO Rotate Board
	game := b.board
	var printBoard string
	for i := 89; i > 10; i-- {
		if i%10 == 0 {
			printBoard += "\n"
			continue
		} else if (i+1)%10 == 0 {
			printBoard += string(game[i]) + ": "
			continue
		}
		printBoard += "|" + strconv.Itoa(i) + "|"
	}
	printBoard += "\n"
	printBoard += "   :a ::b ::c ::d ::e ::f ::g ::h :\n"
	fmt.Println(printBoard)
}

// isUpper is a wrapper to check if byte in
// Board.board is upper case.
// If Uppercase, it is either white player
// [TODO] or it is empty square.
func (b Board) isUpper(x int) bool {
	//compare = []byte(bytes.ToLower(b))[0]
	compare := byte(unicode.ToUpper(rune(b.board[x])))
	if b.board[x] == compare {
		return true
	} else {
		return false
	}
}

/*
Search
*/

// SearchForValid returns lists of int coordinates
// for valid origins and destinations of the current
// player.
// TODO: castling is a mess
func (b *Board) SearchForValid() ([]int, []int) {
	isWhite := b.toMove == "w"
	movers := make([]int, 0, 16)
	targets := make([]int, 0, 64)
	origs := make([]int, 0, 16)
	dests := make([]int, 0, 64)
	var d byte
	var validMoveCount int
	var king int	
	for idx, val := range b.board {
		if king == 0 && isWhite && val == 'K' {
			king = idx
		} else if king == 0 && !isWhite && val == 'k' {
			king = idx
		}
		// Only look for 64 squares
		if idx%10 == 0 || (idx+1)%10 == 0 || idx > 88 || idx < 11 {
			continue
		}

		// TODO:
		// This is why Castle search-valid doens't work
		if isWhite && b.isUpper(idx) && val != '.' {
			movers = append(movers, idx)

		} else if !isWhite && !b.isUpper(idx) && val != '.' {
			movers = append(movers, idx)
		} else {
			targets = append(targets, idx)
		}
	}
	for _, val := range movers {
		p := string(bytes.ToUpper(b.board[val : val+1]))
		for _, target := range targets {
			if val == 14 && target == 11 {
				//fmt.Println("this should work")
			}
			if target%10 == 0 || (target+1)%10 == 0 || target > 88 || target < 11 {
				continue
			}
			if isWhite {
				d = []byte(bytes.ToLower(b.board[target : target+1]))[0]
			} else {
				d = []byte(bytes.ToLower(b.board[target : target+1]))[0]
			}
			switch {
			case p == "P":
				e := b.validPawn(val, target, d)
				if e == nil {
					validMoveCount++
					origs = append(origs, val)
					dests = append(dests, target)
				}
			case p == "N":
				e := b.validKnight(val, target)
				if e == nil {
					origs = append(origs, val)
					dests = append(dests, target)
					validMoveCount++
				}
			case p == "B":
				e := b.validBishop(val, target)
				if e == nil {
					validMoveCount++
					origs = append(origs, val)
					dests = append(dests, target)
				}
			case p == "R":
				e := b.validRook(val, target)
				if e == nil {
					origs = append(origs, val)
					dests = append(dests, target)
					validMoveCount++

				}
			case p == "Q":
				e := b.validQueen(val, target)
				if e == nil {
					origs = append(origs, val)
					dests = append(dests, target)
					validMoveCount++
				}
			case p == "K":
				e := b.validKing(val, target, false)
				if e == nil {
					origs = append(origs, val)
					dests = append(dests, target)
					validMoveCount++
				}
				// Castle
				e = b.validKing(val, target, true)
				// Castling validation is totally messed up
				//fmt.Println(e, target)
				if e == nil {
					validMoveCount++
					origs = append(origs, val)
					dests = append(dests, target)
				}
			}
		}
	}
	//fmt.Println("Valid move count: ", validMoveCount)
	var realOrigs, realDests []int
	for idx, _ := range origs {
		if b.board[origs[idx]] == 'k' || b.board[origs[idx]] == 'K' {
			king = dests[idx]
		}
		possible := *b
		bCopy := make([]byte, 120)
		copy(bCopy, b.board)
		possible.board = bCopy
		err := possible.Move(origs[idx], dests[idx])
		isCheck := possible.isInCheck(king)
		if err != nil || isCheck {
			// this deletes?, nooo it doesn't..
			// Or atleast it doesn't keep the order.
			// TODO: fix or delete this code:
			origs = append(origs[:idx], origs[idx:]...)
			dests = append(dests[:idx], dests[idx:]...)
		} else {
			realOrigs = append(realOrigs, origs[idx])
			realDests = append(realDests, dests[idx])
		}
	}
	return realOrigs, realDests
}

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

/*
TODO: Evaluate
*/
