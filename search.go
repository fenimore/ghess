package ghess

import "bytes"

// SearchValid finds two arrays, of all valid possible
// destinations and origins.
// TODO
func (b *Board) SearchValid() ([]int, []int) {
	movers := make([]int, 0, 16)
	targets := make([]int, 0, 63) // There will only ever be 63 open squares
	origs := make([]int, 0, 16)
	dests := make([]int, 0, 63)
	validO := make([]int, 0, 16)
	validD := make([]int, 0, 63)
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

		// TODO:
		// This is why Castle search return in valid doesn't work
		if b.toMove == "w" && b.isUpper(idx) && val != '.' {
			movers = append(movers, idx)
		} else if b.toMove == "b" && !b.isUpper(idx) && val != '.' {
			movers = append(movers, idx)
		} else {
			targets = append(targets, idx)
		}
	}

	for _, idx := range movers {
		p := bytes.ToUpper(b.board[idx : idx+1])[0]
		for _, target := range targets {
			// TODO: Check for Castling
			var e error
			switch p {
			case 'P':
				e = b.validPawn(idx, target)
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
				}
				e = b.validKing(idx, target, true)
				//fmt.Println("King")
			}
			if e == nil {
				origs = append(origs, idx)
				dests = append(dests, target)
			}
		}
	}

	// Check if it moves into Check
	for i := 0; i < len(origs); i++ {
		// Check if King is the piece moved
		if b.board[origs[i]] == 'k' || b.board[origs[i]] == 'K' {
			king = dests[i]
		}
		// Copy board to "make move" as the move must be made
		// in order to test for check
		possible := CopyBoard(b)
		err := possible.Move(origs[i], dests[i])
		isCheck := possible.isInCheck(king)
		if err == nil && !isCheck {
			validO = append(validO, origs[i])
			validD = append(validD, dests[i])
		}
	}

	return validO, validD
}
