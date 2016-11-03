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
func (b *Board) SearchValid() ([]int, []int) {
	movers := make([]int, 0, 16)
	targets := make([]int, 0, 64)
	origs := make([]int, 0, 16)
	dests := make([]int, 0, 64)
	//isWhite := b.toMove == "w"
	//var king int

	// Find and sort pieces:
	for idx, val := range b.board {
		// Only look for 64 squares
		if idx%10 == 0 || (idx+1)%10 == 0 || idx > 88 || idx < 11 {
			continue
		}

		if b.toMove == "w" && b.isUpper(idx) && val != '.' {
			movers = append(movers, idx)
		} else if b.toMove == "b" && !b.isUpper(idx) && val != '.' {
			movers = append(movers, idx)

		} else {
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

func (b *Board) SearchValidOrdered() ([]int, []int) {
	movers := make([]int, 0, 16)
	targets := make([]int, 0, 64)
	origs := make([]int, 0, 16)
	dests := make([]int, 0, 64)
	//isWhite := b.toMove == "w"
	//var king int

	// Find and sort pieces:
	for idx, val := range b.board {
		// Only look for 64 squares
		if idx%10 == 0 || (idx+1)%10 == 0 || idx > 88 || idx < 11 {
			continue
		}

		if b.toMove == "w" && b.isUpper(idx) && val != '.' {
			movers = append(movers, idx)
		} else if b.toMove == "b" && !b.isUpper(idx) && val != '.' {
			movers = append(movers, idx)

		} else if val == '.' {
			targets = append(targets, idx)
		} else {
			// Is a possible attack
			targets = append([]int{idx}, targets...) // prepend
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
				isCheck = possible.isOpponentInCheck()
			}
			if e == nil && !isCheck {
				if b.board[target] != '.' {
					origs = append([]int{idx}, origs...)
					dests = append([]int{target}, dests...)
				} else {
					origs = append(origs, idx)
					dests = append(dests, target)
				}
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
