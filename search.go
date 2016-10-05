package ghess

import (
	"bytes"
	"fmt"
)

// SearchValid finds two arrays, of all valid possible
// destinations and origins.
func (b *Board) SearchValid() ([]int, []int) {
	movers := make([]int, 0, 16)
	targets := make([]int, 0, 63) // There will only ever be 63 open squares
	origs := make([]int, 0, 16)
	dests := make([]int, 0, 16)

	// Find and sort pieces:
	for idx, val := range b.board {
		// Only look for 64 squares
		if idx%10 == 0 || (idx+1)%10 == 0 || idx > 88 || idx < 11 {
			continue
		}

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
			switch p {
			case 'P':
				fmt.Println(target)
			case 'N':
			case 'B':
			case 'R':
			case 'Q':
			case 'K':

			}
		}
	}

	return origs, dests
}

// Run it in goroutine
//func (b *Board) CheckTargets(orig int, targets []int) {

//}
