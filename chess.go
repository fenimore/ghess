/*
Search and Evaluation


*/
package main

import (
	"errors"
	"fmt"
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
	A1, H1, A8, H8 int
}

func (b Board) String() string {
	var printBoard string

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
		if idx < 9 && idx > 0{
			fmt.Print(": ", string(val), ": ")
		}
		if idx % 10 == 0 && idx != 0{
			fmt.Print("\n")
		}
	}
			
	return printBoard
}

func (b Board) Set(token byte, x, y int) error {
	idx, err := b.squareAt(x, y)
	if err != nil {
		return err
	}
	b.board[idx] = token
	return nil
}

func (b Board) squareAt(x, y int) (int, error) {
	woff := 16
	foff := (y*2+1)*woff + x*2 + 1

	if foff > len(b.board){
		return 0, errors.New("out of range")
	}
	return  (y*2+1)*woff + x*2 + 1, nil
}

func NewBoard() Board {
	b := make([]byte, 120)
	fmt.Println("initializing board")
	b[1], b[2], b[3], b[4], b[5], b[6], b[7], b[8] = 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'
	b[19], b[29], b[39], b[49], b[59], b[69], b[79], b[89] = '1', '2', '3', '4', '5', '6', '7', '8'
	return Board{
		board: b,
		A1: 91,
		A8: 21,
		H1: 98,
		H8: 28,
	}
}


func main() {
//	Board[0] = 70
//	Board[1] = 
	///	fmt.Print(Board)
	board := NewBoard()
	//fmt.Print(board.String())
	//fmt.Println(string(board.board[1:8]))
	board.String()
	fmt.Println("At 20-29", string(board.board[31]))
}
