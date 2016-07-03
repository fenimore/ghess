/*
Go Chess Engine
Fenimore Love 2016
GPLv3
TODO: Search and Evaluation
*/
package main

import (
	"errors"
	"fmt"
	"bytes"
)

// The chessboard type
type Board struct {
	board []byte // piece position
	castle []byte // castle possibility KQkq or ----
	empassant int // square vulnerable to empassant
	coord map[string]int // the pgn format
	toMove string // Next move is w or b
	moves int // the count of moves 
	pieces map[string]string // the unicode fonts
}

// __init__ for Board
func NewBoard() Board {
	b := make([]byte, 120)
	fmt.Println("initializing board")
	// starting position
	b = []byte(`           RNBKQBNR  PPPPPPPP  ........  ........  ........  ........  pppppppp  rnbkqbnr                                `)
	b[91], b[92], b[93], b[94], b[95], b[96], b[97], b[98] = 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'
	b[19], b[29], b[39], b[49], b[59], b[69], b[79], b[89] = '1', '2', '3', '4', '5', '6', '7', '8'
	
	cas := []byte(`KQkq`) // castle possibility
	// Map of PGN notation, incomplete
	m := make(map[string]int)
	m["a1"], m["b1"], m["c1"], m["d1"], m["e1"], m["f"], m["g1"], m["h1"]  = 11, 12, 13, 14, 15, 16, 17, 18
	// Map of unicode fonts
	r := make(map[string]string)
	r["p"] = "\u2659"
	r["P"] = "\u265F"
	r["b"] = "\u2657"
	r["B"] = "\u265D"
	r["n"] = "\u2658"
	r["N"] = "\u265E"
	r["r"] = "\u2656"
	r["R"] = "\u265C"
	r["q"] = "\u2655"
	r["Q"] = "\u265B"
	r["k"] = "\u2654"
	r["K"] = "\u265A"	
	r["."] = "\u2022"
	return Board{
		board: b,
		castle: cas,
		coord: m,
		pieces: r,
		toMove: "w",
		
	}
}

// Return a string of the board
// Todo Unicode chess pieces
func (b *Board) String()string {
	var printBoard string
	//fmt.Println(string(b.castle))
	//fmt.Println(b.coord["a1"])
	for idx, val := range b.board {
		if idx < 100 && idx > 10{
			if idx % 10 != 0 && idx <90{
				if (idx+1) % 10 !=0{// why doesn't an || work?
					font := b.pieces[string(val)]
					printBoard += "|"+ font +"| "
				} else {
					printBoard += ":"+string(val)
				}
			}
		}
		if idx > 90 && idx < 99{
			printBoard += ":"+string(val)+": "
		}
		if idx % 10 == 0 && idx != 0{
			printBoard += "\n"
			fmt.Print("\n")
		}
	}
	return printBoard
}
	
/*
Rules and validation
*/
// Move byte value to new position
func (b *Board) Move(orig, dest int) error {
	val := b.board[orig]
	var o byte
	if b.toMove == "w" {
		// check that orig is Upper
		fmt.Println("white to move")
		o = []byte(bytes.ToUpper(b.board[orig:orig+1]))[0]
	} else if b.toMove == "b" {
		// check if orig is Lower
		fmt.Println("Black to Move")
		o = []byte(bytes.ToLower(b.board[orig:orig+1]))[0]
	}
	if b.board[orig] != o || o == '.' {
		return errors.New("not your turn or its a period")
	}
	// Dest:
	// is it empty
	// are the squares leading up to it empty
	// is i possesed by the right color
	// Piece:
	// validatePawn() etc
	// if w
	// if orig > 30
	// dest - orig == 1
	// else == 1 || 2
	// return valid
	b.board[orig] = '.'
	b.board[dest] = val
	// Update Tickers
	if b.toMove == "w" {
		b.toMove = "b"
		b.moves++ // add one to move 
	} else {
		b.toMove = "w"
	}
	return nil
}





func main() {
	board := NewBoard()
	fmt.Print(board.String())
	//board.Move(24, 44)
	//fmt.Print(board.String())
	//board.Move(74, 54)
	//fmt.Print(board.String())
	//board.Coordinates()

	e := board.Move(22, 42)
	fmt.Print(board.String())
	if e != nil {
		fmt.Print("not nil")
	}
	e = board.Move(24, 42)
	fmt.Print(board.String())
	if e != nil {
		fmt.Print("not nil")
	}
}


/*
Testing method
*/

func (b *Board) Coordinates() {
	fmt.Println(string(b.castle))
	fmt.Println(b.coord["a1"])
	for idx, val := range b.board {
		if idx < 100 && idx > 10{
			if idx % 10 != 0 && idx <90{
				if (idx+1) % 10 !=0{// why doesn't an || work?
					fmt.Print( ":", idx, ": ")
				} else {
					fmt.Print(":",string(val))
				}
			}
		}
		if idx > 90 && idx < 99{
			fmt.Print(": ", string(val), ": ")
		}
		if idx % 10 == 0 && idx != 0{
			fmt.Print("\n")
		}
	}
}
