// Package ghess is a chess engine. This file concerns the Search
// and keeps a tally of board tension.
package ghess

import (
	"bytes"
	"strconv"
)

// SearchValid finds two arrays, of all valid possible
// destinations and origins. These are int coordinates
// which point to the index of the byte slice Board.board
func (b *Board) SearchValidSlowly() ([]int, []int) {
	movers := make([]int, 0, 16)
	targets := make([]int, 0, 64)
	origs := make([]int, 0, 16)
	dests := make([]int, 0, 64)

	// Find and sort pieces:
	for idx, val := range b.board {
		// Only look for 64 squares
		/*if idx%10 == 0 || (idx+1)%10 == 0 ||
			idx > 88 || idx < 11 {
			continue
		}*/
		if val != ' ' && b.toMove == "w" && b.isUpper(idx) && val != '.' {
			movers = append(movers, idx)
		} else if val != ' ' && b.toMove == "b" &&
			!b.isUpper(idx) && val != '.' {
			movers = append(movers, idx)
		} else if val != ' ' {
			targets = append(targets, idx)
		}
	}

	// Add Castle (edge squares to targets
	if b.toMove == "w" {
		if b.castle[1] == 'Q' {
			targets = append(targets, 18)
		}
		if b.castle[0] == 'K' {
			targets = append(targets, 11)
		}
	}
	if b.toMove == "b" {
		if b.castle[3] == 'q' {
			targets = append(targets, 88)
		}
		if b.castle[2] == 'k' {
			targets = append(targets, 81)
		}
	}
	// Check for Valid attacks
	for _, idx := range movers {
		for _, target := range targets {
			var e error
			var isCheck bool
			possible := CopyBoard(b)
			e = possible.Move(idx, target)
			if e == nil {
				// Can't move into check
				isCheck = possible.isOpponentInCheck()
			}
			if e == nil && !isCheck {
				origs = append(origs, idx)
				dests = append(dests, target)
			}
		}
	}
	return origs, dests
}

// SearchValid finds two arrays, of all valid possible
// destinations and origins. These are int coordinates
// which point to the index of the byte slice Board.board
func (b *Board) SearchValid() ([]int, []int) {
	movers := make([]int, 0, 16)
	origs := make([]int, 0, 16)
	dests := make([]int, 0, 64)

	// Find and sort pieces:
	for idx, val := range b.board {
		// Only look for 64 squares
		if val != ' ' && b.toMove == "w" && b.isUpper(idx) && val != '.' {
			movers = append(movers, idx)
		} else if val != ' ' && b.toMove == "b" &&
			!b.isUpper(idx) && val != '.' {
			movers = append(movers, idx)
		}
	}
	for _, idx := range movers {
		switch b.board[idx] {
		case 'p', 'P':
			o, d := b.searchPawn(idx)
			origs = append(origs, o...)
			dests = append(dests, d...)
		case 'n', 'N':
			o, d := b.searchKnight(idx)
			origs = append(origs, o...)
			dests = append(dests, d...)
		case 'b', 'B':
			o, d := b.searchBishop(idx)
			origs = append(origs, o...)
			dests = append(dests, d...)
		case 'r', 'R':
			o, d := b.searchRook(idx)
			origs = append(origs, o...)
			dests = append(dests, d...)
		case 'q', 'Q':
			o, d := b.searchQueen(idx)
			origs = append(origs, o...)
			dests = append(dests, d...)
		case 'k', 'K':
			o, d := b.searchKing(idx)
			origs = append(origs, o...)
			dests = append(dests, d...)
		}
	}
	return origs, dests
}

// checkForCheck avoids validating the move, but simply
// checks if a new position would put the opponent in check
// (which is illegal)
func (b *Board) searchOk(o, d int) bool {
	possible := CopyBoard(b)
	possible.updateBoard(o, d, b.board[o], false, false)
	if possible.isOpponentInCheck() {
		return false
	}
	return true
}

func (b *Board) searchPawn(orig int) ([]int, []int) {
	isWhite := b.isUpper(orig)
	origs := make([]int, 0)
	dests := make([]int, 0)
	var possibilities [4]int
	// TODO: Check for Empassant
	if isWhite {
		possibilities[0] = orig + 10
		possibilities[1] = orig + 20
		possibilities[2] = orig + 11
		possibilities[3] = orig + 9
	} else {
		possibilities[0] = orig - 10
		possibilities[1] = orig - 20
		possibilities[2] = orig - 11
		possibilities[3] = orig - 9
	}

	for idx, possibility := range possibilities {
		if possibility > 89 || possibility < 11 {
			continue
		}
		switch b.board[possibility] {
		case ' ':
			continue
		case '.':
			if idx == 0 {
				if b.searchOk(orig, possibility) {
					origs = append(origs, orig)
					dests = append(dests, possibility)
				}
			} else if idx == 1 && (orig < 29 || orig > 69) {
				if isWhite {
					if b.board[orig+10] != '.' {
						continue
					}
				} else {
					if b.board[orig-10] != '.' {
						continue
					}
				}
				if b.searchOk(orig, possibility) {
					origs = append(origs, orig)
					dests = append(dests, possibility)
				}
			}
		default: // if it's a piece
			if (idx == 2 || idx == 3) &&
				isWhite && !b.isUpper(possibility) {
				if b.searchOk(orig, possibility) {
					origs = append(origs, orig)
					dests = append(dests, possibility)
				}
			} else if (idx == 2 || idx == 3) &&
				!isWhite && b.isUpper(possibility) {
				if b.searchOk(orig, possibility) {
					origs = append(origs, orig)
					dests = append(dests, possibility)
				}
			}
		}
	}
	return origs, dests
}

func (b *Board) searchKnight(orig int) ([]int, []int) {
	isWhite := b.isUpper(orig)
	origs := make([]int, 0)
	dests := make([]int, 0)
	var possibilities [8]int
	possibilities[0], possibilities[1],
		possibilities[2], possibilities[3],
		possibilities[4], possibilities[5],
		possibilities[6], possibilities[7] = orig+21,
		orig+19, orig+12, orig+8, orig-8,
		orig-12, orig-19, orig-21
PossLoop:
	for _, possibility := range possibilities {
		if possibility > 89 || possibility < 11 {
			continue PossLoop
		}
		switch b.board[possibility] {
		case ' ':
			continue PossLoop
		case '.':
			if b.searchOk(orig, possibility) {
				origs = append(origs, orig)
				dests = append(dests, possibility)
			}
		default: // if it's a piece
			if isWhite && !b.isUpper(possibility) {
				if b.searchOk(orig, possibility) {
					origs = append(origs, orig)
					dests = append(dests, possibility)
				}
			} else if !isWhite && b.isUpper(possibility) {
				if b.searchOk(orig, possibility) {
					origs = append(origs, orig)
					dests = append(dests, possibility)
				}
			}
		}

	}
	return origs, dests
}

func (b *Board) searchBishop(orig int) ([]int, []int) {
	isWhite := b.isUpper(orig)
	origs := make([]int, 0)
	dests := make([]int, 0)
	// Not by possibility.. but rather check diagnols..
a1h8Loop:
	for target := orig + 9; target < 89; target = target + 9 {
		// Should stop when off the board
		switch b.board[target] {
		case ' ':
			break a1h8Loop
		case '.':
			if b.searchOk(orig, target) {
				origs = append(origs, orig)
				dests = append(dests, target)
			}
		case 'p', 'n', 'b', 'r', 'q', 'k':
			if isWhite {
				if b.searchOk(orig, target) {
					origs = append(origs, orig)
					dests = append(dests, target)
				}
			}
			break a1h8Loop
		case 'P', 'N', 'B', 'R', 'Q', 'K':
			if !isWhite {
				if b.searchOk(orig, target) {
					origs = append(origs, orig)
					dests = append(dests, target)
				}
			}
			break a1h8Loop
		}
	}
h8a1Loop:
	for target := orig - 9; target > 10; target = target - 9 {
		// Should stop when off the board
		switch b.board[target] {
		case ' ':
			break h8a1Loop
		case '.':
			if b.searchOk(orig, target) {
				origs = append(origs, orig)
				dests = append(dests, target)
			}
		case 'p', 'n', 'b', 'r', 'q', 'k':
			if isWhite {
				if b.searchOk(orig, target) {
					origs = append(origs, orig)
					dests = append(dests, target)
				}
			}
			break h8a1Loop
		case 'P', 'N', 'B', 'R', 'Q', 'K':
			if !isWhite {
				if b.searchOk(orig, target) {
					origs = append(origs, orig)
					dests = append(dests, target)
				}
			}
			break h8a1Loop
		}
	}

h1a8Loop:
	for target := orig + 11; target < 89; target = target + 11 {
		// Should stop when off the board
		switch b.board[target] {
		case ' ':
			break h1a8Loop
		case '.':
			if b.searchOk(orig, target) {
				origs = append(origs, orig)
				dests = append(dests, target)
			}
		case 'p', 'n', 'b', 'r', 'q', 'k':
			if isWhite {
				if b.searchOk(orig, target) {
					origs = append(origs, orig)
					dests = append(dests, target)
				}
			}
			break h1a8Loop
		case 'P', 'N', 'B', 'R', 'Q', 'K':
			if !isWhite {
				if b.searchOk(orig, target) {
					origs = append(origs, orig)
					dests = append(dests, target)
				}
			}
			break h1a8Loop
		}
	}
a8h1Loop:
	for target := orig - 11; target > 10; target = target - 11 {
		// Should stop when off the board
		switch b.board[target] {
		case ' ':
			break a8h1Loop
		case '.':
			if b.searchOk(orig, target) {
				origs = append(origs, orig)
				dests = append(dests, target)
			}
		case 'p', 'n', 'b', 'r', 'q', 'k':
			if isWhite {
				if b.searchOk(orig, target) {
					origs = append(origs, orig)
					dests = append(dests, target)
				}
			}
			break a8h1Loop
		case 'P', 'N', 'B', 'R', 'Q', 'K':
			if !isWhite {
				if b.searchOk(orig, target) {
					origs = append(origs, orig)
					dests = append(dests, target)
				}
			}
			break a8h1Loop
		}
	}
	return origs, dests
}

func (b *Board) searchRook(orig int) ([]int, []int) {
	isWhite := b.isUpper(orig)
	origs := make([]int, 0)
	dests := make([]int, 0)
RightLoop:
	for i := orig + 1; !((i+1)%10 == 0); i = i + 1 {
		switch b.board[i] {
		case ' ':
			break RightLoop
		case '.':
			if b.searchOk(orig, i) {
				origs = append(origs, orig)
				dests = append(dests, i)
			}
		case 'p', 'n', 'b', 'r', 'q', 'k':
			if isWhite {
				if b.searchOk(orig, i) {
					origs = append(origs, orig)
					dests = append(dests, i)
				}
			}
			break RightLoop
		case 'P', 'N', 'B', 'R', 'Q', 'K':
			if !isWhite {
				if b.searchOk(orig, i) {
					origs = append(origs, orig)
					dests = append(dests, i)
				}
			}
			break RightLoop
		}
	}
LeftLoop:
	for i := orig - 1; !(i%10 == 0); i = i - 1 {
		switch b.board[i] {
		case ' ':
			break LeftLoop
		case '.':
			if b.searchOk(orig, i) {
				origs = append(origs, orig)
				dests = append(dests, i)
			}
		case 'p', 'n', 'b', 'r', 'q', 'k':
			if isWhite {
				if b.searchOk(orig, i) {
					origs = append(origs, orig)
					dests = append(dests, i)
				}
			}
			break LeftLoop
		case 'P', 'N', 'B', 'R', 'Q', 'K':
			if !isWhite {
				if b.searchOk(orig, i) {
					origs = append(origs, orig)
					dests = append(dests, i)
				}
			}
			break LeftLoop
		}
	}
UpVerLoop:
	for i := orig + 10; i < 88; i = i + 10 {

		switch b.board[i] {
		case ' ':
			break UpVerLoop
		case '.':
			if b.searchOk(orig, i) {
				origs = append(origs, orig)
				dests = append(dests, i)
			}
		case 'p', 'n', 'b', 'r', 'q', 'k':
			if isWhite {
				if b.searchOk(orig, i) {
					origs = append(origs, orig)
					dests = append(dests, i)
				}
			}
			break UpVerLoop
		case 'P', 'N', 'B', 'R', 'Q', 'K':
			if !isWhite {
				if b.searchOk(orig, i) {
					origs = append(origs, orig)
					dests = append(dests, i)
				}
			}
			break UpVerLoop
		}
	}
DownVerLoop:
	for i := orig - 10; i > 10; i = i - 10 {
		// Should stop when off the board
		switch b.board[i] {
		case ' ':
			break DownVerLoop
		case '.':
			if b.searchOk(orig, i) {
				origs = append(origs, orig)
				dests = append(dests, i)
			}
		case 'p', 'n', 'b', 'r', 'q', 'k':
			if isWhite {
				if b.searchOk(orig, i) {
					origs = append(origs, orig)
					dests = append(dests, i)
				}
			}
			break DownVerLoop
		case 'P', 'N', 'B', 'R', 'Q', 'K':
			if !isWhite {
				if b.searchOk(orig, i) {
					origs = append(origs, orig)
					dests = append(dests, i)
				}
			}
			break DownVerLoop
		}
	}
	return origs, dests
}

func (b *Board) searchQueen(orig int) ([]int, []int) {
	//isWhite := b.isUpper(orig)
	origs := make([]int, 0)
	dests := make([]int, 0)

	o, d := b.searchBishop(orig)
	origs = append(origs, o...)
	dests = append(dests, d...)
	o, d = b.searchRook(orig)
	origs = append(origs, o...)
	dests = append(dests, d...)

	return origs, dests
}

func (b *Board) searchKing(orig int) ([]int, []int) {
	isWhite := b.isUpper(orig)
	origs := make([]int, 0)
	dests := make([]int, 0)
	var possibilities [8]int
	possibilities[0], possibilities[1],
		possibilities[2], possibilities[3],
		possibilities[4], possibilities[5],
		possibilities[6], possibilities[7] = orig+11,
		orig+9, orig+10, orig+1, orig-11,
		orig-10, orig-9, orig-1
	// Check kings close
	for _, possibility := range possibilities {
		switch b.board[possibility] {
		case ' ':
			continue
		case '.':
			if b.searchOk(orig, possibility) {
				origs = append(origs, orig)
				dests = append(dests, possibility)
			}
		default: // if it's a piece
			if isWhite && !b.isUpper(possibility) {
				if b.searchOk(orig, possibility) {
					origs = append(origs, orig)
					dests = append(dests, possibility)
				}
			} else if !isWhite && b.isUpper(possibility) {
				if b.searchOk(orig, possibility) {
					origs = append(origs, orig)
					dests = append(dests, possibility)
				}
			}
		}

	}
	if isWhite {
		if b.castle[1] == 'Q' && b.board[orig+1] == '.' {
			possible := CopyBoard(b)
			e := possible.Move(orig, 18)
			if e == nil {
				origs = append(origs, orig)
				dests = append(dests, 18)
			}
		}
		if b.castle[0] == 'K' && b.board[orig-1] == '.' {
			possible := CopyBoard(b)
			e := possible.Move(orig, 11)
			if e == nil {
				origs = append(origs, orig)
				dests = append(dests, 11)
			}
		}
	} else {
		if b.castle[3] == 'q' && b.board[orig+1] == '.' {
			possible := CopyBoard(b)
			e := possible.Move(orig, 88)
			if e == nil {
				origs = append(origs, orig)
				dests = append(dests, 88)
			}
		}
		if b.castle[2] == 'k' && b.board[orig-1] == '.' {
			possible := CopyBoard(b)
			e := possible.Move(orig, 81)
			if e == nil {
				origs = append(origs, orig)
				dests = append(dests, 81)
			}
		}
	}
	return origs, dests
}

// Tension returns a map of which squares are attacked
// by which side. Negative for black Positive for white.
// TODO: For now it wont take not moving out of check into account?
// TODO: Empassant isn't counted in tension
func (b *Board) Tension() map[int]int {
	tension := make(map[int]int)
	whites := make([]int, 0, 16) // white movers black targets
	blacks := make([]int, 0, 16) // black movers white targets
	// Todo, do I need blanks?
	blanks := make([]int, 0, 63) // everyone's targets

	// Find and sort pieces:
	for idx, val := range b.board {
		// Only look for 64 squares
		if idx%10 == 0 || (idx+1)%10 == 0 || idx > 88 || idx < 11 {
			continue
		}

		// TODO:
		// This is why Castle search return in valid doesn't work
		if b.isUpper(idx) && val != '.' {
			whites = append(whites, idx)
		} else if !b.isUpper(idx) && val != '.' {
			blacks = append(blacks, idx)
		} else {
			blanks = append(blanks, idx)
		}
	}

	// Increment For Squares white is attacking
	for _, idx := range whites {
		whiteTargets := append(blanks, blacks...)
		whiteTargets = append(whiteTargets, whites...)
		p := bytes.ToUpper(b.board[idx : idx+1])[0]
		// Check for Pawn attacks, cause they're weird:
		if p == 'P' {
			if (idx+9)%10 != 0 {
				tension[idx+9]++
			}
			if (idx+11)%10 != 0 {
				tension[idx+11]++
			}
		}
	WhiteValidator:
		for _, target := range whiteTargets {
			// TODO: Check for Castling
			var e error
			switch p {
			case 'P':
				continue WhiteValidator
			case 'N':
				e = b.validKnight(idx, target)
			case 'B':
				e = b.validBishop(idx, target)
			case 'R':
				e = b.validRook(idx, target)
			case 'Q':
				e = b.validQueen(idx, target)
			case 'K':
				e = b.validKing(idx, target, false)
				// Don't check for castling? Cause that's not pressure
			}
			if e == nil {
				// Is valid
				tension[target]++
			}
		}
	}

	//Decrement for Squares black is attacking
	for _, idx := range blacks {
		blackTargets := append(blanks, whites...)
		blackTargets = append(blackTargets, blacks...)
		p := bytes.ToUpper(b.board[idx : idx+1])[0]
		// Check for Pawn attacks, cause they're weird:
		if p == 'P' {
			if (idx-9)%10 != 0 {
				tension[idx-9]--
			}
			if (idx-11)%10 != 0 {
				tension[idx-11]--
			}

		}
	BlackValidator:
		for _, target := range blackTargets {
			// TODO: Check for Castling
			var e error
			switch p {
			case 'P': // Cause pawns are weird
				continue BlackValidator
			case 'N':
				e = b.validKnight(idx, target)
			case 'B':
				e = b.validBishop(idx, target)
			case 'R':
				e = b.validRook(idx, target)
			case 'Q':
				e = b.validQueen(idx, target)
			case 'K':
				e = b.validKing(idx, target, false)
				//e = b.validKing(idx, target, true)
				//fmt.Println("King")
			}
			if e == nil {
				tension[target]--
			}
		}
	}
	return tension
}

func (b *Board) TensionSum() int {
	tension := b.Tension()
	var sum int
	for _, v := range tension {
		sum += v
	}
	return sum
}

// StringTension Prints the board with numbers according to the
// amount of attacks by either black or white, negative for black
// and positive for white.
func (b *Board) StringTension() string {
	tension := b.Tension()
	var nums [8]byte // somehow print these?
	nums[0], nums[1], nums[2], nums[3], nums[4], nums[5], nums[6], nums[7] = '1', '2', '3', '4', '5', '6', '7', '8'

	//p := b.pieces
	var printBoard string
	j := 7
	for i := 89; i > 10; i-- {
		if i%10 == 0 {
			printBoard += string(nums[j]) + ": " + "\n"
			j--
			continue
		} else if (i+1)%10 == 0 {
			continue
		}
		if tension[i] == 0 {
			printBoard += "| 0|"

		} else {
			if tension[i] < 0 {
				printBoard += "|" + strconv.Itoa(tension[i]) + "|"
			} else {
				printBoard += "| " + strconv.Itoa(tension[i]) + "|"
			}
		}
	}

	printBoard += string(nums[j]) + ": " + "\n"
	printBoard += ": a:: b:: c:: d:: e:: f:: g:: h:\n"
	return printBoard
}

// Deprecated:
func (b *Board) SearchValidSlow() ([]int, []int) {
	movers := make([]int, 0, 16)
	targets := make([]int, 0, 64)
	origs := make([]int, 0, 16)
	dests := make([]int, 0, 64)
	validO := make([]int, 0, 16)
	validD := make([]int, 0, 64)
	isWhite := b.toMove == "w"
	var king int

	// Find and sort pieces:
	for idx, val := range b.board {
		// Only look for 64 squares
		if idx%10 == 0 || (idx+1)%10 == 0 || idx > 88 || idx < 11 {
			continue
		}

		if val == 'K' || val == 'k' {
			if isWhite && val == 'K' {
				king = idx
			} else if !isWhite && val == 'k' {
				king = idx
			}
		}

		if b.toMove == "w" && b.isUpper(idx) && val != '.' {
			movers = append(movers, idx)
		} else if b.toMove == "b" && !b.isUpper(idx) && val != '.' {
			movers = append(movers, idx)

		} else {
			targets = append(targets, idx)

		}

	}

	// Add Castle to targets
	if b.toMove == "w" {
		if b.castle[1] == 'Q' {
			targets = append(targets, 18)
		}
		if b.castle[0] == 'K' {
			targets = append(targets, 11)
		}
	}
	if b.toMove == "b" {
		if b.castle[3] == 'q' {
			targets = append(targets, 88)
		}
		if b.castle[2] == 'k' {
			targets = append(targets, 81)
		}
	}

	// Check for Valid attacks
	for _, idx := range movers {
		p := bytes.ToUpper(b.board[idx : idx+1])[0]
		for _, target := range targets {
			var e error
			switch p {
			case 'P':
				// why am I not looking at pawns?
				//e = b.validPawn(idx, target)
			case 'N':
				e = b.validKnight(idx, target)
			case 'B':
				e = b.validBishop(idx, target)
			case 'R':
				e = b.validRook(idx, target)
			case 'Q':
				e = b.validQueen(idx, target)
			case 'K':

				e = b.validKing(idx, target, false)
				if e == nil {
					origs = append(origs, idx)
					dests = append(dests, target)
					continue
				}
				err := b.validKing(idx, target, true)
				if err == nil {
					origs = append(origs, idx)
					dests = append(dests, target)
					continue
				}
			}
			if e == nil {
				origs = append(origs, idx)
				dests = append(dests, target)
			}
		}
	}

	// Check if it moves into Check
	//for i := 0; i < len(origs); i++ {
	for i, v := range origs {
		var k int
		// Check if King is the piece moved
		if b.board[v] == 'k' || b.board[v] == 'K' {
			k = dests[i]
		} else {
			k = king
		}
		// Copy board to "make move" as the move must be made
		// in order to test for check
		possible := CopyBoard(b)
		err := possible.Move(origs[i], dests[i])
		isCheck := possible.isInCheck(k)
		//if dests[i] == 75 && isCheck {
		//fmt.Println(origs[i], "What?", king)
		//}
		if err == nil && !isCheck {
			validO = append(validO, origs[i])
			validD = append(validD, dests[i])
		}
	}

	return validO, validD
}
