/*
Search and Evaluation


*/
package main

import (
	"errors"
	"fmt"
	"bytes"
	//"time"
)

//var Board [128]byte

/*

    ' rnbqkbnr\n'  #  20 - 29
    ' pppppppp\n'  #  30 - 39
    ' ........\n'  #  40 - 49
    ' ........\n'  #  50 - 59
    ' ........\n'  #  60 - 69
    ' ........\n'  #  70 - 79
    ' PPPPPPPP\n'  #  80 - 89
    ' RNBQKBNR\n'  #  90 - 99


*/

type Board struct {
	board []byte
	castle []byte
	empassant int
	coord map[string]int
	A1, H1, A8, H8 int
	toMove string
	moves int
}


func NewBoard() Board {
	b := make([]byte, 120)
	fmt.Println("initializing board")
	b[91], b[92], b[93], b[94], b[95], b[96], b[97], b[98] = 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'
	b[19], b[29], b[39], b[49], b[59], b[69], b[79], b[89] = '1', '2', '3', '4', '5', '6', '7', '8'

	// starting position
	b = []byte(`           RNBKQBNR  PPPPPPPP  ........  ........  ........  ........  pppppppp  rnbkqbnr                                `)
	b[91], b[92], b[93], b[94], b[95], b[96], b[97], b[98] = 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'
	b[19], b[29], b[39], b[49], b[59], b[69], b[79], b[89] = '1', '2', '3', '4', '5', '6', '7', '8'
	
	cas := []byte(`KQkq`)

	m := make(map[string]int)
	m["a1"], m["b1"], m["c1"], m["d1"], m["e1"], m["f"], m["g1"], m["h1"]  = 11, 12, 13, 14, 15, 16, 17, 18
	
	return Board{
		board: b,
		A1: 91,
		A8: 21,
		H1: 98,
		H8: 28,
		castle: cas,
		coord: m,
		toMove: "w",
	}
}


func (b *Board) String()string {
	var printBoard string
	//fmt.Println(string(b.castle))
	//fmt.Println(b.coord["a1"])
	for idx, val := range b.board {
		if idx < 100 && idx > 10{
			if idx % 10 != 0 && idx <90{
				if (idx+1) % 10 !=0{// why doesn't an || work?
					printBoard += "|"+ string(val)+"| "
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

func (b *Board) Move(orig, dest int) error {
	val := b.board[orig]
	//err := nil
	if b.toMove == "w" {
		// check that orig is Upper
		fmt.Println("white to move")
		o := []byte(bytes.ToUpper(b.board[orig:orig+1]))[0]
		
		if b.board[orig] == o && o != '.' {
			fmt.Println("not a peroid")
		} else {
			return errors.New("not your turn")
		}
	} else if b.toMove == "b" {
		// check if orig is Lower
		fmt.Println("Black to Move")
		o := []byte(bytes.ToLower(b.board[orig:orig+1]))[0]
		
		if b.board[orig] == o && o != '.' {
			fmt.Println("not a peroid")
		} else {
			return errors.New("not your turn")
		}
	}
	b.board[orig] = '.'
	// is it empty
	// are the squares leading to it empty
	//
	// Change to Move
	b.board[dest] = val

	fmt.Println(string(b.board[dest]))
	if b.toMove == "w" {
		b.toMove = "b"
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

	e := board.Move(24, 44)
	fmt.Print(board.String())
	if e != nil {
		fmt.Print("not nil")
	}
	e = board.Move(24, 44)
	fmt.Print(board.String())
	if e != nil {
		fmt.Print("not nil")
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
