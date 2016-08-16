/*
TODO:
Create templates:
Index with ids,
get set player data
templates with forms
add javascript chessboard
*/
package main

import (
	"fmt"
	"github.com/polypmer/ghess"
	"html/template"
	"net/http"
)

type ChessHandler struct {
	g ghess.Board
}

type ChessBoard struct {
	Board string
	Fen   string
	Pgn   string
	Move  string
}

// boardHandler for playing game
// Takes url param pgn move
func (h *ChessHandler) playGameHandler(w http.ResponseWriter,
	r *http.Request) {
	// If no board, redirect to board
	//http.Redirect(w, r, "/new/", http.StatusFound)

	move := r.URL.Path[len("/play/"):]
	e := h.g.ParseMove(move)
	if e != nil {
		fmt.Fprintln(w, e.Error())
	}
	fmt.Fprintln(w, h.g.String())
}

func (h *ChessHandler) newGameHandler(w http.ResponseWriter,
	r *http.Request) {
	h.g = ghess.NewBoard()
	fmt.Fprintln(w, h.g.String())
}

func (h *ChessHandler) showGameHandler(w http.ResponseWriter,
	r *http.Request) {
	// Must it be a pointer?
	b := ChessBoard{Board: h.g.String(), Fen: h.g.Position(), Pgn: h.g.PgnString()}
	t, err := template.ParseFiles("templates/board.html")
	if err != nil {
		fmt.Printf("Error %s Templates", err)
	}
	t.Execute(w, b)
	//fmt.Fprintln(w, getPanel(h.g.Stats())+h.g.String())
}

func main() {
	// So HandlFunc takes a custom Handler
	// Which is forcement takes into a reader and writer
	// and then it will print whatever is written to the
	// writer
	PORT := "0.0.0.0:8080"
	h := new(ChessHandler)

	// Server Part
	http.HandleFunc("/play/", h.playGameHandler)
	http.HandleFunc("/board/", h.showGameHandler)
	http.HandleFunc("/new/", h.newGameHandler)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("static/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("static/js"))))
	fmt.Printf("Listening on %s\n", PORT)
	http.ListenAndServe(PORT, nil)
}

func getPanel(m map[string]string) string {
	return "|Move:  " + m["move"] + "     Turn: " + m["turn"] +
		"\n|Check: " + m["check"] + " Castle: " + m["castling"] +
		"\n|Mate:  " + m["checkmate"] + " Score: " + m["score"] + "\n"
}
