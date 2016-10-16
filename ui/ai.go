package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/polypmer/ghess"
)

type Chess struct {
	g    ghess.Board
	init bool
}

// Index page, link to new game
func (h *Chess) Index(w http.ResponseWriter,
	r *http.Request) {
	html := `
<html>
<link href="/css/style.css" rel="stylesheet">
<h1>Ghess Index</h1>
<a href=/new/black >New Game Computer Vs Human</a><br>
<a href=/new/white >New Game Human Vs Computer</a><br>
<hr>
<a href=/board >View Current Game</a><br>
<br><br><br>
<a href="https://github.com/polypmer/go-chess">Source Code</a>
</html>
`
	fmt.Fprintln(w, html)
}

// newGameHandler creates a new Board object
// And redirects to a board instance /board
func (h *Chess) newGame(w http.ResponseWriter,
	r *http.Request) {
	h.g = ghess.NewBoard()
	h.init = true
	choice := r.URL.Path[len("/new/"):]
	if choice == "white" {
		http.Redirect(w, r, "/board", http.StatusSeeOther)
	} else {
		h.g.Move(24, 44)
		http.Redirect(w, r, "/board", http.StatusSeeOther)
	}
}

func (h *Chess) playGame(w http.ResponseWriter,
	r *http.Request) {

	t, err := template.ParseFiles("templates/computer.html")
	if err != nil {
		fmt.Printf("Error %s Templates", err)
	}

	t.Execute(w, h.g.Position())
}

type Move struct {
	Position string `json:"position"`
	Message  string `json:"message"`
}

func (h *Chess) movePiece(w http.ResponseWriter,
	r *http.Request) {
	if r.Method != "POST" {
		fmt.Println("Not Post wtf")
	}
	mv := &Move{}
	move := strings.Split(r.URL.Path[len("/move/"):], "/")
	err := h.g.ParseStand(move[0], move[1])
	if err != nil {
		mv = &Move{
			Position: h.g.Position(),
			Message:  "> That's Not a Valid Move!",
		}
	} else {
		now := time.Now()
		state, err := ghess.MiniMax(0, 3, ghess.GetState(&h.g))
		if err != nil {
			fmt.Println("> Minimax broken")
		}
		h.g.Move(state.Init[0], state.Init[1])
		msg := fmt.Sprintf("> Your Move, my move took %s",
			time.Since(now))
		mv = &Move{
			Position: h.g.Position(),
			Message:  msg,
		}
	}
	js, _ := json.Marshal(mv)
	w.Write([]byte(js))

	//

	// parse move

}

func main() {
	PORT := ":8888"
	h := new(Chess)
	h.g = ghess.NewBoard() // This means only playin' one game attime
	// Server Routes
	http.HandleFunc("/", h.Index)          // link to new game
	http.HandleFunc("/board/", h.playGame) // view
	http.HandleFunc("/new/", h.newGame)    // new board
	http.HandleFunc("/move/", h.movePiece) // ajax call
	// Handle Static Files
	// TODO: Combine into one function?
	http.Handle("/css/", http.StripPrefix("/css/",
		http.FileServer(http.Dir("static/css"))))
	http.Handle("/js/", http.StripPrefix("/js/",
		http.FileServer(http.Dir("static/js"))))
	http.Handle("/img/", http.StripPrefix("/img/",
		http.FileServer(http.Dir("static/img"))))
	// Why can't I just link them all in the same Handle()?
	//http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Printf("Listening on %s\n", PORT)
	http.ListenAndServe(PORT, nil)

}
