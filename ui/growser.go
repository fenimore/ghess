/*
TODO:
Plenty
Add database?
*/
package main

import (
	"fmt"
	"github.com/polypmer/ghess"
	"html/template"
	"net/http"
	"time"
)

// ChessHandler for persistant game
// in handler functions
type ChessHandler struct {
	g    ghess.Board
	init bool
}

// Chessboard struct is for
// sending chessboard.js info
type ChessBoard struct {
	Board    string
	Fen      string
	Pgn      string
	Move     string
	wToMove  bool
	Feedback string
}

// boardHandler for playing game
// Takes url param pgn move
func (h *ChessHandler) playGameHandler(w http.ResponseWriter,
	r *http.Request) {
	// If no board, redirect to board
	move := ""
	feedback := ""
	if r.Method == "POST" {
		r.ParseForm()
		move = r.Form.Get("move")
	} else if r.Method == "GET" {
		move = r.URL.Path[len("/play/"):]
	}
	e := h.g.ParseMove(move)
	if e != nil {
		feedback = e.Error()
	}
	// Display with GUI chessboard.js
	// TODO: Get White and Black
	// And rotate board accordingly
	b := ChessBoard{Board: h.g.String(), Fen: h.g.Position(), Pgn: h.g.PgnString(), Feedback: feedback}
	t, err := template.ParseFiles("templates/board.html")
	if err != nil {
		fmt.Printf("Error %s Templates", err)
	}
	t.Execute(w, b)
}

// newGameHandler creates a new Board object
// and links to /board/ route
func (h *ChessHandler) newGameHandler(w http.ResponseWriter,
	r *http.Request) {
	h.g = ghess.NewBoard()
	h.init = true
	fmt.Fprintln(w, "<a href=/board>New Game Created</a>")
}

// /board/ route, displays board and new move form.
func (h *ChessHandler) showGameHandler(w http.ResponseWriter,
	r *http.Request) {
	if h.init != true {
		// TODO: FAILURE
		http.Redirect(w, r, "/new/", http.StatusSeeOther)
	}
	b := ChessBoard{Board: h.g.String(), Fen: h.g.Position(), Pgn: h.g.PgnString()}
	t, err := template.ParseFiles("templates/board.html")
	if err != nil {
		fmt.Printf("Error %s Templates", err)
	}
	t.Execute(w, b)
}

// Index page, link to new game
func (h *ChessHandler) indexHandler(w http.ResponseWriter,
	r *http.Request) {
	fmt.Fprintln(w, "<a href=/new >New Game</a>")
}

// AJAX Handler for Updating board
// This does not update all open connections
// TODO: Websockets!?
func (h *ChessHandler) makeMoveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}
	// Get Form Value
	field := r.FormValue("move")
	// make Move
	h.g.ParseMove(field)
	pos := h.g.Position()
	// Write to Client
	w.Write([]byte(pos))
	// TODO: write to all open connections
}

func main() {
	// So HandlFunc takes a custom Handler
	// Which is forcement takes into a reader and writer
	// and then it will print whatever is written to the
	// writer
	PORT := "0.0.0.0:8080"
	h := new(ChessHandler)
	// Server Routes
	http.HandleFunc("/", h.indexHandler)            // link to new game
	http.HandleFunc("/play/", h.playGameHandler)    // deprecated
	http.HandleFunc("/board/", h.showGameHandler)   // view
	http.HandleFunc("/new/", h.newGameHandler)      // new board
	http.HandleFunc("/makemove", h.makeMoveHandler) //ajax
	// Handle Static Files
	// TODO: Combine into one function?
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("static/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("static/js"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("static/img"))))
	// Why can't I just link them all in the same Handle()?
	//http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	//Listen and Server on PORT 8080
	fmt.Printf("Listening on %s\n", PORT)
	http.ListenAndServe(PORT, nil)
}

func getPanel(m map[string]string) string {
	return "|Move:  " + m["move"] + "     Turn: " + m["turn"] +
		"\n|Check: " + m["check"] + " Castle: " + m["castling"] +
		"\n|Mate:  " + m["checkmate"] + " Score: " + m["score"] + "\n"
}
