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


type Board struct {
	board []byte
}

func (b Board) String() string {
	return string(b.board)
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
	wth := 8 // +2 for right column + \n
	lng := 8 // +1 for bottom row
	

	b := make([]byte, wth*lng)
	for i := 0; i < 8; i +=2 {
		row0 := i * 8
		row1 := (i + 1) * 8
		for j := 0; j < 8; j += 2 {
			b[row0+j], b[row0+j+1] = 'b', 'w'
			if row1+j+1 <= 64 {
				b[row1+j], b[row1+j+1] = 'w', 'b'
			}
		}
	}
	return Board{
		board: b,
	}
}


func main() {
//	Board[0] = 70
//	Board[1] = 
	///	fmt.Print(Board)
	board := NewBoard()
	fmt.Print(string(board.board))
}
