/*
Go Chess Engine
Fenimore Love 2016
GPLv3
TODO: Search and Evaluation
TODO: Fen PGN reading
TODO: Fen output

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
	fen string
	pgn string
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
	m["a1"], m["b1"], m["c1"], m["d1"], m["e1"], m["f1"], m["g1"], m["h1"]  = 11, 12, 13, 14, 15, 16, 17, 18
	m["a2"], m["b2"], m["c2"], m["d2"], m["e2"], m["f2"], m["g2"], m["h2"]  = 21, 22, 23, 24, 25, 26, 27, 28
	m["a3"], m["b3"], m["c3"], m["d3"], m["e3"], m["f3"], m["g3"], m["h3"]  = 31, 32, 33, 34, 35, 36, 37, 38
	m["a4"], m["b4"], m["c4"], m["d4"], m["e4"], m["f4"], m["g4"], m["h4"]  = 41, 42, 43, 44, 45, 46, 47, 48
	m["a5"], m["b5"], m["c5"], m["d5"], m["e5"], m["f5"], m["g5"], m["h5"]  = 51, 52, 53, 54, 55, 56, 57, 58
	m["a6"], m["b6"], m["c6"], m["d6"], m["e6"], m["f6"], m["g6"], m["h6"]  = 61, 62, 63, 64, 65, 66, 67, 68
	m["a7"], m["b7"], m["c7"], m["d7"], m["e7"], m["f7"], m["g7"], m["h7"]  = 71, 72, 73, 74, 75, 76, 77, 78
	m["a8"], m["b8"], m["c8"], m["d8"], m["e8"], m["f8"], m["g8"], m["h8"]  = 81, 82, 83, 84, 85, 86, 87, 88


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

		}
	}
	return printBoard
}
	
/*
Rules and validation
*/
// Wrapper in standard notation
func (b *Board) pgnMove(orig, dest string) error{
	e := b.Move(b.coord[orig], b.coord[dest])
	if e != nil {
		return e
	}
	return nil
}
// Move byte value to new position
func (b *Board) Move(orig, dest int) error {
	val := b.board[orig]
	var o byte // supposed starting square
	var d byte // supposed destination
	if b.toMove == "w" {
		// check that orig is Upper
		// and dest is Enemy or Empty
		fmt.Println("White moves")
		o = []byte(bytes.ToUpper(b.board[orig:orig+1]))[0]
		d = []byte(bytes.ToLower(b.board[dest:dest+1]))[0]
	} else if b.toMove == "b" {
		// check if orig is Lower
		// and dest is Enemy or Empty
		fmt.Println("Black moves")
		o = []byte(bytes.ToLower(b.board[orig:orig+1]))[0]
		d = []byte(bytes.ToUpper(b.board[dest:dest+1]))[0]
	}
	// Check if it is the right turn
	if b.board[orig] != o {
		return errors.New("Not your turn")
	}
	// Check if Origin is Empty
	if o == '.' {
		return errors.New("Empty square")
	}
	// Check if destination is Enemy
	if b.board[dest] != d{ // 
		return errors.New("Can't attack your own piece")
	}
	p := string(bytes.ToUpper(b.board[orig:orig+1]))
	fmt.Print(orig, dest)
	switch {
	case p == "P":
		fmt.Print("is pawn")
		e := b.validPawn(orig, dest, d)
		if e != nil{
			return e
		}
	case p == "N":
		fmt.Print("is knight")
	case p == "B":
		fmt.Print("is bishop")
	case p == "R":
		fmt.Print("is rook")
	case p == "Q":
		fmt.Print("is queen")
	case p == "K":
		fmt.Print("is king")
	}
	// Dest:
	// are the squares leading up to it empty
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

// validate Pawn Move
func (b *Board) validPawn(orig int, dest int, d byte) error {
	err := errors.New("Illegal Move")
	var remainder int
	if b.toMove == "w" {
		remainder = dest-orig
	} else if b.toMove == "b"{
		remainder = orig-dest
	}
	if remainder == 10 {
		// regular move
	} else if remainder == 20 { // two spaces
		// double starter move
		if orig > 28 && b.toMove == "w"{ // Only from 2nd rank
			return err
		}
		if orig < 70 && b.toMove == "b"{
			return err
		}
	} else if remainder == 9 || remainder == 11 {
		// Attack vector
		// check if b.board[orig+10] == '.'
		if b.board[dest] == d &&  d != '.'{
			// Proper attack
		} else {
			return err
		}
	}
	return nil
}

// Valid Knight

/*
TODO: Export fen
TODO: Parse fen
TODO: Parse pgn
*/



func main() {
	board := NewBoard()
	board.Coordinates()
	//TestGame(board)
	e := board.pgnMove("e2", "e4")
	if e != nil {
		fmt.Print(e)
	}
	fmt.Print(board.String())
	e = board.pgnMove("d7", "d5")
	if e != nil {
		fmt.Print(e)
	}
	fmt.Print(board.String())
	e = board.pgnMove("e4", "d5")
	if e != nil {
		fmt.Print(e)
	}
	fmt.Print(board.String())
}


/*
Testing method
*/

func TestGame(b Board) {
	fmt.Println("New Game")
	fmt.Println(b.toMove, string(b.castle))
	fmt.Print(b.String())
	e := b.Move(22, 42) // white
	fmt.Print(b.String())
	if e != nil {
		fmt.Print("not nil")
	}
	e = b.Move(72, 52)// black
	fmt.Print(b.String())
	if e != nil {
		fmt.Print("not nil")
	}
	e = b.Move(23, 43) // white
	fmt.Print(b.String())
	if e != nil {
		fmt.Print("not nil")
	}
	e = b.Move(52, 43) // black
	fmt.Print(b.String())
	if e != nil {
		fmt.Println(e)
	}
	e = b.Move(21, 43) // white
	fmt.Print(b.String())
	if e != nil {
		fmt.Print(e)
	}
}


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
